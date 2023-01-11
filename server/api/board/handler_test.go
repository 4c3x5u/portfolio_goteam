package board

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
)

func TestHandler(t *testing.T) {
	t.Run("MethodNotAllowed", func(t *testing.T) {
		for _, method := range []string{
			http.MethodConnect, http.MethodDelete, http.MethodGet,
			http.MethodHead, http.MethodOptions, http.MethodPatch,
			http.MethodPut, http.MethodTrace,
		} {
			t.Run(method, func(t *testing.T) {
				sut := NewHandler()
				req, err := http.NewRequest(http.MethodGet, "/login", nil)
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
}
