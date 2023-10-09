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
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(taskTitleValidator, log)

	t.Run("TaskTitleEmpty", func(t *testing.T) {
		wantStatusCode := http.StatusBadRequest
		wantErrMsg := "Task title cannot be empty."

		taskTitleValidator.Err = errTitleEmpty

		reqBody, err := json.Marshal(map[string]any{
			"column":      0,
			"title":       "",
			"description": "",
			"subtasks":    []string{},
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

		assert.OnResErr(wantErrMsg)(t, res, "")
	})

	t.Run("TaskTitleTooLong", func(t *testing.T) {
		wantStatusCode := http.StatusBadRequest
		wantErrMsg := "Task title cannot be longer than 50 characters."

		taskTitleValidator.Err = errTitleTooLong

		reqBody, err := json.Marshal(map[string]any{
			"title": "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqwe" +
				"asd",
			"description": "",
			"column":      0,
			"subtasks":    []string{},
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

		assert.OnResErr(wantErrMsg)(t, res, "")
	})
}
