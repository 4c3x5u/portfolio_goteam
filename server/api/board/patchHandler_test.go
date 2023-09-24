package board

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	pkgLog "server/log"
)

func TestPATCHHandler(t *testing.T) {
	t.Run("IDValidatorErr", func(t *testing.T) {
		log := &pkgLog.FakeErrorer{}
		idValidator := &fakeIDValidator{}
		sut := NewPATCHHandler(idValidator, log)

		wantErrMsg := "Board ID cannot be empty."
		idValidator.OutErr = errors.New(wantErrMsg)

		wantStatusCode := http.StatusBadRequest

		reqBody, err := json.Marshal(ReqBody{})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(
			http.MethodPatch, "", bytes.NewReader(reqBody),
		)
		if err != nil {
			t.Fatal(err)
		}
		req.RequestURI = "board/?id=0"
		w := httptest.NewRecorder()

		sut.Handle(w, req, "")
		res := w.Result()

		if err = assert.Equal(wantStatusCode, res.StatusCode); err != nil {
			t.Error(err)
		}

		var resBody ResBody
		if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Fatal(err)
		}
		if err = assert.Equal(wantErrMsg, resBody.Error); err != nil {
			t.Error(err)
		}
	})
}
