package board

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
)

func TestHandler(t *testing.T) {
	sut := NewHandler()

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

	t.Run("AuthCookieNotFound", func(t *testing.T) {
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
