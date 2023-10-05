//go:build utest

package column

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/api"
	"server/assert"
	"server/auth"
	pkgLog "server/log"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestHandler(t *testing.T) {
	authHeaderReader := &auth.FakeHeaderReader{}
	authTokenValidator := &auth.FakeTokenValidator{}
	idValidator := &api.FakeStringValidator{}
	log := &pkgLog.FakeErrorer{}
	sut := NewHandler(authHeaderReader, authTokenValidator, idValidator, log)

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

	t.Run("IDValidatorErr", func(t *testing.T) {
		// Set pre-determinate return values for sut's dependencies.
		const wantErrMsg = "Invalid ID."
		authTokenValidator.OutSub = "bob123"
		idValidator.OutErr = errors.New(wantErrMsg)

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
			http.StatusBadRequest, res.StatusCode,
		); err != nil {
			t.Error(err)
		}

		// Assert on error message.
		var resBody ResBody
		if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Fatal(err)
		}
		if err = assert.Equal(wantErrMsg, resBody.Error); err != nil {
			t.Error(err)
		}
	})
}
