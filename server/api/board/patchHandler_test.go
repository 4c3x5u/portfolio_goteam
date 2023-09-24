package board

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/dbaccess"
	pkgLog "server/log"
)

func TestPATCHHandler(t *testing.T) {
	log := &pkgLog.FakeErrorer{}
	idValidator := &fakeStringValidator{}
	nameValidator := &fakeStringValidator{}
	sut := NewPATCHHandler(idValidator, nameValidator, log)
	boardSelector := &dbaccess.FakeBoardSelector{}
	sut := NewPATCHHandler(idValidator, nameValidator, boardSelector, log)

	for _, c := range []struct {
		name                string
		idValidatorOutErr   error
		nameValidatorOutErr error
		boardSelectorOutErr error
		wantStatusCode      int
		wantErrMsg          string
	}{
		{
			name:                "IDValidatorErr",
			idValidatorOutErr:   errors.New("Board ID cannot be empty."),
			nameValidatorOutErr: nil,
			boardSelectorOutErr: nil,
			wantStatusCode:      http.StatusBadRequest,
			wantErrMsg:          "Board ID cannot be empty.",
		},
		{
			name:                "NameValidatorErr",
			idValidatorOutErr:   nil,
			nameValidatorOutErr: errors.New("Board name cannot be empty."),
			boardSelectorOutErr: nil,
			wantStatusCode:      http.StatusBadRequest,
			wantErrMsg:          "Board name cannot be empty.",
		},
		{
			name:                "BoardNotFound",
			idValidatorOutErr:   nil,
			nameValidatorOutErr: nil,
			boardSelectorOutErr: sql.ErrNoRows,
			wantStatusCode:      http.StatusNotFound,
			wantErrMsg:          "Board not found.",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			idValidator.OutErr = c.idValidatorOutErr
			nameValidator.OutErr = c.nameValidatorOutErr
			boardSelector.OutErr = c.boardSelectorOutErr

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
			w := httptest.NewRecorder()

			sut.Handle(w, req, "")
			res := w.Result()

			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			var resBody ResBody
			if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
				t.Fatal(err)
			}
			if err = assert.Equal(c.wantErrMsg, resBody.Error); err != nil {
				t.Error(err)
			}
		})
	}
}
