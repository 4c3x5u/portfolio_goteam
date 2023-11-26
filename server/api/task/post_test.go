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
	userSelector := &userTable.FakeSelector{}
	taskTitleValidator := &api.FakeStringValidator{}
	subtaskTitleValidator := &api.FakeStringValidator{}
	columnSelector := &columnTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	taskInserter := &taskTable.FakeInserter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPOSTHandler(
		userSelector,
		taskTitleValidator,
		subtaskTitleValidator,
		columnSelector,
		boardSelector,
		taskInserter,
		log,
	)

	for _, c := range []struct {
		name                     string
		user                     userTable.Record
		userSelectorErr          error
		taskTitleValidatorErr    error
		subtaskTitleValidatorErr error
		columnSelectorErr        error
		board                    boardTable.Record
		boardSelectorErr         error
		taskInserterErr          error
		wantStatusCode           int
		assertFunc               func(*testing.T, *http.Response, string)
	}{
		{
			name:                     "UserNotRecognised",
			user:                     userTable.Record{},
			userSelectorErr:          sql.ErrNoRows,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusUnauthorized,
			assertFunc: assert.OnResErr(
				"Username is not recognised.",
			),
		},
		{
			name:                     "UserSelectorErr",
			user:                     userTable.Record{},
			userSelectorErr:          sql.ErrConnDone,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                     "NotAdmin",
			user:                     userTable.Record{IsAdmin: false},
			userSelectorErr:          nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can create tasks.",
			),
		},
		{
			name:                     "TaskTitleEmpty",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskTitleValidatorErr:    api.ErrStrEmpty,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be empty.",
			),
		},
		{
			name:                     "TaskTitleTooLong",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskTitleValidatorErr:    api.ErrStrTooLong,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:                     "TaskTitleUnexpectedErr",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskTitleValidatorErr:    api.ErrStrNotInt,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				api.ErrStrNotInt.Error(),
			),
		},
		{
			name:                     "SubtaskTitleEmpty",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrEmpty,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be empty.",
			),
		},
		{
			name:                     "SubtaskTitleTooLong",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrTooLong,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be longer than 50 characters.",
			),
		},
		{
			name:                     "SubtaskTitleUnexpectedErr",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrNotInt,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				api.ErrStrNotInt.Error(),
			),
		},
		{
			name:                     "ColumnNotFound",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        sql.ErrNoRows,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusNotFound,
			assertFunc:               assert.OnResErr("Column not found."),
		},
		{
			name:                     "ColumnSelectorErr",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        sql.ErrConnDone,
			board:                    boardTable.Record{},
			boardSelectorErr:         nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                     "BoardSelectorErr",
			user:                     userTable.Record{IsAdmin: true},
			userSelectorErr:          nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{},
			boardSelectorErr:         sql.ErrConnDone,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name: "BoardWrongTeam",
			user: userTable.Record{
				IsAdmin: true, TeamID: 2,
			},
			userSelectorErr:          nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{TeamID: 1},
			boardSelectorErr:         nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name: "TaskInserterErr",
			user: userTable.Record{
				IsAdmin: true, TeamID: 1,
			},
			userSelectorErr:          nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			board:                    boardTable.Record{TeamID: 1},
			boardSelectorErr:         nil,
			taskInserterErr:          sql.ErrConnDone,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userSelector.Rec = c.user
			userSelector.Err = c.userSelectorErr
			taskTitleValidator.Err = c.taskTitleValidatorErr
			subtaskTitleValidator.Err = c.subtaskTitleValidatorErr
			columnSelector.Err = c.columnSelectorErr
			boardSelector.Board = c.board
			boardSelector.Err = c.boardSelectorErr
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
