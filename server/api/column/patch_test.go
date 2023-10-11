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

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	userboardTable "github.com/kxplxn/goteam/server/dbaccess/userboard"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestHandler(t *testing.T) {
	idValidator := &api.FakeStringValidator{}
	columnSelector := &columnTable.FakeSelector{}
	userBoardSelector := &userboardTable.FakeSelector{}
	columnUpdater := &columnTable.FakeUpdater{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(
		idValidator,
		columnSelector,
		userBoardSelector,
		columnUpdater,
		log,
	)

	for _, c := range []struct {
		name                 string
		idValidatorErr       error
		columnSelectorErr    error
		userIsAdmin          bool
		userBoardSelectorErr error
		columnUpdaterErr     error
		wantStatusCode       int
		assertFunc           func(*testing.T, *http.Response, string)
	}{
		{
			name:                 "IDValidatorErr",
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
			sut.Handle(w, req, "")
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
