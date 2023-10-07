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
	titleValidator := &api.FakeStringValidator{}
	columnSelector := &columnTable.FakeSelector{}
	userBoardSelector := &userboardTable.FakeSelector{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPOSTHandler(
		titleValidator, columnSelector, userBoardSelector, log,
	)

	for _, c := range []struct {
		name                    string
		titleValidatorOutErr    error
		columnSelectorOutErr    error
		userBoardSelectorOutErr error
		wantStatusCode          int
		assertFunc              func(*testing.T, *http.Response, string)
	}{
		{
			name:                    "TitleEmpty",
			titleValidatorOutErr:    errTitleEmpty,
			columnSelectorOutErr:    nil,
			userBoardSelectorOutErr: nil,
			wantStatusCode:          http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be empty.",
			),
		},
		{
			name:                    "TitleTooLong",
			titleValidatorOutErr:    errTitleTooLong,
			columnSelectorOutErr:    nil,
			userBoardSelectorOutErr: nil,
			wantStatusCode:          http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:                    "ColumnNotFound",
			titleValidatorOutErr:    nil,
			columnSelectorOutErr:    sql.ErrNoRows,
			userBoardSelectorOutErr: nil,
			wantStatusCode:          http.StatusNotFound,
			assertFunc:              assert.OnResErr("Column not found."),
		},
		{
			name:                    "ColumnSelectorErr",
			titleValidatorOutErr:    nil,
			columnSelectorOutErr:    sql.ErrConnDone,
			userBoardSelectorOutErr: nil,
			wantStatusCode:          http.StatusInternalServerError,
			assertFunc:              assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:                    "NoAccess",
			titleValidatorOutErr:    nil,
			columnSelectorOutErr:    nil,
			userBoardSelectorOutErr: sql.ErrNoRows,
			wantStatusCode:          http.StatusUnauthorized,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			titleValidator.OutErr = c.titleValidatorOutErr
			columnSelector.OutErr = c.columnSelectorOutErr
			userBoardSelector.OutErr = c.userBoardSelectorOutErr

			reqBody, err := json.Marshal(map[string]any{
				"title":       "",
				"description": "",
				"column":      0,
				"subtasks":    []map[string]any{},
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
