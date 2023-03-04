//go:build itest

package itest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/api"
	boardAPI "server/api/board"
	"server/assert"
	"server/auth"
	"server/db"
	"server/log"
)

func TestBoard(t *testing.T) {
	// Create board API handler.
	logger := log.NewAppLogger()
	sut := boardAPI.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodPost: boardAPI.NewPOSTHandler(
				boardAPI.NewPOSTValidator(),
				db.NewUserBoardCounter(dbConnPool),
				db.NewBoardInserter(dbConnPool),
				logger,
			),
			http.MethodDelete: boardAPI.NewDELETEHandler(
				boardAPI.NewDELETEValidator(),
				db.NewUserBoardSelector(dbConnPool),
				db.NewBoardDeleter(dbConnPool),
				logger,
			),
		},
	)

	for _, c := range []struct {
		name           string
		authHeader     string
		boardName      string
		wantStatusCode int
		wantErrMsg     string
		userID         string // only used for board count on POST success
	}{
		// Auth Cases
		{
			name:           "NoAuthHeader",
			authHeader:     "",
			boardName:      "",
			wantStatusCode: http.StatusUnauthorized,
			wantErrMsg:     "",
			userID:         "",
		},
		{
			name: "InvalidAuthHeader",
			authHeader: "eyJhbGciOiJIUzI1NiNowAsEtqKSQauaqow1.eyJzdWIiOiJib2I" +
				"xMjMifQ.Y8_6K50EHUEJlJf4X21fNCFhYWhVIqN3Tw1niz8XwZc",
			boardName:      "",
			wantStatusCode: http.StatusUnauthorized,
			wantErrMsg:     "",
			userID:         "",
		},

		// POST Cases
		{
			name: "EmptyBoardName",
			authHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2I" +
				"xMjMifQ.Y8_6K50EHUEJlJf4X21fNCFhYWhVIqN3Tw1niz8XwZc",
			boardName:      "",
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     "Board name cannot be empty.",
			userID:         "",
		},
		{
			name: "TooLongBoardName",
			authHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2I" +
				"xMjMifQ.Y8_6K50EHUEJlJf4X21fNCFhYWhVIqN3Tw1niz8XwZc",
			boardName:      "A Board Whose Name Is Just Too Long!",
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     "Board name cannot be longer than 35 characters.",
			userID:         "",
		},
		{
			name: "TooManyBoards",
			authHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2I" +
				"xMjMifQ.Y8_6K50EHUEJlJf4X21fNCFhYWhVIqN3Tw1niz8XwZc",
			boardName:      "bob123's new board",
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg: "You have already created the maximum amount of " +
				"boards allowed per user. Please delete one of your boards " +
				"to create a new one.",
			userID: "",
		},
		{
			name: "Succeess",
			authHeader: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2I" +
				"xMjQifQ.LqENrj9APUHgQ3X0HRN6-IFMIg6nyo0_n74KfoxA0qI",
			boardName:      "bob124's new board",
			wantStatusCode: http.StatusOK,
			wantErrMsg:     "",
			userID:         "bob124",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			reqBody, err := json.Marshal(boardAPI.POSTReqBody{
				Name: c.boardName,
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
			if c.authHeader != "" {
				req.Header.Add("Authorization", "Bearer "+c.authHeader)
			}
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, req)

			res := w.Result()

			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			switch c.wantStatusCode {
			case http.StatusUnauthorized:
				authResHeader := res.Header.Values("WWW-Authenticate")[0]
				if err = assert.Equal("Bearer", authResHeader); err != nil {
					t.Error(err)
				}
			case http.StatusBadRequest:
				resBody := boardAPI.POSTResBody{}
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
			case http.StatusOK:
				var boardCount int
				err = dbConnPool.QueryRow(
					"SELECT COUNT(*) FROM app.user_board "+
						"WHERE userID = $1 AND isAdmin = TRUE",
					c.userID,
				).Scan(&boardCount)
				if err != nil {
					t.Error(err)
				}
				if err = assert.Equal(1, boardCount); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
