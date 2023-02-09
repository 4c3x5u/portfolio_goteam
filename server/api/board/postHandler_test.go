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
				validator.OutErrMsg = c.validatorOutErrMsg
				userBoardCounter.OutRes = c.userBoardCounterOutRes
				userBoardCounter.OutErr = c.userBoardCounterOutErr
				dbBoardInserter.OutErr = c.boardInserterOutErr

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

				sut.Handle(w, req, sub)

				if err = assert.Equal(
					c.wantStatusCode, w.Result().StatusCode,
				); err != nil {
					t.Error(err)
				}

				// if 400 is expected - there must be a validation error in
				// response body
				if c.wantStatusCode == http.StatusBadRequest {
					if c.wantErrMsg == "" {
						t.Error(
							"status was 400 but no error messages were " +
								"expected",
						)
					} else {
						resBody := POSTResBody{}
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

				// DEPENDENCY-INPUT-BASED ASSERTIONS

				if c.wantStatusCode == http.StatusInternalServerError {
					errFound := false
					for _, err := range []error{
						c.userBoardCounterOutErr, c.boardInserterOutErr,
					} {
						if err != nil && err != sql.ErrNoRows {
							errFound = true

							if levelErr := assert.Equal(
								log.LevelError, logger.InLevel,
							); levelErr != nil {
								t.Error(levelErr)
							}

							if msgErr := assert.Equal(
								err.Error(),
								logger.InMessage,
							); msgErr != nil {
								t.Error(msgErr)
							}
						}
					}
					if !errFound {
						t.Error(
							"c.wantStatusCode was 500 but no errors were " +
								"logged.",
						)
					}
					return
				}

				// if max boards is not reached, board creator must be called
				if c.userBoardCounterOutRes >= maxBoards ||
					c.validatorOutErrMsg != "" {
					return
				}
				if err := assert.Equal(
					db.NewBoard(reqBody.Name, sub),
					dbBoardInserter.InBoard,
				); err != nil {
					t.Error(err)
				}
			})
		}
	})
}
