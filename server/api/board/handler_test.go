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
		authHeaderReader, authTokenValidator, map[string]api.MethodHandler{
			http.MethodPost:   postHandler,
			http.MethodDelete: deleteHandler,
		},
	)

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
		wantMethodHandler    *api.FakeMethodHandler
		wantStatusCode       int
	}{
		{
			name:                 "InvalidAuthToken",
			httpMethod:           http.MethodPost,
			tokenValidatorOutSub: "",
			wantMethodHandler:    nil,
			wantStatusCode:       http.StatusUnauthorized,
		},
		{
			name:                 "Success" + http.MethodPost,
			httpMethod:           http.MethodPost,
			tokenValidatorOutSub: "bob123",
			wantMethodHandler:    postHandler,
			wantStatusCode:       http.StatusOK,
		},
		{
			name:                 "Success" + http.MethodDelete,
			httpMethod:           http.MethodDelete,
			tokenValidatorOutSub: "bob123",
			wantMethodHandler:    deleteHandler,
			wantStatusCode:       http.StatusOK,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			authTokenValidator.OutSub = c.tokenValidatorOutSub

			req, err := http.NewRequest(c.httpMethod, "/board", nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()

			sut.ServeHTTP(w, req)

			if err = assert.Equal(
				c.wantStatusCode, w.Result().StatusCode,
			); err != nil {
				t.Error(err)
			}

			if c.wantStatusCode == http.StatusUnauthorized {
				// 401 is expected - WWWAuthenticated cookie must be set.
				name, value := auth.WWWAuthenticate()
				if err = assert.Equal(
					value, w.Result().Header.Get(name),
				); err != nil {
					t.Error(err)
				}
			} else {
				// 401 is NOT expected - a method handler must be called.
				if err := assert.Equal(
					w, c.wantMethodHandler.InResponseWriter,
				); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(
					req, c.wantMethodHandler.InReq,
				); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(
					c.tokenValidatorOutSub, c.wantMethodHandler.InSub,
				); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
