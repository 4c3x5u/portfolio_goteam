//go:build utest

package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/auth"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestHandler(t *testing.T) {
	authHeaderReader := &auth.FakeHeaderReader{}
	authTokenValidator := &auth.FakeTokenValidator{}
	postHandler := &FakeMethodHandler{}
	deleteHandler := &FakeMethodHandler{}
	patchHandler := &FakeMethodHandler{}
	sut := NewHandler(
		authHeaderReader,
		authTokenValidator,
		map[string]MethodHandler{
			http.MethodPost:   postHandler,
			http.MethodDelete: deleteHandler,
			http.MethodPatch:  patchHandler,
		},
	)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		for _, httpMethod := range []string{
			http.MethodConnect, http.MethodGet, http.MethodHead,
			http.MethodOptions, http.MethodPut, http.MethodTrace,
		} {
			t.Run(httpMethod, func(t *testing.T) {
				req, err := http.NewRequest(httpMethod, "", nil)
				if err != nil {
					t.Fatal(err)
				}
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				if err = assert.Equal(
					http.StatusMethodNotAllowed, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				// Assert that all allowed methods were set in the correct
				// header.
				allowedMethods := res.Header.Get("Access-Control-Allow-Methods")
				for method := range sut.methodHandlers {
					if err := assert.True(
						strings.Contains(allowedMethods, method),
					); err != nil {
						t.Error(err)
					}
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
		if err = assert.Equal(
			http.StatusUnauthorized, res.StatusCode,
		); err != nil {
			t.Error(err)
		}

		// Assert the WWW-Authenticate header was set correctly
		if err := assert.Equal(
			"Bearer", w.Result().Header.Get("WWW-Authenticate"),
		); err != nil {
			t.Error(err)
		}
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
				w := httptest.NewRecorder()

				// Handle request with sut and get the result.
				sut.ServeHTTP(w, req)
				res := w.Result()

				// Assert on the status code.
				if err = assert.Equal(
					http.StatusOK, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				fakeMethodHandler := methodHandler.(*FakeMethodHandler)

				// Run case-specific assertions.
				if err := assert.Equal(
					w, fakeMethodHandler.InResponseWriter,
				); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(
					req, fakeMethodHandler.InReq,
				); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(
					wantSub, fakeMethodHandler.InSub,
				); err != nil {
					t.Error(err)
				}
			})
		}
	})
}
