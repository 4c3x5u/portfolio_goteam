//go:build utest

package column

import (
	"bytes"
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
	userboardTable "server/dbaccess/userboard"
	pkgLog "server/log"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestHandler(t *testing.T) {
	authHeaderReader := &auth.FakeHeaderReader{}
	authTokenValidator := &auth.FakeTokenValidator{}
	idValidator := &api.FakeStringValidator{}
	columnSelector := &columnTable.FakeSelector{}
	userBoardSelector := &userboardTable.FakeSelector{}
	columnUpdater := &columnTable.FakeUpdater{}
	log := &pkgLog.FakeErrorer{}
	sut := NewHandler(
		authHeaderReader,
		authTokenValidator,
		idValidator,
		columnSelector,
		userBoardSelector,
		columnUpdater,
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

	assertOnLoggedErr := func(wantErrMsg string) func(*testing.T, *http.Response) {
		return func(t *testing.T, _ *http.Response) {
			if err := assert.Equal(
				wantErrMsg, log.InMessage,
			); err != nil {
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
		name                        string
		authTokenValidatorOutSub    string
		idValidatorOutErr           error
		columnSelectorOutErr        error
		userBoardSelectorOutIsAdmin bool
		userBoardSelectorOutErr     error
		columnUpdaterOutErr         error
		wantStatusCode              int
		assertFunc                  func(t *testing.T, r *http.Response)
	}{
		{
			name:                        "InvalidAuthToken",
			authTokenValidatorOutSub:    "",
			idValidatorOutErr:           nil,
			columnSelectorOutErr:        nil,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     nil,
			columnUpdaterOutErr:         nil,
			wantStatusCode:              http.StatusUnauthorized,
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
			name:                        "IDValidatorErr",
			authTokenValidatorOutSub:    "bob123",
			idValidatorOutErr:           errors.New("invalid id"),
			columnSelectorOutErr:        nil,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     nil,
			columnUpdaterOutErr:         nil,
			wantStatusCode:              http.StatusBadRequest,
			assertFunc:                  assertOnResErr("invalid id"),
		},
		{
			name:                        "ColumnNotFound",
			authTokenValidatorOutSub:    "bob123",
			idValidatorOutErr:           nil,
			columnSelectorOutErr:        sql.ErrNoRows,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     nil,
			columnUpdaterOutErr:         nil,
			wantStatusCode:              http.StatusBadRequest,
			assertFunc:                  assertOnResErr("Column not found."),
		},
		{
			name:                        "ColumnSelectorErr",
			authTokenValidatorOutSub:    "bob123",
			idValidatorOutErr:           nil,
			columnSelectorOutErr:        sql.ErrConnDone,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     nil,
			columnUpdaterOutErr:         nil,
			wantStatusCode:              http.StatusInternalServerError,
			assertFunc: assertOnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                        "UserBoardNotFound",
			authTokenValidatorOutSub:    "bob123",
			idValidatorOutErr:           nil,
			columnSelectorOutErr:        nil,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     sql.ErrNoRows,
			columnUpdaterOutErr:         nil,
			wantStatusCode:              http.StatusUnauthorized,
			assertFunc: assertOnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:                        "UserBoardSelectorErr",
			authTokenValidatorOutSub:    "bob123",
			idValidatorOutErr:           nil,
			columnSelectorOutErr:        nil,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     sql.ErrConnDone,
			columnUpdaterOutErr:         nil,
			wantStatusCode:              http.StatusInternalServerError,
			assertFunc: assertOnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                        "UserIsNotAdmin",
			authTokenValidatorOutSub:    "bob123",
			idValidatorOutErr:           nil,
			columnSelectorOutErr:        nil,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     nil,
			columnUpdaterOutErr:         nil,
			wantStatusCode:              http.StatusUnauthorized,
			assertFunc: assertOnResErr(
				"Only board admins can move tasks.",
			),
		},
		{
			name:                        "TaskNotFound",
			authTokenValidatorOutSub:    "bob123",
			idValidatorOutErr:           nil,
			columnSelectorOutErr:        nil,
			userBoardSelectorOutIsAdmin: true,
			userBoardSelectorOutErr:     nil,
			columnUpdaterOutErr:         sql.ErrNoRows,
			wantStatusCode:              http.StatusNotFound,
			assertFunc:                  assertOnResErr("Task not found."),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for sut's dependencies.
			authTokenValidator.OutSub = c.authTokenValidatorOutSub
			idValidator.OutErr = c.idValidatorOutErr
			columnSelector.OutErr = c.columnSelectorOutErr
			userBoardSelector.OutIsAdmin = c.userBoardSelectorOutIsAdmin
			userBoardSelector.OutErr = c.userBoardSelectorOutErr
			columnUpdater.OutErr = c.columnUpdaterOutErr

			// Prepare request and response recorder.
			tasks, err := json.Marshal([]map[string]int{{"id": 0, "order": 0}})
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPatch, "", bytes.NewReader(tasks),
			)
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
