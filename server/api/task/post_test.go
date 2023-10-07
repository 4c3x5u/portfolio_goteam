//go:build utest

package task

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"server/api"
	"server/assert"
	"testing"
)

// TestPOSTHandler tests the Handle method of POSTHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPOSTHandler(t *testing.T) {
	titleValidator := &api.FakeStringValidator{}
	sut := NewPOSTHandler(titleValidator)

	t.Run("ValidatorErr", func(t *testing.T) {
		wantStatusCode := http.StatusBadRequest
		wantErr := errors.New("task title invalid")

		titleValidator.OutErr = wantErr

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
		if err = assert.Equal(wantErr.Error(), resBody.Error); err != nil {
			t.Error(err)
		}
	})
}
