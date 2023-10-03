//go:build itest

package itest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	loginAPI "server/api/login"
	"server/assert"
	"server/auth"
	"server/dbaccess"
	"server/log"

	"github.com/golang-jwt/jwt/v4"
)

func TestLogin(t *testing.T) {
	sut := loginAPI.NewHandler(
		loginAPI.NewValidator(),
		dbaccess.NewUserSelector(db),
		loginAPI.NewPasswordComparator(),
		auth.NewJWTGenerator(jwtKey),
		log.New(),
	)

	for _, c := range []struct {
		name           string
		username       string
		password       string
		wantStatusCode int
		assertFunc     func(*testing.T, *http.Response)
	}{
		{
			name:           "UsnEmpty",
			username:       "",
			password:       "P4ssw@rd123",
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     func(*testing.T, *http.Response) {},
		},
		{
			name:           "PwdEmpty",
			username:       "bob123",
			password:       "",
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     func(*testing.T, *http.Response) {},
		},
		{
			name:           "UsnIncorrect",
			username:       "bob321",
			password:       "P4ssw@rd123",
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     func(*testing.T, *http.Response) {},
		},
		{
			name:           "PwdIncorrect",
			username:       "bob123",
			password:       "P4ssw@rd321",
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     func(*testing.T, *http.Response) {},
		},
		{
			name:           "Success",
			username:       "bob123",
			password:       "P4ssw@rd123",
			wantStatusCode: http.StatusOK,
			assertFunc: func(t *testing.T, res *http.Response) {
				var token string
				for _, ck := range res.Cookies() {
					if ck.Name == "auth-token" {
						token = ck.Value
					}
				}

				claims := jwt.RegisteredClaims{}
				if _, err := jwt.ParseWithClaims(
					token, &claims, func(token *jwt.Token) (any, error) {
						return []byte(jwtKey), nil
					},
				); err != nil {
					t.Fatal(err)
				}
				if err := assert.Equal("bob123", claims.Subject); err != nil {
					t.Error(err)
				}
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			reqBody, err := json.Marshal(map[string]string{
				"username": c.username,
				"password": c.password,
			})
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPost, "", bytes.NewReader(reqBody),
			)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, req)

			res := w.Result()

			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			// Run case-specific assertions.
			c.assertFunc(t, res)
		})
	}
}
