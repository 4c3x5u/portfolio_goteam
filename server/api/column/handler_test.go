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
		name                 string
		sub                  string
		idValidatorErr       error
		columnSelectorErr    error
		userIsAdmin          bool
		userBoardSelectorErr error
		columnUpdaterErr     error
		wantStatusCode       int
		assertFunc           func(*testing.T, *http.Response, string)
	}{
		{
			name:                 "InvalidAuthToken",
			sub:                  "",
			idValidatorErr:       nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			columnUpdaterErr:     nil,
			wantStatusCode:       http.StatusUnauthorized,
			assertFunc: func(t *testing.T, res *http.Response, _ string) {
				name, value := auth.WWWAuthenticate()
				if err := assert.Equal(
					value, res.Header.Get(name),
				); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:                 "IDValidatorErr",
			sub:                  "bob123",
			idValidatorErr:       errors.New("invalid id"),
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			columnUpdaterErr:     nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc:           assert.OnResErr("invalid id"),
		},
		{
			name:                 "ColumnNotFound",
			sub:                  "bob123",
			idValidatorErr:       nil,
			columnSelectorErr:    sql.ErrNoRows,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			columnUpdaterErr:     nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc:           assert.OnResErr("Column not found."),
		},
		{
			name:                 "ColumnSelectorErr",
			sub:                  "bob123",
			idValidatorErr:       nil,
			columnSelectorErr:    sql.ErrConnDone,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			columnUpdaterErr:     nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                 "UserBoardNotFound",
			sub:                  "bob123",
			idValidatorErr:       nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: sql.ErrNoRows,
			columnUpdaterErr:     nil,
			wantStatusCode:       http.StatusUnauthorized,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:                 "UserBoardSelectorErr",
			sub:                  "bob123",
			idValidatorErr:       nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: sql.ErrConnDone,
			columnUpdaterErr:     nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                 "NotAdmin",
			sub:                  "bob123",
			idValidatorErr:       nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			columnUpdaterErr:     nil,
			wantStatusCode:       http.StatusUnauthorized,
			assertFunc: assert.OnResErr(
				"Only board admins can move tasks.",
			),
		},
		{
			name:                 "TaskNotFound",
			sub:                  "bob123",
			idValidatorErr:       nil,
			columnSelectorErr:    nil,
			userIsAdmin:          true,
			userBoardSelectorErr: nil,
			columnUpdaterErr:     sql.ErrNoRows,
			wantStatusCode:       http.StatusNotFound,
			assertFunc:           assert.OnResErr("Task not found."),
		},
		{
			name:                 "ColumnUpdaterErr",
			sub:                  "bob123",
			idValidatorErr:       nil,
			columnSelectorErr:    nil,
			userIsAdmin:          true,
			userBoardSelectorErr: nil,
			columnUpdaterErr:     sql.ErrConnDone,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for sut's dependencies.
			authTokenValidator.Sub = c.sub
			idValidator.Err = c.idValidatorErr
			columnSelector.Err = c.columnSelectorErr
			userBoardSelector.IsAdmin = c.userIsAdmin
			userBoardSelector.Err = c.userBoardSelectorErr
			columnUpdater.Err = c.columnUpdaterErr

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
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
