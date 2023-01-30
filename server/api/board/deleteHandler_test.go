package board

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/db"
	"server/log"
)

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDELETEHandler(t *testing.T) {
	validator := &fakeDELETEReqValidator{}
	userBoardSelector := &db.FakeRelSelector{}
	userBoardDeleter := &db.FakeDeleter{}
	logger := &log.FakeLogger{}
	sut := NewDELETEHandler(validator, userBoardSelector, userBoardDeleter, logger)

	for _, c := range []struct {
		name                        string
		validatorOutErr             error
		userBoardSelectorOutIsAdmin bool
		userBoardSelectorOutErr     error
		boardDeleterOutErr          error
		wantStatusCode              int
		wantErrMsg                  string
	}{
		{
			name:                        "ValidatorErr",
			validatorOutErr:             errEmptyBoardID,
			userBoardSelectorOutIsAdmin: true,
			userBoardSelectorOutErr:     nil,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusBadRequest,
			wantErrMsg:                  errEmptyBoardID.Error(),
		},
		{
			name:                        "NoRows",
			validatorOutErr:             nil,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     sql.ErrNoRows,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusNotFound,
			wantErrMsg:                  "",
		},
		{
			name:                        "ConnDone",
			validatorOutErr:             nil,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     sql.ErrConnDone,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusInternalServerError,
			wantErrMsg:                  "",
		},
		{
			name:                        "NotAdmin",
			validatorOutErr:             nil,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     nil,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusUnauthorized,
			wantErrMsg:                  "",
		},
		{
			name:                        "DeleteErr",
			validatorOutErr:             nil,
			userBoardSelectorOutIsAdmin: true,
			userBoardSelectorOutErr:     nil,
			boardDeleterOutErr:          errors.New("delete board error"),
			wantStatusCode:              http.StatusInternalServerError,
			wantErrMsg:                  "",
		},
		{
			name:                        "Success",
			validatorOutErr:             nil,
			userBoardSelectorOutIsAdmin: true,
			userBoardSelectorOutErr:     nil,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusOK,
			wantErrMsg:                  "",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			validator.OutErr = c.validatorOutErr
			userBoardSelector.OutIsAdmin = c.userBoardSelectorOutIsAdmin
			userBoardSelector.OutErr = c.userBoardSelectorOutErr
			userBoardDeleter.OutErr = c.boardDeleterOutErr

			req, err := http.NewRequest(http.MethodPost, "/board?id=123", nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()

			sut.Handle(w, req, "")

			if err := assert.Equal(
				c.wantStatusCode, w.Result().StatusCode,
			); err != nil {
				t.Error(err)
			}

			if c.wantStatusCode == http.StatusBadRequest {
				if c.wantErrMsg == "" {
					t.Error(
						"status was 400 but no error messages were expected",
					)
				} else {
					resBody := DELETEResBody{}
					if err := json.NewDecoder(w.Result().Body).Decode(
						&resBody,
					); err != nil {
						t.Error(err)
					}

					if err := assert.Equal(
						c.wantErrMsg, resBody.Error,
					); err != nil {
						t.Error(err)
					}
				}
			}

			if c.wantStatusCode == http.StatusInternalServerError {
				errFound := false
				for _, err := range []error{
					c.userBoardSelectorOutErr,
					c.boardDeleterOutErr,
				} {
					if err != nil {
						errFound = true

						if err := assert.Equal(
							log.LevelError, logger.InLevel,
						); err != nil {
							t.Error(err)
						}

						if err := assert.Equal(
							err.Error(), logger.InMessage,
						); err != nil {
							t.Error(err)
						}
					}
				}
				if !errFound {
					t.Errorf(
						"c.wantStatusCode was %d but no errors were logged.",
						http.StatusInternalServerError,
					)
				}
			}
		})
	}
}
