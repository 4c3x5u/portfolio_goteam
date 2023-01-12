package board

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/auth"
)

func TestHandler(t *testing.T) {
	tokenValidator := &auth.FakeTokenValidator{}
	sut := NewHandler(tokenValidator)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		for _, httpMethod := range []string{
			http.MethodConnect, http.MethodDelete, http.MethodGet,
			http.MethodHead, http.MethodOptions, http.MethodPatch,
			http.MethodPut, http.MethodTrace,
		} {
			t.Run(httpMethod, func(t *testing.T) {
				req, err := http.NewRequest(httpMethod, "/board", nil)
				if err != nil {
					t.Fatal(err)
				}
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)

				if err = assert.Equal(http.StatusMethodNotAllowed, w.Result().StatusCode); err != nil {
					t.Error(err)
				}
			})
		}
	})

	for _, c := range []struct {
		name                 string
		tokenValidatorOutErr error
	}{
		{name: "NoAuthCookie", tokenValidatorOutErr: nil},
		{name: "InvalidAuthCookie", tokenValidatorOutErr: errors.New("token validator error")},
	} {
		t.Run(c.name, func(t *testing.T) {
			tokenValidator.OutErr = c.tokenValidatorOutErr

			req, err := http.NewRequest(http.MethodPost, "/board", nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, req)

			if err = assert.Equal(http.StatusUnauthorized, w.Result().StatusCode); err != nil {
				t.Error(err)
			}
		})
	}
}
