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
		})
	}
}
