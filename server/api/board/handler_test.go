package board

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/auth"
	"server/db"
)

func TestHandler(t *testing.T) {
	tokenValidator := &auth.FakeTokenValidator{}
	userBoardCounter := &db.FakeCounter{}
	boardInserter := &db.FakeBoardInserter{}
	sut := NewHandler(tokenValidator, userBoardCounter, boardInserter)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		for _, httpMethod := range []string{
			http.MethodConnect, http.MethodDelete, http.MethodGet,
			http.MethodHead, http.MethodOptions, http.MethodPatch,
			http.MethodPut, http.MethodTrace,
		} {
			t.Run(httpMethod, func(t *testing.T) {
				req, err := http.NewRequest(httpMethod, "/board", nil)
				if err != nil {
					t.Fatal(err)
				}
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)

				if err = assert.Equal(http.StatusMethodNotAllowed, w.Result().StatusCode); err != nil {
					t.Error(err)
				}
			})
		}
	})

	authCookie := &http.Cookie{Name: auth.CookieName, Value: "dummytoken"}
	for _, c := range []struct {
		name                   string
		authCookie             *http.Cookie
		reqBody                ReqBody
		tokenValidatorOutSub   string
		tokenValidatorOutErr   error
		userBoardCounterOutRes int
		boardInserterOutErr    error
		wantStatusCode         int
		wantErr                string
	}{
		{
			name:                   "NoAuthCookie",
			authCookie:             nil,
			reqBody:                ReqBody{},
			tokenValidatorOutSub:   "",
			tokenValidatorOutErr:   nil,
			userBoardCounterOutRes: 0,
			boardInserterOutErr:    nil,
			wantStatusCode:         http.StatusUnauthorized,
			wantErr:                "",
		},
		{
			name:                   "InvalidAuthCookie",
			authCookie:             authCookie,
			reqBody:                ReqBody{},
			tokenValidatorOutSub:   "",
			tokenValidatorOutErr:   errors.New("token validator error"),
			userBoardCounterOutRes: 0,
			boardInserterOutErr:    nil,
			wantStatusCode:         http.StatusUnauthorized,
			wantErr:                "",
		},
		{
			name:                   "BoardNameNil",
			authCookie:             authCookie,
			reqBody:                ReqBody{},
			tokenValidatorOutErr:   nil,
			tokenValidatorOutSub:   "bob21",
			userBoardCounterOutRes: 0,
			boardInserterOutErr:    nil,
			wantStatusCode:         http.StatusBadRequest,
			wantErr:                errNameEmpty,
		},
		{
			name:                   "BoardNameEmpty",
			authCookie:             authCookie,
			reqBody:                ReqBody{Name: ""},
			tokenValidatorOutErr:   nil,
			tokenValidatorOutSub:   "bob21",
			userBoardCounterOutRes: 0,
			boardInserterOutErr:    nil,
			wantStatusCode:         http.StatusBadRequest,
			wantErr:                errNameEmpty,
		},
		{
			name:                   "BoardNameTooLong",
			authCookie:             authCookie,
			reqBody:                ReqBody{Name: "boardyboardsyboardkyboardishboardxyz"},
			tokenValidatorOutErr:   nil,
			tokenValidatorOutSub:   "bob21",
			userBoardCounterOutRes: 0,
			boardInserterOutErr:    nil,
			wantStatusCode:         http.StatusBadRequest,
			wantErr:                errNameTooLong,
		},
		{
			name:                   "MaxBoardsCreated",
			authCookie:             authCookie,
			reqBody:                ReqBody{Name: "someboard"},
			tokenValidatorOutErr:   nil,
			tokenValidatorOutSub:   "bob21",
			userBoardCounterOutRes: 3,
			boardInserterOutErr:    nil,
			wantStatusCode:         http.StatusBadRequest,
			wantErr:                errMaxBoards,
		},
		{
			name:                   "BoardCreatorError",
			authCookie:             authCookie,
			reqBody:                ReqBody{Name: "someboard"},
			tokenValidatorOutErr:   nil,
			tokenValidatorOutSub:   "bob21",
			userBoardCounterOutRes: 0,
			boardInserterOutErr:    errors.New("board creator error"),
			wantStatusCode:         http.StatusInternalServerError,
			wantErr:                "",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			tokenValidator.OutSub = c.tokenValidatorOutSub
			tokenValidator.OutErr = c.tokenValidatorOutErr
			userBoardCounter.OutRes = c.userBoardCounterOutRes
			boardInserter.OutErr = c.boardInserterOutErr

			reqBodyJSON, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPost, "/board", bytes.NewReader(reqBodyJSON),
			)
			if err != nil {
				t.Fatal(err)
			}

			if c.authCookie != nil {
				req.AddCookie(c.authCookie)
			}

			w := httptest.NewRecorder()

			sut.ServeHTTP(w, req)

			if err = assert.Equal(
				c.wantStatusCode, w.Result().StatusCode,
			); err != nil {
				t.Error(err)
			}

			// if 400 is expected - there must be a validation error in response body
			if c.wantStatusCode == http.StatusBadRequest {
				resBody := ResBody{}
				if err := json.NewDecoder(w.Result().Body).Decode(
					&resBody,
				); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(c.wantErr, resBody.Error); err != nil {
					t.Error(err)
				}
			}

			// DEPENDENCY-INPUT-BASED ASSERTIONS

			// if no auth cookie was present, token validator must be called
			if c.authCookie == nil {
				return
			}
			if err := assert.Equal(
				c.authCookie.Value, tokenValidator.InToken,
			); err != nil {
				t.Error(err)
			}

			// if no token validator or board name validation error is expected, board
			// counter must be called
			if c.tokenValidatorOutErr != nil ||
				c.wantErr == errNameEmpty ||
				c.wantErr == errNameTooLong {
				return
			}
			if err := assert.Equal(
				c.tokenValidatorOutSub, userBoardCounter.InID,
			); err != nil {
				t.Error(err)
			}

			// if max boards is not reached, board creator must be called
			if c.userBoardCounterOutRes >= maxBoards {
				return
			}
			if err := assert.Equal(
				db.NewBoard(c.reqBody.Name, c.tokenValidatorOutSub),
				boardInserter.InBoard,
			); err != nil {
				t.Error(err)
			}
		})
	}
}
