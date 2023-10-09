//go:build utest

package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"database/sql"
	"server/api"
	"server/assert"
	columnTable "server/dbaccess/column"
	taskTable "server/dbaccess/task"
	userboardTable "server/dbaccess/userboard"
	pkgLog "server/log"
)

// TestPOSTHandler tests the Handle method of POSTHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPOSTHandler(t *testing.T) {
	taskTitleValidator := &api.FakeStringValidator{}
	subtaskTitleValidator := &api.FakeStringValidator{}
	columnSelector := &columnTable.FakeSelector{}
	userBoardSelector := &userboardTable.FakeSelector{}
	taskInserter := &taskTable.FakeInserter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPOSTHandler(
		taskTitleValidator,
		subtaskTitleValidator,
		columnSelector,
		userBoardSelector,
		taskInserter,
		log,
	)

	for _, c := range []struct {
		name                     string
		taskTitleValidatorErr    error
		subtaskTitleValidatorErr error
		columnSelectorErr        error
		userIsAdmin              bool
		userBoardSelectorErr     error
		taskInserterErr          error
		wantStatusCode           int
		assertFunc               func(*testing.T, *http.Response, string)
	}{
		{
			name:                     "TaskTitleEmpty",
			taskTitleValidatorErr:    api.ErrStrEmpty,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			userIsAdmin:              false,
			userBoardSelectorErr:     nil,
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
			userIsAdmin:              false,
			userBoardSelectorErr:     nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:                     "SubtaskTitleEmpty",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrEmpty,
			columnSelectorErr:        nil,
			userIsAdmin:              false,
			userBoardSelectorErr:     nil,
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
			userIsAdmin:              false,
			userBoardSelectorErr:     nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be longer than 50 characters.",
			),
		},
		{
			name:                     "ColumnNotFound",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        sql.ErrNoRows,
			userIsAdmin:              false,
			userBoardSelectorErr:     nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusNotFound,
			assertFunc:               assert.OnResErr("Column not found."),
		},
		{
			name:                     "ColumnSelectorErr",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        sql.ErrConnDone,
			userIsAdmin:              false,
			userBoardSelectorErr:     nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                     "NoAccess",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			userIsAdmin:              false,
			userBoardSelectorErr:     sql.ErrNoRows,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusUnauthorized,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:                     "UserBoardSelectorErr",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			userIsAdmin:              false,
			userBoardSelectorErr:     sql.ErrConnDone,
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
			userIsAdmin:              false,
			userBoardSelectorErr:     nil,
			taskInserterErr:          nil,
			wantStatusCode:           http.StatusUnauthorized,
			assertFunc: assert.OnResErr(
				"Only board admins can create tasks.",
			),
		},
		{
			name:                     "TaskInserterErr",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			columnSelectorErr:        nil,
			userIsAdmin:              true,
			userBoardSelectorErr:     nil,
			taskInserterErr:          sql.ErrConnDone,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			taskTitleValidator.Err = c.taskTitleValidatorErr
			subtaskTitleValidator.Err = c.subtaskTitleValidatorErr
			columnSelector.Err = c.columnSelectorErr
			userBoardSelector.IsAdmin = c.userIsAdmin
			userBoardSelector.Err = c.userBoardSelectorErr
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
