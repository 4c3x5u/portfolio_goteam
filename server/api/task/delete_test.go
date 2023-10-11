//go:build utest

package task

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/api"
	"server/assert"
	"server/dbaccess"
	columnTable "server/dbaccess/column"
	taskTable "server/dbaccess/task"
	userboardTable "server/dbaccess/userboard"
	pkgLog "server/log"
)

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDELETEHandler(t *testing.T) {
	idValidator := &api.FakeStringValidator{}
	taskSelector := &taskTable.FakeSelector{}
	columnSelector := &columnTable.FakeSelector{}
	userBoardSelector := &userboardTable.FakeSelector{}
	taskDeleter := &dbaccess.FakeDeleter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDELETEHandler(
		idValidator,
		taskSelector,
		columnSelector,
		userBoardSelector,
		taskDeleter,
		log,
	)

	for _, c := range []struct {
		name                 string
		idValidatorErr       error
		taskSelectorErr      error
		columnSelectorErr    error
		userIsAdmin          bool
		userBoardSelectorErr error
		taskDeleterErr       error
		wantStatusCode       int
		assertFunc           func(*testing.T, *http.Response, string)
	}{
		{
			name:                 "IDEmpty",
			idValidatorErr:       api.ErrStrEmpty,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			taskDeleterErr:       nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc:           assert.OnResErr("Task ID cannot be empty."),
		},
		{
			name:                 "IDNotInt",
			idValidatorErr:       api.ErrStrNotInt,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			taskDeleterErr:       nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc:           assert.OnResErr("Task ID must be an integer."),
		},
		{
			name:                 "IDUnexpectedErr",
			idValidatorErr:       api.ErrStrTooLong,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			taskDeleterErr:       nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr(api.ErrStrTooLong.Error()),
		},
		{
			name:                 "TaskSelectorErr",
			idValidatorErr:       nil,
			taskSelectorErr:      sql.ErrConnDone,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			taskDeleterErr:       nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:                 "TaskNotFound",
			idValidatorErr:       nil,
			taskSelectorErr:      sql.ErrNoRows,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			taskDeleterErr:       nil,
			wantStatusCode:       http.StatusNotFound,
			assertFunc:           assert.OnResErr("Task not found."),
		},
		{
			name:                 "ColumnSelectorErr",
			idValidatorErr:       nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    sql.ErrNoRows,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			taskDeleterErr:       nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:                 "UserBoardSelectorErr",
			idValidatorErr:       nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: sql.ErrConnDone,
			taskDeleterErr:       nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:                 "NoAccess",
			idValidatorErr:       nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: sql.ErrNoRows,
			taskDeleterErr:       nil,
			wantStatusCode:       http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:                 "NotAdmin",
			idValidatorErr:       nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			taskDeleterErr:       nil,
			wantStatusCode:       http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only board admins can delete tasks.",
			),
		},
		{
			name:                 "TaskDeleterErr",
			idValidatorErr:       nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          true,
			userBoardSelectorErr: nil,
			taskDeleterErr:       sql.ErrNoRows,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:                 "Success",
			idValidatorErr:       nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          true,
			userBoardSelectorErr: nil,
			taskDeleterErr:       nil,
			wantStatusCode:       http.StatusOK,
			assertFunc:           func(*testing.T, *http.Response, string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			idValidator.Err = c.idValidatorErr
			taskSelector.Err = c.taskSelectorErr
			columnSelector.Err = c.columnSelectorErr
			userBoardSelector.IsAdmin = c.userIsAdmin
			userBoardSelector.Err = c.userBoardSelectorErr
			taskDeleter.Err = c.taskDeleterErr

			r, err := http.NewRequest("", "", nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			sut.Handle(w, r, "")
			res := w.Result()

			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			c.assertFunc(t, res, log.InMessage)
		})
	}
}
