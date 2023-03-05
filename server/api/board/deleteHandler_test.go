//go:build utest

package board

import (
	"database/sql"
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
	sut := NewDELETEHandler(
		validator, userBoardSelector, userBoardDeleter, logger,
	)

	// used in 500 error cases to assert on the logged error message
	assertOnLoggedErr := func(wantErrMsg string) func(*testing.T) {
		return func(t *testing.T) {
			if levelErr := assert.Equal(
				log.LevelError, logger.InLevel,
			); levelErr != nil {
				t.Error(levelErr)
			}

			if msgErr := assert.Equal(
				wantErrMsg, logger.InMessage,
			); msgErr != nil {
				t.Error(msgErr)
			}
		}
	}

	for _, c := range []struct {
		name                        string
		validatorOutOK              bool
		userBoardSelectorOutIsAdmin bool
		userBoardSelectorOutErr     error
		boardDeleterOutErr          error
		wantStatusCode              int
		assertFunc                  func(*testing.T)
	}{
		{
			name:                        "ValidatorErr",
			validatorOutOK:              false,
			userBoardSelectorOutIsAdmin: true,
			userBoardSelectorOutErr:     nil,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusBadRequest,
			assertFunc:                  func(t *testing.T) {},
		},
		{
			name:                        "NoRowsInUserBoardTable",
			validatorOutOK:              true,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     sql.ErrNoRows,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusUnauthorized,
			assertFunc:                  func(t *testing.T) {},
		},
		{
			name:                        "ConnDone",
			validatorOutOK:              true,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     sql.ErrConnDone,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusInternalServerError,
			assertFunc:                  assertOnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:                        "NotAdmin",
			validatorOutOK:              true,
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     nil,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusUnauthorized,
			assertFunc:                  func(t *testing.T) {},
		},
		{
			name:                        "DeleteErr",
			validatorOutOK:              true,
			userBoardSelectorOutIsAdmin: true,
			userBoardSelectorOutErr:     nil,
			boardDeleterOutErr:          errors.New("delete board error"),
			wantStatusCode:              http.StatusInternalServerError,
			assertFunc:                  assertOnLoggedErr("delete board error"),
		},
		{
			name:                        "Success",
			validatorOutOK:              true,
			userBoardSelectorOutIsAdmin: true,
			userBoardSelectorOutErr:     nil,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusOK,
			assertFunc:                  func(t *testing.T) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for sut's dependencies.
			validator.OutOK = c.validatorOutOK
			userBoardSelector.OutIsAdmin = c.userBoardSelectorOutIsAdmin
			userBoardSelector.OutErr = c.userBoardSelectorOutErr
			userBoardDeleter.OutErr = c.boardDeleterOutErr

			// Prepare request and response recorder.
			req, err := http.NewRequest(http.MethodPost, "?id=123", nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.Handle(w, req, "")
			res := w.Result()

			// Assert on the status code.
			if err := assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			// Run case-specific assertions.
			c.assertFunc(t)
		})
	}
}
