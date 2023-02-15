//go:build utest

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
	"server/db"
	"server/log"
)

// TestPOSTHandler tests the Handle method of POSTHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPOSTHandler(t *testing.T) {
	validator := &fakePOSTReqValidator{}
	userBoardCounter := &db.FakeCounter{}
	dbBoardInserter := &db.FakeBoardInserter{}
	logger := &log.FakeLogger{}
	sut := NewPOSTHandler(validator, userBoardCounter, dbBoardInserter, logger)
	sub := "bob123"

	boardInserterErr := errors.New("create board error")

	t.Run(http.MethodPost, func(t *testing.T) {
		for _, c := range []struct {
			name                   string
			validatorOutErrMsg     string
			userBoardCounterOutRes int
			userBoardCounterOutErr error
			boardInserterOutErr    error
			wantStatusCode         int
			wantErrMsg             string
		}{
			{
				name:                   "InvalidRequest",
				validatorOutErrMsg:     msgNameEmpty,
				userBoardCounterOutRes: 0,
				userBoardCounterOutErr: nil,
				boardInserterOutErr:    nil,
				wantStatusCode:         http.StatusBadRequest,
				wantErrMsg:             msgNameEmpty,
			},
			{
				name:                   "UserBoardCounterErr",
				validatorOutErrMsg:     "",
				userBoardCounterOutRes: 0,
				userBoardCounterOutErr: sql.ErrConnDone,
				boardInserterOutErr:    nil,
				wantStatusCode:         http.StatusInternalServerError,
				wantErrMsg:             msgMaxBoards,
			},
			{
				name:                   "MaxBoardsCreated",
				validatorOutErrMsg:     "",
				userBoardCounterOutRes: 3,
				userBoardCounterOutErr: nil,
				boardInserterOutErr:    nil,
				wantStatusCode:         http.StatusBadRequest,
				wantErrMsg:             msgMaxBoards,
			},
			{
				name:                   "BoardInserterErr",
				validatorOutErrMsg:     "",
				userBoardCounterOutRes: 0,
				userBoardCounterOutErr: sql.ErrNoRows,
				boardInserterOutErr:    boardInserterErr,
				wantStatusCode:         http.StatusInternalServerError,
				wantErrMsg:             boardInserterErr.Error(),
			},
			{
				name:                   "Success",
				validatorOutErrMsg:     "",
				userBoardCounterOutRes: 0,
				userBoardCounterOutErr: sql.ErrNoRows,
				boardInserterOutErr:    nil,
				wantStatusCode:         http.StatusOK,
				wantErrMsg:             "",
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				// Set pre-determinate return values for sut's dependencies.
				validator.OutErrMsg = c.validatorOutErrMsg
				userBoardCounter.OutRes = c.userBoardCounterOutRes
				userBoardCounter.OutErr = c.userBoardCounterOutErr
				dbBoardInserter.OutErr = c.boardInserterOutErr

				// Prepare request and response recorder.
				reqBody := POSTReqBody{Name: "My Board"}
				reqBodyJSON, err := json.Marshal(reqBody)
				if err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest(
					http.MethodPost, "/board", bytes.NewReader(reqBodyJSON),
				)
				if err != nil {
					t.Fatal(err)
				}
				w := httptest.NewRecorder()

				// Handle request with sut and get the result.
				sut.Handle(w, req, sub)
				res := w.Result()

				// Assert on the status code.
				if err = assert.Equal(
					c.wantStatusCode, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				switch c.wantStatusCode {
				case http.StatusBadRequest:
					// 400 is expected - there must be a validation error in
					// response body
					if c.wantErrMsg == "" {
						t.Error(
							"400 was expected but no error messages were " +
								"expected",
						)
					} else {
						resBody := POSTResBody{}
						if err = json.NewDecoder(w.Result().Body).Decode(
							&resBody,
						); err != nil {
							t.Error(err)
						}
						if err = assert.Equal(
							c.wantErrMsg, resBody.Error,
						); err != nil {
							t.Error(err)
						}
					}
				case http.StatusInternalServerError:
					// 500 was expected - an error must be logged.
					errFound := false
					for _, depErr := range []error{
						c.userBoardCounterOutErr,
						c.boardInserterOutErr,
					} {
						if depErr != nil && depErr != sql.ErrNoRows {
							errFound = true
							if err = assert.Equal(
								log.LevelError, logger.InLevel,
							); err != nil {
								t.Error(err)
							}
							if err = assert.Equal(
								depErr.Error(), logger.InMessage,
							); err != nil {
								t.Error(err)
							}
						}
					}
					if !errFound {
						t.Errorf(
							"500 was expected but no errors were returned " +
								"from sut's dependencies",
						)
					}
					return
				}
			})
		}
	})
}
