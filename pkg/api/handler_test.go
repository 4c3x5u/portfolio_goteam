//go:build utest

package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestHandler(t *testing.T) {
	postHandler := &FakeMethodHandler{}
	deleteHandler := &FakeMethodHandler{}
	patchHandler := &FakeMethodHandler{}
	sut := NewHandler(
		map[string]MethodHandler{
			http.MethodPost:   postHandler,
			http.MethodDelete: deleteHandler,
			http.MethodPatch:  patchHandler,
		},
	)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		for _, httpMethod := range []string{
			http.MethodConnect, http.MethodGet, http.MethodHead, http.MethodPut,
			http.MethodTrace,
		} {
			t.Run(httpMethod, func(t *testing.T) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest(httpMethod, "/", nil)

				sut.ServeHTTP(w, r)

				resp := w.Result()
				assert.Equal(
					t.Error, resp.StatusCode, http.StatusMethodNotAllowed,
				)
				allowedMethods := resp.Header.Get(
					"Access-Control-Allow-Methods",
				)
				for method := range sut.methodHandlers {
					assert.True(t.Error,
						strings.Contains(allowedMethods, method),
					)
				}
			})
		}
	})

	t.Run("OK", func(t *testing.T) {
		for httpMethod, methodHandler := range sut.methodHandlers {
			t.Run(httpMethod, func(t *testing.T) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest(httpMethod, "/", nil)
				r.AddCookie(&http.Cookie{Name: cookie.AuthName, Value: ""})

				sut.ServeHTTP(w, r)

				resp := w.Result()
				assert.Equal(t.Error, resp.StatusCode, http.StatusOK)
				fakeMethodHandler := methodHandler.(*FakeMethodHandler)
				assert.Equal(t.Error, fakeMethodHandler.InResponseWriter, w)
				assert.Equal(t.Error, fakeMethodHandler.InR, r)
			})
		}
	})
}
