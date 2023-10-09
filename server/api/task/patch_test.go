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
		wantStatusCode           int
		assertFunc               func(*testing.T, *http.Response, string)
	}{
		{
			name:                     "TaskIDEmpty",
			taskIDValidatorErr:       api.ErrStrEmpty,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task ID cannot be empty.",
			),
		},
		{
			name:                     "TaskIDNotInt",
			taskIDValidatorErr:       api.ErrStrNotInt,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task ID must be an integer.",
			),
		},
		{
			name:                     "TaskIDUnexpectedErr",
			taskIDValidatorErr:       api.ErrStrTooLong,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				api.ErrStrTooLong.Error(),
			),
		},
		{
			name:                     "TaskTitleEmpty",
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    api.ErrStrEmpty,
			subtaskTitleValidatorErr: nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be empty.",
			),
		},
		{
			name:                     "TaskTitleTooLong",
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    api.ErrStrTooLong,
			subtaskTitleValidatorErr: nil,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:                     "TaskTitleUnexpectedErr",
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    api.ErrStrNotInt,
			subtaskTitleValidatorErr: nil,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				api.ErrStrNotInt.Error(),
			),
		},
		{
			name:                     "SubtaskTitleEmpty",
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrEmpty,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be empty.",
			),
		},
		{
			name:                     "SubtaskTitleTooLong",
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrTooLong,
			wantStatusCode:           http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be longer than 50 characters.",
			),
		},
		{
			name:                     "SubtaskTitleUnexpectedErr",
			taskIDValidatorErr:       nil,
			taskTitleValidatorErr:    nil,
			subtaskTitleValidatorErr: api.ErrStrNotInt,
			wantStatusCode:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				api.ErrStrNotInt.Error(),
			),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
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

			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			c.assertFunc(t, res, log.InMessage)
		})
	}
}
