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
	taskIDValidator := &api.FakeStringValidator{}
	taskTitleValidator := &api.FakeStringValidator{}
	subtaskTitleValidator := &api.FakeStringValidator{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(
		taskIDValidator,
		taskTitleValidator,
		subtaskTitleValidator,
		log,
	)

	for _, c := range []struct {
		name                     string
		taskIDValidatorErr       error
		taskTitleValidatorErr    error
		subtaskTitleValidatorErr error
		wantErrMsg               string
	}{
		{
			name:                     "TaskIDEmpty",
			taskIDValidatorErr:       api.ErrValueEmpty,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			wantErrMsg:               "Task ID cannot be empty.",
		},
		{
			name:                     "TaskTitleEmpty",
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    api.ErrValueEmpty,
			subtaskTitleValidatorErr: nil,
			wantErrMsg:               "Task title cannot be empty.",
		},
		{
			name:                     "TaskTitleTooLong",
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    api.ErrValueTooLong,
			subtaskTitleValidatorErr: nil,
			wantErrMsg: "Task title cannot be longer than 50 " +
				"characters.",
		},
		{
			name:                     "SubtaskTitleEmpty",
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrValueEmpty,
			wantErrMsg:               "Subtask title cannot be empty.",
		},
		{
			name:                     "SubtaskTitleTooLong",
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrValueTooLong,
			wantErrMsg: "Subtask title cannot be longer than 50 " +
				"characters.",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			wantStatusCode := http.StatusBadRequest

			taskIDValidator.Err = c.taskIDValidatorErr
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
