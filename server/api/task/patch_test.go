//go:build utest

package task

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/api"
	"server/assert"
	pkgLog "server/log"
)

// TestPATCHHandler tests the Handle method of PATCHHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	taskTitleValidator := &api.FakeStringValidator{}
	subtaskTitleValidator := &api.FakeStringValidator{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(taskTitleValidator, subtaskTitleValidator, log)

	for _, c := range []struct {
		name                     string
		taskTitleValidatorErr    error
		subtaskTitleValidatorErr error
		wantErrMsg               string
	}{
		{
			name:                     "TaskTitleEmpty",
			taskTitleValidatorErr:    errTitleEmpty,
			subtaskTitleValidatorErr: nil,
			wantErrMsg:               "Task title cannot be empty.",
		},
		{
			name:                     "TaskTitleTooLong",
			taskTitleValidatorErr:    errTitleTooLong,
			subtaskTitleValidatorErr: nil,
			wantErrMsg: "Task title cannot be longer than 50 " +
				"characters.",
		},
		{
			name:                     "SubtaskTitleEmpty",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: errTitleEmpty,
			wantErrMsg:               "Subtask title cannot be empty.",
		},
		{
			name:                     "SubtaskTitleTooLong",
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: errTitleTooLong,
			wantErrMsg: "Subtask title cannot be longer than 50 " +
				"characters.",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			wantStatusCode := http.StatusBadRequest

			taskTitleValidator.Err = c.taskTitleValidatorErr
			subtaskTitleValidator.Err = c.subtaskTitleValidatorErr

			reqBody, err := json.Marshal(map[string]any{
				"column":      0,
				"title":       "",
				"description": "",
				"subtasks":    []string{""},
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

			if err = assert.Equal(wantStatusCode, res.StatusCode); err != nil {
				t.Error(err)
			}

			assert.OnResErr(c.wantErrMsg)(t, res, "")
		})
	}
}
