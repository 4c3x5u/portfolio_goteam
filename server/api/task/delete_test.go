//go:build utest

package task

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	taskTable "github.com/kxplxn/goteam/server/dbaccess/task"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDELETEHandler(t *testing.T) {
	idValidator := &api.FakeStringValidator{}
	taskSelector := &taskTable.FakeSelector{}
	columnSelector := &columnTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	userSelector := &userTable.FakeSelector{}
	taskDeleter := &dbaccess.FakeDeleter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDELETEHandler(
		idValidator,
		taskSelector,
		columnSelector,
		boardSelector,
		userSelector,
		taskDeleter,
		log,
	)

	for _, c := range []struct {
		name              string
		idValidatorErr    error
		taskSelectorErr   error
		columnSelectorErr error
		board             boardTable.Record
		boardSelectorErr  error
		user              userTable.Record
		userSelectorErr   error
		taskDeleterErr    error
		wantStatusCode    int
		assertFunc        func(*testing.T, *http.Response, string)
	}{
		{
			name:              "IDEmpty",
			idValidatorErr:    api.ErrStrEmpty,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			user:              userTable.Record{},
			userSelectorErr:   nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusBadRequest,
			assertFunc:        assert.OnResErr("Task ID cannot be empty."),
		},
		{
			name:              "IDNotInt",
			idValidatorErr:    api.ErrStrNotInt,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			user:              userTable.Record{},
			userSelectorErr:   nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusBadRequest,
			assertFunc:        assert.OnResErr("Task ID must be an integer."),
		},
		{
			name:              "IDUnexpectedErr",
			idValidatorErr:    api.ErrStrTooLong,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			user:              userTable.Record{},
			userSelectorErr:   nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(api.ErrStrTooLong.Error()),
		},
		{
			name:              "TaskSelectorErr",
			idValidatorErr:    nil,
			taskSelectorErr:   sql.ErrConnDone,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			user:              userTable.Record{},
			userSelectorErr:   nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:              "TaskNotFound",
			idValidatorErr:    nil,
			taskSelectorErr:   sql.ErrNoRows,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			user:              userTable.Record{},
			userSelectorErr:   nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusNotFound,
			assertFunc:        assert.OnResErr("Task not found."),
		},
		{
			name:              "ColumnSelectorErr",
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: sql.ErrNoRows,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			user:              userTable.Record{},
			userSelectorErr:   nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:              "BoardSelectorErr",
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  sql.ErrNoRows,
			user:              userTable.Record{},
			userSelectorErr:   nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:              "UserSelectorErr",
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			user:              userTable.Record{},
			userSelectorErr:   sql.ErrConnDone,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:              "UserNotRecognised",
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			user:              userTable.Record{},
			userSelectorErr:   sql.ErrNoRows,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusUnauthorized,
			assertFunc:        assert.OnResErr("Username is not recognised."),
		},
		{
			name:              "NotAdmin",
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			user:              userTable.Record{IsAdmin: false},
			userSelectorErr:   nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only board admins can delete tasks.",
			),
		},
		{
			name:              "BoardWrongTeam",
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{TeamID: 1},
			boardSelectorErr:  nil,
			user:              userTable.Record{IsAdmin: true, TeamID: 2},
			userSelectorErr:   nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:              "TaskDeleterErr",
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{TeamID: 1},
			boardSelectorErr:  nil,
			user:              userTable.Record{IsAdmin: true, TeamID: 1},
			userSelectorErr:   nil,
			taskDeleterErr:    sql.ErrNoRows,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:              "Success",
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{TeamID: 1},
			boardSelectorErr:  nil,
			user:              userTable.Record{IsAdmin: true, TeamID: 1},
			userSelectorErr:   nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusOK,
			assertFunc:        func(*testing.T, *http.Response, string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			idValidator.Err = c.idValidatorErr
			taskSelector.Err = c.taskSelectorErr
			columnSelector.Err = c.columnSelectorErr
			boardSelector.Board = c.board
			boardSelector.Err = c.boardSelectorErr
			userSelector.User = c.user
			userSelector.Err = c.userSelectorErr
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
