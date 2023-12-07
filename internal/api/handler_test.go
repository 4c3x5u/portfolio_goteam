//go:build utest

package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/auth"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestHandler(t *testing.T) {
	authTokenValidator := &auth.FakeTokenValidator{}
	postHandler := &FakeMethodHandler{}
	deleteHandler := &FakeMethodHandler{}
	patchHandler := &FakeMethodHandler{}
	sut := NewHandler(
		authTokenValidator,
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
				req, err := http.NewRequest(httpMethod, "", nil)
				if err != nil {
					t.Fatal(err)
				}
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

	t.Run("InvalidAuthToken", func(t *testing.T) {
		// Set pre-determinate return values for sut's dependencies.
		authTokenValidator.Sub = ""

		// Prepare request and response recorder.
		req, err := http.NewRequest(http.MethodPost, "", nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()

		// Handle request with sut and get the result.
		sut.ServeHTTP(w, req)
		res := w.Result()

		// Assert on the status code.
		assert.Equal(t.Error, res.StatusCode, http.StatusUnauthorized)

		// Assert the WWW-Authenticate header was set correctly
		assert.Equal(t.Error,
			w.Result().Header.Get("WWW-Authenticate"), "Bearer",
		)
	})

	t.Run("Success", func(t *testing.T) {
		for httpMethod, methodHandler := range sut.methodHandlers {
			t.Run(httpMethod, func(t *testing.T) {
				wantSub := "somesub"

				// Set pre-determinate return values for sut's dependencies.
				authTokenValidator.Sub = wantSub

				// Prepare request and response recorder.
				req, err := http.NewRequest(httpMethod, "", nil)
				if err != nil {
					t.Fatal(err)
				}
				req.AddCookie(&http.Cookie{Name: auth.CookieName, Value: ""})
				w := httptest.NewRecorder()

				// Handle request with sut and get the result.
				sut.ServeHTTP(w, req)
				res := w.Result()

				// Assert on the status code.
				assert.Equal(t.Error, res.StatusCode, http.StatusOK)

				fakeMethodHandler := methodHandler.(*FakeMethodHandler)
				assert.Equal(t.Error, fakeMethodHandler.InResponseWriter, w)
				assert.Equal(t.Error, fakeMethodHandler.InReq, req)
				assert.Equal(t.Error, fakeMethodHandler.InSub, wantSub)
			})
		}
	})
}
