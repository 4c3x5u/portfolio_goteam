//go:build utest

package board

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/db"
	boardTable "server/db/board"
	pkgLog "server/log"
)

// TestPOSTHandler tests the Handle method of POSTHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPOSTHandler(t *testing.T) {
	validator := &fakeStringValidator{}
	userBoardCounter := &db.FakeCounter{}
	dbBoardInserter := &boardTable.FakeInserter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPOSTHandler(validator, userBoardCounter, dbBoardInserter, log)

	// Used in status 500 cases to assert on the logged error message.
	assertOnLoggedErr := func(
		wantErrMsg string,
	) func(*testing.T, *pkgLog.FakeErrorer, io.ReadCloser) {
		return func(t *testing.T, l *pkgLog.FakeErrorer, _ io.ReadCloser) {
			if err := assert.Equal(wantErrMsg, l.InMessage); err != nil {
				t.Error(err)
			}
		}
	}

	// Used in status 400 cases to assert on the error returned in res body.
	assertOnResErr := func(
		wantErrMsg string,
	) func(*testing.T, *pkgLog.FakeErrorer, io.ReadCloser) {
		return func(
			t *testing.T, _ *pkgLog.FakeErrorer, rawResBody io.ReadCloser,
		) {
			resBody := ResBody{}
			if err := json.NewDecoder(rawResBody).Decode(&resBody); err != nil {
				t.Fatal(err)
			}
			if err := assert.Equal(wantErrMsg, resBody.Error); err != nil {
				t.Error(err)
			}
		}
	}

	t.Run(http.MethodPost, func(t *testing.T) {
		for _, c := range []struct {
			name                   string
			validatorOutErr        error
			userBoardCounterOutRes int
			userBoardCounterOutErr error
			boardInserterOutErr    error
			wantStatusCode         int
			assertFunc             func(
				*testing.T, *pkgLog.FakeErrorer, io.ReadCloser,
			)
		}{
			{
				name:                   "InvalidRequest",
				validatorOutErr:        errors.New("Board name cannot be empty."),
				userBoardCounterOutRes: 0,
				userBoardCounterOutErr: nil,
				boardInserterOutErr:    nil,
				wantStatusCode:         http.StatusBadRequest,
				assertFunc: assertOnResErr(
					"Board name cannot be empty.",
				),
			},
			{
				name:                   "UserBoardCounterErr",
				validatorOutErr:        nil,
				userBoardCounterOutRes: 0,
				userBoardCounterOutErr: sql.ErrConnDone,
				boardInserterOutErr:    nil,
				wantStatusCode:         http.StatusInternalServerError,
				assertFunc: assertOnLoggedErr(
					sql.ErrConnDone.Error(),
				),
			},
			{
				name:                   "MaxBoardsCreated",
				validatorOutErr:        nil,
				userBoardCounterOutRes: 3,
				userBoardCounterOutErr: nil,
				boardInserterOutErr:    nil,
				wantStatusCode:         http.StatusBadRequest,
				assertFunc: assertOnResErr(
					"You have already created the maximum amount of boards " +
						"allowed per user. Please delete one of your boards " +
						"to create a new one.",
				),
			},
			{
				name:                   "BoardInserterErr",
				validatorOutErr:        nil,
				userBoardCounterOutRes: 0,
				userBoardCounterOutErr: sql.ErrNoRows,
				boardInserterOutErr:    errors.New("create board error"),
				wantStatusCode:         http.StatusInternalServerError,
				assertFunc:             assertOnLoggedErr("create board error"),
			},
			{
				name:                   "Success",
				validatorOutErr:        nil,
				userBoardCounterOutRes: 0,
				userBoardCounterOutErr: sql.ErrNoRows,
				boardInserterOutErr:    nil,
				wantStatusCode:         http.StatusOK,
				assertFunc: func(
					*testing.T, *pkgLog.FakeErrorer, io.ReadCloser,
				) {
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				// Set pre-determinate return values for sut's dependencies.
				validator.OutErr = c.validatorOutErr
				userBoardCounter.OutRes = c.userBoardCounterOutRes
				userBoardCounter.OutErr = c.userBoardCounterOutErr
				dbBoardInserter.OutErr = c.boardInserterOutErr

				// Prepare request and response recorder.
				reqBody, err := json.Marshal(ReqBody{})
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

				// Handle request with sut and get the result.
				sut.Handle(w, req, "")
				res := w.Result()

				// Assert on the status code.
				if err = assert.Equal(
					c.wantStatusCode, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				// Run case-specific assertions.
				c.assertFunc(t, log, w.Result().Body)
			})
		}
	})
}
