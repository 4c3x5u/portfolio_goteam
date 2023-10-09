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

	for _, c := range []struct {
		name                  string
		taskTitleValidatorErr error
		wantErrMsg            string
	}{} {
		t.Run(c.name, func(t *testing.T) {
			wantStatusCode := http.StatusBadRequest

			taskTitleValidator.Err = c.taskTitleValidatorErr

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

			assert.OnResErr(c.wantErrMsg)(t, res, "")
		})
	}
}
