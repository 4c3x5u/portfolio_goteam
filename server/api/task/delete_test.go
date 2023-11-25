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
	userSelector := &userTable.FakeSelector{}
	idValidator := &api.FakeStringValidator{}
	taskSelector := &taskTable.FakeSelector{}
	columnSelector := &columnTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	taskDeleter := &dbaccess.FakeDeleter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDELETEHandler(
		userSelector,
		idValidator,
		taskSelector,
		columnSelector,
		boardSelector,
		taskDeleter,
		log,
	)

	for _, c := range []struct {
		user              userTable.Record
		userSelectorErr   error
		name              string
		idValidatorErr    error
		taskSelectorErr   error
		columnSelectorErr error
		board             boardTable.Record
		boardSelectorErr  error
		taskDeleterErr    error
		wantStatusCode    int
		assertFunc        func(*testing.T, *http.Response, string)
	}{
		{
			name:              "UserSelectorErr",
			user:              userTable.Record{},
			userSelectorErr:   sql.ErrConnDone,
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:              "UserNotRecognised",
			user:              userTable.Record{},
			userSelectorErr:   sql.ErrNoRows,
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusUnauthorized,
			assertFunc:        assert.OnResErr("Username is not recognised."),
		},
		{
			name:              "NotAdmin",
			user:              userTable.Record{IsAdmin: false},
			userSelectorErr:   nil,
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only board admins can delete tasks.",
			),
		},
		{
			name:              "IDEmpty",
			user:              userTable.Record{IsAdmin: true},
			userSelectorErr:   nil,
			idValidatorErr:    api.ErrStrEmpty,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusBadRequest,
			assertFunc:        assert.OnResErr("Task ID cannot be empty."),
		},
		{
			name:              "IDNotInt",
			user:              userTable.Record{IsAdmin: true},
			userSelectorErr:   nil,
			idValidatorErr:    api.ErrStrNotInt,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusBadRequest,
			assertFunc:        assert.OnResErr("Task ID must be an integer."),
		},
		{
			name:              "IDUnexpectedErr",
			user:              userTable.Record{IsAdmin: true},
			userSelectorErr:   nil,
			idValidatorErr:    api.ErrStrTooLong,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(api.ErrStrTooLong.Error()),
		},
		{
			name:              "TaskSelectorErr",
			user:              userTable.Record{IsAdmin: true},
			userSelectorErr:   nil,
			idValidatorErr:    nil,
			taskSelectorErr:   sql.ErrConnDone,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:              "TaskNotFound",
			user:              userTable.Record{IsAdmin: true},
			userSelectorErr:   nil,
			idValidatorErr:    nil,
			taskSelectorErr:   sql.ErrNoRows,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusNotFound,
			assertFunc:        assert.OnResErr("Task not found."),
		},
		{
			name:              "ColumnSelectorErr",
			user:              userTable.Record{IsAdmin: true},
			userSelectorErr:   nil,
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: sql.ErrNoRows,
			board:             boardTable.Record{},
			boardSelectorErr:  nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:              "BoardSelectorErr",
			user:              userTable.Record{IsAdmin: true},
			userSelectorErr:   nil,
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{},
			boardSelectorErr:  sql.ErrNoRows,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:              "BoardWrongTeam",
			user:              userTable.Record{IsAdmin: true, TeamID: 2},
			userSelectorErr:   nil,
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{TeamID: 1},
			boardSelectorErr:  nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:              "TaskDeleterErr",
			user:              userTable.Record{IsAdmin: true, TeamID: 1},
			userSelectorErr:   nil,
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{TeamID: 1},
			boardSelectorErr:  nil,
			taskDeleterErr:    sql.ErrNoRows,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:              "Success",
			user:              userTable.Record{IsAdmin: true, TeamID: 1},
			userSelectorErr:   nil,
			idValidatorErr:    nil,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			board:             boardTable.Record{TeamID: 1},
			boardSelectorErr:  nil,
			taskDeleterErr:    nil,
			wantStatusCode:    http.StatusOK,
			assertFunc:        func(*testing.T, *http.Response, string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userSelector.User = c.user
			userSelector.Err = c.userSelectorErr
			idValidator.Err = c.idValidatorErr
			taskSelector.Err = c.taskSelectorErr
			columnSelector.Err = c.columnSelectorErr
			boardSelector.Board = c.board
			boardSelector.Err = c.boardSelectorErr
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
