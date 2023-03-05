//go:build utest

package board

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"server/api"
	"server/assert"
	"server/auth"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestHandler(t *testing.T) {
	authHeaderReader := &auth.FakeHeaderReader{}
	authTokenValidator := &auth.FakeTokenValidator{}
	postHandler := &api.FakeMethodHandler{}
	deleteHandler := &api.FakeMethodHandler{}
	sut := NewHandler(
		authHeaderReader,
		authTokenValidator,
		map[string]api.MethodHandler{
			http.MethodPost:   postHandler,
			http.MethodDelete: deleteHandler,
		},
	)

	// Used in success cases to assert on the parameters received by the method
	// handler.
	assertOnMethodHandler := func(
		h *api.FakeMethodHandler, sub string,
	) func(*testing.T, *http.Request, *httptest.ResponseRecorder) {
		return func(
			t *testing.T, r *http.Request, w *httptest.ResponseRecorder,
		) {
			if err := assert.Equal(w, h.InResponseWriter); err != nil {
				t.Error(err)
			}
			if err := assert.Equal(r, h.InReq); err != nil {
				t.Error(err)
			}
			if err := assert.Equal(sub, h.InSub); err != nil {
				t.Error(err)
			}
		}
	}

	t.Run("MethodNotAllowed", func(t *testing.T) {
		for _, httpMethod := range []string{
			http.MethodConnect, http.MethodGet, http.MethodHead,
			http.MethodOptions, http.MethodPatch, http.MethodPut,
			http.MethodTrace,
		} {
			t.Run(httpMethod, func(t *testing.T) {
				req, err := http.NewRequest(httpMethod, "/board", nil)
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
		name                 string
		httpMethod           string
		tokenValidatorOutSub string
		wantStatusCode       int
		assertFunc           func(
			*testing.T, *http.Request, *httptest.ResponseRecorder,
		)
	}{
		{
			name:                 "InvalidAuthToken",
			httpMethod:           http.MethodPost,
			tokenValidatorOutSub: "",
			wantStatusCode:       http.StatusUnauthorized,
			assertFunc: func(
				t *testing.T, _ *http.Request, w *httptest.ResponseRecorder,
			) {
				name, value := auth.WWWAuthenticate()
				if err := assert.Equal(
					value, w.Result().Header.Get(name),
				); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:                 "Success" + http.MethodPost,
			httpMethod:           http.MethodPost,
			tokenValidatorOutSub: "bob123",
			wantStatusCode:       http.StatusOK,
			assertFunc:           assertOnMethodHandler(postHandler, "bob123"),
		},
		{
			name:                 "Success" + http.MethodDelete,
			httpMethod:           http.MethodDelete,
			tokenValidatorOutSub: "bob123",
			wantStatusCode:       http.StatusOK,
			assertFunc: assertOnMethodHandler(
				deleteHandler, "bob123",
			),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for sut's dependencies.
			authTokenValidator.OutSub = c.tokenValidatorOutSub

			// Prepare request and response recorder.
			req, err := http.NewRequest(c.httpMethod, "", nil)
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
			c.assertFunc(t, req, w)
		})
	}
}
