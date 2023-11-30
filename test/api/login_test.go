//go:build itest

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v4"

	loginAPI "github.com/kxplxn/goteam/internal/api/login"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/auth"
	userTable "github.com/kxplxn/goteam/pkg/dbaccess/user"
	"github.com/kxplxn/goteam/pkg/log"
)

// TestLoginHandler tests the http.Handler for the login API route and asserts
// that it behaves correctly during various execution paths.
func TestLoginHandler(t *testing.T) {
	sut := loginAPI.NewPOSTHandler(
		loginAPI.NewValidator(),
		userTable.NewSelector(db),
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
			username:       "team1Member",
			password:       "P4ssw@rd123",
			wantStatusCode: http.StatusOK,
			assertFunc: func(t *testing.T, res *http.Response) {
				cookie := res.Cookies()[0]

				assert.True(t.Error, cookie.Secure)
				assert.Equal(t.Error, cookie.SameSite, http.SameSiteNoneMode)

				claims := jwt.RegisteredClaims{}
				if _, err := jwt.ParseWithClaims(
					cookie.Value, &claims, func(token *jwt.Token) (any, error) {
						return []byte(jwtKey), nil
					},
				); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t.Error, claims.Subject, "team1Member")
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
			req := httptest.NewRequest(
				http.MethodPost, "/", bytes.NewReader(reqBody),
			)
			w := httptest.NewRecorder()

			sut.Handle(w, req, "")

			res := w.Result()

			assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

			// Run case-specific assertions.
			c.assertFunc(t, res)
		})
	}
}
