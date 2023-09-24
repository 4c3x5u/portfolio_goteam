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
	boardSelector := &dbaccess.FakeBoardSelector{}
	sut := NewPATCHHandler(idValidator, nameValidator, boardSelector, log)

	for _, c := range []struct {
		name                string
		idValidatorOutErr   error
		nameValidatorOutErr error
		boardSelectorOutErr error
		wantStatusCode      int
		assertFunc          func(*testing.T, *http.Response)
	}{
		{
			name:                "IDValidatorErr",
			idValidatorOutErr:   errors.New("Board ID cannot be empty."),
			nameValidatorOutErr: nil,
			boardSelectorOutErr: nil,
			wantStatusCode:      http.StatusBadRequest,
			assertFunc: func(t *testing.T, res *http.Response) {
				var resBody ResBody
				if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}
				if err := assert.Equal(
					"Board ID cannot be empty.", resBody.Error,
				); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:                "NameValidatorErr",
			idValidatorOutErr:   nil,
			nameValidatorOutErr: errors.New("Board name cannot be empty."),
			boardSelectorOutErr: nil,
			wantStatusCode:      http.StatusBadRequest,
			assertFunc: func(t *testing.T, res *http.Response) {
				var resBody ResBody
				if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}
				if err := assert.Equal(
					"Board name cannot be empty.", resBody.Error,
				); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:                "BoardNotFound",
			idValidatorOutErr:   nil,
			nameValidatorOutErr: nil,
			boardSelectorOutErr: sql.ErrNoRows,
			wantStatusCode:      http.StatusNotFound,
			assertFunc: func(t *testing.T, res *http.Response) {
				var resBody ResBody
				if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}
				if err := assert.Equal(
					"Board not found.", resBody.Error,
				); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:                "BoardSelectorErr",
			idValidatorOutErr:   nil,
			nameValidatorOutErr: nil,
			boardSelectorOutErr: sql.ErrConnDone,
			wantStatusCode:      http.StatusInternalServerError,
			assertFunc: func(t *testing.T, res *http.Response) {
				if err := assert.Equal(
					sql.ErrConnDone.Error(), log.InMessage,
				); err != nil {
					t.Error(err)
				}
			},
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

			c.assertFunc(t, res)
		})
	}
}
