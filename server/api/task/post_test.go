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
	log := &pkgLog.FakeErrorer{}
	sut := NewPOSTHandler(
		taskTitleValidator,
		subtaskTitleValidator,
		columnSelector,
		userBoardSelector,
		log,
	)

	for _, c := range []struct {
		name                        string
		taskTitleValidatorOutErr    error
		subtaskTitleValidatorOutErr error
		columnSelectorOutErr        error
		userBoardSelectorOutRes     bool
		userBoardSelectorOutErr     error
		wantStatusCode              int
		assertFunc                  func(*testing.T, *http.Response, string)
	}{
		{
			name:                        "TaskTitleEmpty",
			taskTitleValidatorOutErr:    errTitleEmpty,
			subtaskTitleValidatorOutErr: nil,
			columnSelectorOutErr:        nil,
			userBoardSelectorOutRes:     false,
			userBoardSelectorOutErr:     nil,
			wantStatusCode:              http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be empty.",
			),
		},
		{
			name:                        "TaskTitleTooLong",
			taskTitleValidatorOutErr:    errTitleTooLong,
			subtaskTitleValidatorOutErr: nil,
			columnSelectorOutErr:        nil,
			userBoardSelectorOutRes:     false,
			userBoardSelectorOutErr:     nil,
			wantStatusCode:              http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:                        "SubtaskTitleEmpty",
			taskTitleValidatorOutErr:    nil,
			subtaskTitleValidatorOutErr: errTitleEmpty,
			columnSelectorOutErr:        nil,
			userBoardSelectorOutRes:     false,
			userBoardSelectorOutErr:     nil,
			wantStatusCode:              http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be empty.",
			),
		},
		{
			name:                        "SubtaskTitleTooLong",
			taskTitleValidatorOutErr:    nil,
			subtaskTitleValidatorOutErr: errTitleTooLong,
			columnSelectorOutErr:        nil,
			userBoardSelectorOutRes:     false,
			userBoardSelectorOutErr:     nil,
			wantStatusCode:              http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be longer than 50 characters.",
			),
		},
		{
			name:                        "ColumnNotFound",
			taskTitleValidatorOutErr:    nil,
			subtaskTitleValidatorOutErr: nil,
			columnSelectorOutErr:        sql.ErrNoRows,
			userBoardSelectorOutRes:     false,
			userBoardSelectorOutErr:     nil,
			wantStatusCode:              http.StatusNotFound,
			assertFunc:                  assert.OnResErr("Column not found."),
		},
		{
			name:                        "ColumnSelectorErr",
			taskTitleValidatorOutErr:    nil,
			subtaskTitleValidatorOutErr: nil,
			columnSelectorOutErr:        sql.ErrConnDone,
			userBoardSelectorOutRes:     false,
			userBoardSelectorOutErr:     nil,
			wantStatusCode:              http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                        "NoAccess",
			taskTitleValidatorOutErr:    nil,
			subtaskTitleValidatorOutErr: nil,
			columnSelectorOutErr:        nil,
			userBoardSelectorOutRes:     false,
			userBoardSelectorOutErr:     sql.ErrNoRows,
			wantStatusCode:              http.StatusUnauthorized,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:                        "UserBoardSelectorErr",
			taskTitleValidatorOutErr:    nil,
			subtaskTitleValidatorOutErr: nil,
			columnSelectorOutErr:        nil,
			userBoardSelectorOutRes:     false,
			userBoardSelectorOutErr:     sql.ErrConnDone,
			wantStatusCode:              http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                        "NotAdmin",
			taskTitleValidatorOutErr:    nil,
			subtaskTitleValidatorOutErr: nil,
			columnSelectorOutErr:        nil,
			userBoardSelectorOutRes:     false,
			userBoardSelectorOutErr:     nil,
			wantStatusCode:              http.StatusUnauthorized,
			assertFunc: assert.OnResErr(
				"Only board admins can create tasks.",
			),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			taskTitleValidator.OutErr = c.taskTitleValidatorOutErr
			subtaskTitleValidator.OutErr = c.subtaskTitleValidatorOutErr
			columnSelector.OutErr = c.columnSelectorOutErr
			userBoardSelector.OutErr = c.userBoardSelectorOutErr

			reqBody, err := json.Marshal(map[string]any{
				"title":       "",
				"description": "",
				"column":      0,
				"subtasks":    []map[string]any{{"title": ""}},
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
