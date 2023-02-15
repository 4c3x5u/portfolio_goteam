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
	"server/db"
	"server/log"

	"github.com/golang-jwt/jwt/v4"
)

func TestLogin(t *testing.T) {
	sut := loginAPI.NewHandler(
		loginAPI.NewValidator(),
		db.NewUserSelector(dbConnPool),
		loginAPI.NewPasswordComparer(),
		auth.NewJWTGenerator(jwtKey),
		log.NewAppLogger(),
	)

	for _, c := range []struct {
		name           string
		username       string
		password       string
		wantStatusCode int
	}{
		{
			name:           "UsnEmpty",
			username:       "",
			password:       "P4ssw@rd123",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "PwdEmpty",
			username:       "bob123",
			password:       "",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "UsnIncorrect",
			username:       "bob321",
			password:       "P4ssw@rd123",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "PwdIncorrect",
			username:       "bob123",
			password:       "P4ssw@rd321",
			wantStatusCode: http.StatusBadRequest,
		},
		{
			name:           "Success",
			username:       "bob123",
			password:       "P4ssw@rd123",
			wantStatusCode: http.StatusOK,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			reqBody, err := json.Marshal(loginAPI.ReqBody{
				Username: c.username,
				Password: c.password,
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

			if c.wantStatusCode == http.StatusOK {
				// assert that the returned JWT is valid and has the correct
				// subject
				var token string
				for _, ck := range res.Cookies() {
					if ck.Name == "auth-token" {
						token = ck.Value
					}
				}

				claims := jwt.RegisteredClaims{}
				if _, err = jwt.ParseWithClaims(
					token, &claims, func(token *jwt.Token) (any, error) {
						return []byte(jwtKey), nil
					},
				); err != nil {
					t.Fatal(err)
				}
				if err = assert.Equal(c.username, claims.Subject); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
