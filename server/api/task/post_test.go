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
	pkgLog "server/log"
)

// TestPOSTHandler tests the Handle method of POSTHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPOSTHandler(t *testing.T) {
	titleValidator := &api.FakeStringValidator{}
	columnSelector := &columnTable.FakeSelector{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPOSTHandler(titleValidator, columnSelector, log)

	for _, c := range []struct {
		name                 string
		titleValidatorOutErr error
		columnSelectorOutErr error
		wantStatusCode       int
		wantErrMsg           string
	}{
		{
			name:                 "TitleEmpty",
			titleValidatorOutErr: errTitleEmpty,
			columnSelectorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
			wantErrMsg:           "Task title cannot be empty.",
		},
		{
			name:                 "TitleTooLong",
			titleValidatorOutErr: errTitleTooLong,
			columnSelectorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
			wantErrMsg: "Task title cannot be longer than 50 " +
				"characters.",
		},
		{
			name:                 "ColumnNotFound",
			titleValidatorOutErr: nil,
			columnSelectorOutErr: sql.ErrNoRows,
			wantStatusCode:       http.StatusNotFound,
			wantErrMsg:           "Column not found.",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			titleValidator.OutErr = c.titleValidatorOutErr
			columnSelector.OutErr = c.columnSelectorOutErr

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

			if err = assert.Equal(
				c.wantStatusCode, w.Result().StatusCode,
			); err != nil {
				t.Error(err)
			}

			var resBody ResBody
			if err = json.NewDecoder(w.Body).Decode(&resBody); err != nil {
				t.Fatal(err)
			}
			if err = assert.Equal(c.wantErrMsg, resBody.Error); err != nil {
				t.Error(err)
			}
		})
	}
}
