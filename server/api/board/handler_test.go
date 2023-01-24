package board

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/api"
	"server/assert"
	"server/auth"
)

func TestHandler(t *testing.T) {
	authHeaderReader := &auth.FakeHeaderReader{}
	authTokenValidator := &auth.FakeTokenValidator{}
	postHandler := &api.FakeMethodHandler{}
	sut := NewHandler(authHeaderReader, authTokenValidator, postHandler)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		for _, httpMethod := range []string{
			http.MethodConnect, http.MethodDelete, http.MethodGet,
			http.MethodHead, http.MethodOptions, http.MethodPatch,
			http.MethodPut, http.MethodTrace,
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
			})
		}
	})

	for _, c := range []struct {
		name                     string
		reqBody                  ReqBody
		authHeaderReaderOutToken string
		tokenValidatorOutSub     string
		wantStatusCode           int
	}{
		{
			name:                     "AuthTokenEmpty",
			reqBody:                  ReqBody{},
			authHeaderReaderOutToken: "",
			tokenValidatorOutSub:     "",
			wantStatusCode:           http.StatusUnauthorized,
		},
		{
			name:                     "TokenValidatorErr",
			reqBody:                  ReqBody{},
			authHeaderReaderOutToken: "ABCDEFG",
			tokenValidatorOutSub:     "",
			wantStatusCode:           http.StatusUnauthorized,
		},
		{
			name:                     "Success",
			reqBody:                  ReqBody{Name: "someboard"},
			authHeaderReaderOutToken: "ABCDEFG",
			tokenValidatorOutSub:     "bob123",
			wantStatusCode:           http.StatusOK,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			authHeaderReader.OutToken = c.authHeaderReaderOutToken
			authTokenValidator.OutSub = c.tokenValidatorOutSub

			reqBodyJSON, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPost, "/board", bytes.NewReader(reqBodyJSON),
			)
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

			// If 401 is expected, WWWAuthenticated cookie must be set.
			if c.wantStatusCode == http.StatusUnauthorized {
				name, value := auth.WWWAuthenticate()
				if err = assert.Equal(value, w.Result().Header.Get(name)); err != nil {
					t.Error(err)
				}
			}

			// DEPENDENCY-INPUT-BASED ASSERTIONS

			// If 200 is expected, postHandler must be called.
			if c.wantStatusCode == http.StatusOK {
				if err := assert.Equal(w, postHandler.InResponseWriter); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(req, postHandler.InReq); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(
					c.tokenValidatorOutSub, postHandler.InSub,
				); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
