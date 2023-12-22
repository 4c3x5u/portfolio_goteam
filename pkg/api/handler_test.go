//go:build utest

package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/token"
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
				req := httptest.NewRequest(httpMethod, "/", nil)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				assert.Equal(t.Error,
					res.StatusCode, http.StatusMethodNotAllowed,
				)

				// Assert that all allowed methods were set in the correct
				// header.
				allowedMethods := res.Header.Get("Access-Control-Allow-Methods")
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
				// Prepare request and response recorder.
				r := httptest.NewRequest(httpMethod, "/", nil)
				r.AddCookie(&http.Cookie{Name: token.AuthName, Value: ""})
				w := httptest.NewRecorder()

				// Handle request with sut and get the result.
				sut.ServeHTTP(w, r)
				res := w.Result()

				// Assert on the status code.
				assert.Equal(t.Error, res.StatusCode, http.StatusOK)

				fakeMethodHandler := methodHandler.(*FakeMethodHandler)
				assert.Equal(t.Error, fakeMethodHandler.InResponseWriter, w)
				assert.Equal(t.Error, fakeMethodHandler.InReq, r)
			})
		}
	})
}
