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
)

// TestPOSTHandler tests the Handle method of POSTHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPOSTHandler(t *testing.T) {
	titleValidator := &api.FakeStringValidator{}
	sut := NewPOSTHandler(titleValidator)

	for _, c := range []struct {
		name                 string
		titleValidatorOutErr error
		wantErrMsg           string
	}{
		{
			name:                 "TitleEmpty",
			titleValidatorOutErr: errTitleEmpty,
			wantErrMsg:           "Task title cannot be empty.",
		},
		{
			name:                 "TitleTooLong",
			titleValidatorOutErr: errTitleTooLong,
			wantErrMsg: "Task title cannot be longer than 50 " +
				"characters.",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			wantStatusCode := http.StatusBadRequest

			titleValidator.OutErr = c.titleValidatorOutErr

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
				wantStatusCode, w.Result().StatusCode,
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
