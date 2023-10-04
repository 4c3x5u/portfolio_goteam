//go:build utest

package column

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/auth"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestHandler(t *testing.T) {
	authHeaderReader := &auth.FakeHeaderReader{}
	authTokenValidator := &auth.FakeTokenValidator{}
	sut := NewHandler(authHeaderReader, authTokenValidator)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		for _, httpMethod := range []string{
			http.MethodConnect, http.MethodGet, http.MethodPost,
			http.MethodDelete, http.MethodHead, http.MethodOptions,
			http.MethodPut, http.MethodTrace,
		} {
			t.Run(httpMethod, func(t *testing.T) {
				req, err := http.NewRequest(httpMethod, "", nil)
				if err != nil {
					t.Fatal(err)
				}
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)

				if err = assert.Equal(
					http.StatusMethodNotAllowed, w.Result().StatusCode,
				); err != nil {
					t.Error(err)
				}

				if err := assert.Equal(
					http.MethodPost,
					w.Result().Header.Get("Access-Control-Allow-Methods"),
				); err != nil {
					t.Error(err)
				}
			})
		}
	})

	t.Run("InvalidAuthToken", func(t *testing.T) {
		// Set pre-determinate return values for sut's dependencies.
		authTokenValidator.OutSub = ""

		// Prepare request and response recorder.
		req, err := http.NewRequest(http.MethodPatch, "", nil)
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

		// Run case-specific assertions.
		name, value := auth.WWWAuthenticate()
		if err := assert.Equal(
			value, w.Result().Header.Get(name),
		); err != nil {
			t.Error(err)
		}
	})
}
