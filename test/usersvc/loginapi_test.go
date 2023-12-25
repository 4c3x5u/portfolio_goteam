//go:build itest

package test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/golang-jwt/jwt/v4"

	loginAPI "github.com/kxplxn/goteam/internal/usersvc/loginapi"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db/usertbl"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/test"
)

func TestLoginAPI(t *testing.T) {
	sut := loginAPI.NewPostHandler(
		loginAPI.NewValidator(),
		usertbl.NewRetriever(db),
		loginAPI.NewPasswordComparator(),
		test.AuthEncoder,
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
			assertFunc: func(t *testing.T, resp *http.Response) {
				ckAuth := resp.Cookies()[0]

				assert.True(t.Error, ckAuth.Secure)
				assert.Equal(t.Error, ckAuth.SameSite, http.SameSiteNoneMode)

				claims := jwt.MapClaims{}
				if _, err := jwt.ParseWithClaims(
					ckAuth.Value, &claims, func(token *jwt.Token) (any, error) {
						return test.JWTKey, nil
					},
				); err != nil {
					t.Fatal(err)
				}

				_, ok := claims["username"].(string)
				if !ok {
					t.Error()
				}

				_, ok = claims["isAdmin"].(bool)
				if !ok {
					t.Error()
				}

				_, ok = claims["teamID"].(string)
				if !ok {
					t.Error()
				}

				assert.Equal(t.Error,
					claims["username"].(string), "team1Member",
				)
				assert.Equal(t.Error, claims["isAdmin"].(bool), false)
				assert.Equal(t.Error,
					claims["teamID"].(string),
					"afeadc4a-68b0-4c33-9e83-4648d20ff26a",
				)
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost, "/", strings.NewReader(`{
                    "username": "`+c.username+`",
                    "password": "`+c.password+`"
                }`),
			)

			sut.Handle(w, r, "")

			resp := w.Result()
			assert.Equal(t.Error, resp.StatusCode, c.wantStatusCode)
			c.assertFunc(t, resp)
		})
	}
}
