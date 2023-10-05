//go:build utest

package column

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/api"
	"server/assert"
	"server/auth"
	columnTable "server/dbaccess/column"
	pkgLog "server/log"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestHandler(t *testing.T) {
	authHeaderReader := &auth.FakeHeaderReader{}
	authTokenValidator := &auth.FakeTokenValidator{}
	idValidator := &api.FakeStringValidator{}
	columnSelector := &columnTable.FakeSelector{}
	log := &pkgLog.FakeErrorer{}
	sut := NewHandler(
		authHeaderReader,
		authTokenValidator,
		idValidator,
		columnSelector,
		log,
	)

	// Used in status 400 cases to assert on the error returned in res body.
	assertOnResErr := func(
		wantErrMsg string,
	) func(*testing.T, *http.Response) {
		return func(
			t *testing.T, res *http.Response,
		) {
			resBody := ResBody{}
			if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
				t.Fatal(err)
			}
			if err := assert.Equal(wantErrMsg, resBody.Error); err != nil {
				t.Error(err)
			}
		}
	}

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
		columnSelectorOutErr     error
		wantStatusCode           int
		assertFunc               func(t *testing.T, r *http.Response)
	}{
		{
			name:                     "InvalidAuthToken",
			authTokenValidatorOutSub: "",
			idValidatorOutErr:        nil,
			columnSelectorOutErr:     nil,
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
			columnSelectorOutErr:     nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc:               assertOnResErr("invalid id"),
		},
		{
			name:                     "ColumnNotFound",
			authTokenValidatorOutSub: "bob123",
			idValidatorOutErr:        nil,
			columnSelectorOutErr:     sql.ErrNoRows,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc:               assertOnResErr("Column not found."),
		},
		{
			name:                     "ColumnSelectorErr",
			authTokenValidatorOutSub: "bob123",
			idValidatorOutErr:        nil,
			columnSelectorOutErr:     sql.ErrConnDone,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: func(t *testing.T, _ *http.Response) {
				if err := assert.Equal(
					sql.ErrConnDone.Error(), log.InMessage,
				); err != nil {
					t.Error(err)
				}
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for sut's dependencies.
			authTokenValidator.OutSub = c.authTokenValidatorOutSub
			idValidator.OutErr = c.idValidatorOutErr
			columnSelector.OutErr = c.columnSelectorOutErr

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
