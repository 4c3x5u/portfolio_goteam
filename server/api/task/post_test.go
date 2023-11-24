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

// TestPOSTHandler tests the Handle method of POSTHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPOSTHandler(t *testing.T) {
	taskTitleValidator := &api.FakeStringValidator{}
	subtaskTitleValidator := &api.FakeStringValidator{}
	columnSelector := &columnTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	userSelector := &userTable.FakeSelector{}
	taskInserter := &taskTable.FakeInserter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPOSTHandler(
		taskTitleValidator,
		subtaskTitleValidator,
		columnSelector,
		boardSelector,
		userSelector,
		taskInserter,
		log,
	)

	for _, c := range []struct {
		name                     string
		taskTitleValidatorErr    error
		subtaskTitleValidatorErr error
		columnSelectorErr        error
		board                    boardTable.Record
		boardSelectorErr         error
		user                     userTable.Record
		userSelectorErr          error
		taskInserterErr          error
		wantStatusCode           int
		assertFunc               func(*testing.T, *http.Response, string)
	}{
		{
			name:                     "TaskTitleEmpty",
			taskTitleValidatorErr:    api.ErrStrEmpty,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			user:                     userTable.Record{},
			userSelectorErr:          nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be empty.",
			),
		},
		{
			name:                     "TaskTitleTooLong",
			taskTitleValidatorErr:    api.ErrStrTooLong,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			user:                     userTable.Record{},
			userSelectorErr:          nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:                     "TaskTitleUnexpectedErr",
			taskTitleValidatorErr:    api.ErrStrNotInt,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			user:                     userTable.Record{},
			userSelectorErr:          nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				api.ErrStrNotInt.Error(),
			),
		},
		{
			name:                     "SubtaskTitleEmpty",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrEmpty,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			user:                     userTable.Record{},
			userSelectorErr:          nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be empty.",
			),
		},
		{
			name:                     "SubtaskTitleTooLong",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrTooLong,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			user:                     userTable.Record{},
			userSelectorErr:          nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be longer than 50 characters.",
			),
		},
		{
			name:                     "SubtaskTitleUnexpectedErr",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrNotInt,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			user:                     userTable.Record{},
			userSelectorErr:          nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				api.ErrStrNotInt.Error(),
			),
		},
		{
			name:                     "ColumnNotFound",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        sql.ErrNoRows,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			user:                     userTable.Record{},
			userSelectorErr:          nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusNotFound,
			assertFunc:               assert.OnResErr("Column not found."),
		},
		{
			name:                     "ColumnSelectorErr",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        sql.ErrConnDone,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			user:                     userTable.Record{},
			userSelectorErr:          nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                     "BoardSelectorErr",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         sql.ErrConnDone,
			user:                     userTable.Record{},
			userSelectorErr:          nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                     "UserSelectorErr",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			user:                     userTable.Record{},
			userSelectorErr:          sql.ErrConnDone,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                     "NotAdmin",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			user:                     userTable.Record{IsAdmin: false},
			userSelectorErr:          nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can create tasks.",
			),
		},
		{
			name:                     "BoardWrongTeam",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{TeamID: 1},
			boardSelectorErr:         nil,
			user: userTable.Record{
				IsAdmin: true, TeamID: 2,
			},
			userSelectorErr: nil,
			taskInserterErr: nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:                     "TaskInserterErr",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{TeamID: 1},
			boardSelectorErr:         nil,
			user: userTable.Record{
				IsAdmin: true, TeamID: 1,
			},
			userSelectorErr: nil,
			taskInserterErr: sql.ErrConnDone,
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			taskTitleValidator.Err = c.taskTitleValidatorErr
			subtaskTitleValidator.Err = c.subtaskTitleValidatorErr
			columnSelector.Err = c.columnSelectorErr
			boardSelector.Board = c.board
			boardSelector.Err = c.boardSelectorErr
			userSelector.User = c.user
			userSelector.Err = c.userSelectorErr
			taskInserter.Err = c.taskInserterErr

			reqBody, err := json.Marshal(map[string]any{
				"title":       "",
				"description": "",
				"column":      0,
				"subtasks":    []string{""},
			})
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPost, "", bytes.NewReader(reqBody),
			)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()

			sut.Handle(w, req, "")
			res := w.Result()

			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			c.assertFunc(t, w.Result(), log.InMessage)
		})
	}
}
