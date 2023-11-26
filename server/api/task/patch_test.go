//go:build utest

package task

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	taskTable "github.com/kxplxn/goteam/server/dbaccess/task"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestPATCHHandler tests the Handle method of PATCHHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	userSelector := &userTable.FakeSelector{}
	taskIDValidator := &api.FakeStringValidator{}
	taskTitleValidator := &api.FakeStringValidator{}
	subtaskTitleValidator := &api.FakeStringValidator{}
	taskSelector := &taskTable.FakeSelector{}
	columnSelector := &columnTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	taskUpdater := &taskTable.FakeUpdater{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(
		userSelector,
		taskIDValidator,
		taskTitleValidator,
		subtaskTitleValidator,
		taskSelector,
		columnSelector,
		boardSelector,
		taskUpdater,
		log,
	)

	for _, c := range []struct {
		name                     string
		user                     userTable.Record
		userSelectorErr          error
		taskIDValidatorErr       error
		taskTitleValidatorErr    error
		subtaskTitleValidatorErr error
		taskSelectorErr          error
		columnSelectorErr        error
		board                    boardTable.Record
		boardSelectorErr         error
		taskUpdaterErr           error
		wantStatusCode           int
		assertFunc               func(*testing.T, *http.Response, string)
	}{
		{
			name:                     "UserNotRecognised",
			user:                     userTable.Record{},
			userSelectorErr:          sql.ErrNoRows,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusUnauthorized,
			assertFunc: assert.OnResErr(
				"Username is not recognised.",
			),
		},
		{
			name:                     "UserSelectorErr",
			user:                     userTable.Record{},
			userSelectorErr:          sql.ErrConnDone,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                     "NotAdmin",
			user:                     userTable.Record{IsAdmin: false},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can edit tasks.",
			),
		},
		{
			name:                     "TaskIDEmpty",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       api.ErrStrEmpty,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task ID cannot be empty.",
			),
		},
		{
			name:                     "TaskIDNotInt",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       api.ErrStrNotInt,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task ID must be an integer.",
			),
		},
		{
			name:                     "TaskIDUnexpectedErr",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       api.ErrStrTooLong,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				api.ErrStrTooLong.Error(),
			),
		},
		{
			name:                     "TaskTitleEmpty",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    api.ErrStrEmpty,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be empty.",
			),
		},
		{
			name:                     "TaskTitleTooLong",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    api.ErrStrTooLong,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:                     "TaskTitleUnexpectedErr",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    api.ErrStrNotInt,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				api.ErrStrNotInt.Error(),
			),
		},
		{
			name:                     "SubtaskTitleEmpty",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrEmpty,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be empty.",
			),
		},
		{
			name:                     "SubtaskTitleTooLong",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrTooLong,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be longer than 50 characters.",
			),
		},
		{
			name:                     "SubtaskTitleUnexpectedErr",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrNotInt,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				api.ErrStrNotInt.Error(),
			),
		},
		{
			name:                     "TaskNotFound",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          sql.ErrNoRows,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusNotFound,
			assertFunc:               assert.OnResErr("Task not found."),
		},
		{
			name:                     "TaskSelectorErr",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          sql.ErrConnDone,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                     "ColumnSelectorErr",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        sql.ErrNoRows,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc:               assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:                     "BoardSelectorErr",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         sql.ErrNoRows,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc:               assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name: "NoAccess",
			user: userTable.Record{
				IsAdmin: true, TeamID: 2,
			},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{TeamID: 1},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name: "TaskUpdaterErr",
			user: userTable.Record{
				IsAdmin: true, TeamID: 1,
			},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{TeamID: 1},
			boardSelectorErr:         nil,
			taskUpdaterErr:           sql.ErrConnDone,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name: "Success",
			user: userTable.Record{
				IsAdmin: true, TeamID: 1,
			},
			userSelectorErr:          nil,
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			taskSelectorErr:          nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{TeamID: 1},
			boardSelectorErr:         nil,
			taskUpdaterErr:           nil,
			wantStatusCode:           http.StatusOK,
			assertFunc:               func(*testing.T, *http.Response, string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userSelector.Rec = c.user
			userSelector.Err = c.userSelectorErr
			taskIDValidator.Err = c.taskIDValidatorErr
			taskTitleValidator.Err = c.taskTitleValidatorErr
			subtaskTitleValidator.Err = c.subtaskTitleValidatorErr
			taskSelector.Err = c.taskSelectorErr
			columnSelector.Err = c.columnSelectorErr
			boardSelector.Board = c.board
			boardSelector.Err = c.boardSelectorErr
			taskUpdater.Err = c.taskUpdaterErr

			reqBody, err := json.Marshal(map[string]any{
				"column":      0,
				"title":       "",
				"description": "",
				"subtasks":    []map[string]any{{"title": ""}},
			})
			if err != nil {
				t.Fatal(err)
			}
			r, err := http.NewRequest("", "", bytes.NewReader(reqBody))
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
