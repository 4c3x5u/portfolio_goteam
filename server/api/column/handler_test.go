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

	for _, c := range []struct {
		name                     string
		authTokenValidatorOutSub string
		idValidatorOutErr        error
		wantStatusCode           int
		assertFunc               func(t *testing.T, r *http.Response)
	}{
		{
			name:                     "InvalidAuthToken",
			authTokenValidatorOutSub: "",
			idValidatorOutErr:        nil,
			wantStatusCode:           http.StatusUnauthorized,
			assertFunc: func(t *testing.T, res *http.Response) {
				name, value := auth.WWWAuthenticate()
				if err := assert.Equal(
					value, res.Header.Get(name),
				); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:                     "IDValidatorErr",
			authTokenValidatorOutSub: "bob123",
			idValidatorOutErr:        errors.New("invalid id"),
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: func(t *testing.T, res *http.Response) {
				var resBody ResBody
				if err := json.NewDecoder(res.Body).Decode(
					&resBody,
				); err != nil {
					t.Fatal(err)
				}
				if err := assert.Equal(
					"invalid id", resBody.Error,
				); err != nil {
					t.Error(err)
				}
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			authTokenValidator.OutSub = c.authTokenValidatorOutSub
			idValidator.OutErr = c.idValidatorOutErr

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
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			// Run case-specific assertions.
			c.assertFunc(t, res)
		})
	}
}
