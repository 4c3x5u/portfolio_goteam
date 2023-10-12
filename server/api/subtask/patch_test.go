//go:build utest

package subtask

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	subtaskTable "github.com/kxplxn/goteam/server/dbaccess/subtask"
	taskTable "github.com/kxplxn/goteam/server/dbaccess/task"
	userboardTable "github.com/kxplxn/goteam/server/dbaccess/userboard"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestPATCHHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	idValidator := &api.FakeStringValidator{}
	subtaskSelector := &subtaskTable.FakeSelector{}
	taskSelector := &taskTable.FakeSelector{}
	columnSelector := &columnTable.FakeSelector{}
	userBoardSelector := &userboardTable.FakeSelector{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(
		idValidator,
		subtaskSelector,
		taskSelector,
		columnSelector,
		userBoardSelector,
		log,
	)

	for _, c := range []struct {
		name                 string
		idValidatorErr       error
		subtaskSelectorErr   error
		taskSelectorErr      error
		columnSelectorErr    error
		userIsAdmin          bool
		userBoardSelectorErr error
		wantStatusCode       int
		assertFunc           func(*testing.T, *http.Response, string)
	}{
		{
			name:                 "IDEmpty",
			idValidatorErr:       api.ErrStrEmpty,
			subtaskSelectorErr:   nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc:           assert.OnResErr("Subtask ID cannot be empty."),
		},
		{
			name:                 "IDNotInt",
			idValidatorErr:       api.ErrStrNotInt,
			subtaskSelectorErr:   nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc:           assert.OnResErr("Subtask ID must be an integer."),
		},
		{
			name:                 "IDUnexpectedErr",
			idValidatorErr:       api.ErrStrTooLong,
			subtaskSelectorErr:   nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr(api.ErrStrTooLong.Error()),
		},
		{
			name:                 "SubtaskSelectorErr",
			idValidatorErr:       nil,
			subtaskSelectorErr:   sql.ErrConnDone,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:                 "SubtaskNotFound",
			idValidatorErr:       nil,
			subtaskSelectorErr:   sql.ErrNoRows,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			wantStatusCode:       http.StatusNotFound,
			assertFunc:           assert.OnResErr("Subtask not found."),
		},
		{
			name:                 "TaskSelectorErr",
			idValidatorErr:       nil,
			subtaskSelectorErr:   nil,
			taskSelectorErr:      sql.ErrNoRows,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:                 "ColumnSelectorErr",
			idValidatorErr:       nil,
			subtaskSelectorErr:   nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    sql.ErrNoRows,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:                 "UserBoardSelectorErr",
			idValidatorErr:       nil,
			subtaskSelectorErr:   nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: sql.ErrConnDone,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:                 "NoAccess",
			idValidatorErr:       nil,
			subtaskSelectorErr:   nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: sql.ErrNoRows,
			wantStatusCode:       http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:                 "NotAdmin",
			idValidatorErr:       nil,
			subtaskSelectorErr:   nil,
			taskSelectorErr:      nil,
			columnSelectorErr:    nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			wantStatusCode:       http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only board admins can edit subtasks.",
			),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			idValidator.Err = c.idValidatorErr
			subtaskSelector.Err = c.subtaskSelectorErr
			taskSelector.Err = c.taskSelectorErr
			columnSelector.Err = c.columnSelectorErr
			userBoardSelector.Err = c.userBoardSelectorErr

			r, err := http.NewRequest("", "?id=", nil)
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
