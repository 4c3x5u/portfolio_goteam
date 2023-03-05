//go:build itest

package itest

import (
	"bytes"
	"database/sql"
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
				db.NewUserBoardCounter(dbConn),
				db.NewBoardInserter(dbConn),
				logger,
			),
			http.MethodDelete: boardAPI.NewDELETEHandler(
				boardAPI.NewDELETEValidator(),
				db.NewUserBoardSelector(dbConn),
				db.NewBoardDeleter(dbConn),
				logger,
			),
		},
	)

	// used in varioues test cases to authenticate the request sent
	addBearerAuth := func(token string) func(*http.Request) {
		return func(req *http.Request) {
			req.Header.Add("Authorization", "Bearer "+token)
		}
	}

	// used in 400 error cases to assert on the error message
	assertOnErrMsg := func(
		wantErrMsg string,
	) func(*testing.T, *httptest.ResponseRecorder) {
		return func(t *testing.T, w *httptest.ResponseRecorder) {
			resBody := boardAPI.POSTResBody{}
			if err := json.NewDecoder(w.Result().Body).Decode(
				&resBody,
			); err != nil {
				t.Error(err)
			}
			if err := assert.Equal(
				wantErrMsg, resBody.Error,
			); err != nil {
				t.Error(err)
			}
		}
	}

	t.Run("Auth", func(t *testing.T) {
		for _, c := range []struct {
			name     string
			authFunc func(*http.Request)
		}{
			// Auth Cases
			{name: "HeaderEmpty", authFunc: func(*http.Request) {}},
			{name: "HeaderInvalid", authFunc: addBearerAuth("asdfasldfkjasd")},
		} {
			t.Run(c.name, func(t *testing.T) {
				for _, method := range []string{
					http.MethodPost, http.MethodDelete,
				} {
					t.Run(method, func(t *testing.T) {
						req, err := http.NewRequest(method, "", nil)
						if err != nil {
							t.Fatal(err)
						}
						c.authFunc(req)
						w := httptest.NewRecorder()

						sut.ServeHTTP(w, req)
						res := w.Result()

						if err = assert.Equal(
							http.StatusUnauthorized, res.StatusCode,
						); err != nil {
							t.Error(err)
						}

						if err = assert.Equal(
							"Bearer", res.Header.Values("WWW-Authenticate")[0],
						); err != nil {
							t.Error(err)
						}
					})
				}
			})
		}
	})

	t.Run("POST", func(t *testing.T) {
		for _, c := range []struct {
			name           string
			authFunc       func(*http.Request)
			boardName      string
			wantStatusCode int
			assertFunc     func(*testing.T, *httptest.ResponseRecorder)
		}{
			{
				name: "EmptyBoardName",
				authFunc: addBearerAuth(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2I" +
						"xMjMifQ.Y8_6K50EHUEJlJf4X21fNCFhYWhVIqN3Tw1niz8XwZc",
				),
				boardName:      "",
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assertOnErrMsg("Board name cannot be empty."),
			},
			{
				name: "TooLongBoardName",
				authFunc: addBearerAuth(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2I" +
						"xMjMifQ.Y8_6K50EHUEJlJf4X21fNCFhYWhVIqN3Tw1niz8XwZc",
				),
				boardName:      "A Board Whose Name Is Just Too Long!",
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg(
					"Board name cannot be longer than 35 characters.",
				)},
			{
				name: "TooManyBoards",
				authFunc: addBearerAuth(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2I" +
						"xMjMifQ.Y8_6K50EHUEJlJf4X21fNCFhYWhVIqN3Tw1niz8XwZc",
				),
				boardName:      "bob123's new board",
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg(
					"You have already created the maximum amount of boards " +
						"allowed per user. Please delete one of your boards " +
						"to create a new one.",
				)},
			{
				name: "Success",
				authFunc: addBearerAuth(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2I" +
						"xMjQifQ.LqENrj9APUHgQ3X0HRN6-IFMIg6nyo0_n74KfoxA0qI",
				),
				boardName:      "bob124's new board",
				wantStatusCode: http.StatusOK,
				assertFunc: func(*testing.T, *httptest.ResponseRecorder) {
					var boardCount int
					err := dbConn.QueryRow(
						"SELECT COUNT(*) FROM app.user_board " +
							"WHERE userID = 'bob124' AND isAdmin = TRUE",
					).Scan(&boardCount)
					if err != nil {
						t.Error(err)
					}
					if err = assert.Equal(1, boardCount); err != nil {
						t.Error(err)
					}
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				reqBodyBytes, err := json.Marshal(boardAPI.POSTReqBody{
					Name: c.boardName,
				})
				if err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest(
					http.MethodPost, "", bytes.NewReader(reqBodyBytes),
				)
				if err != nil {
					t.Fatal(err)
				}
				c.authFunc(req)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				if err = assert.Equal(
					c.wantStatusCode, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				c.assertFunc(t, w)
			})
		}
	})

	t.Run("DELETE", func(t *testing.T) {
		for _, c := range []struct {
			name           string
			id             string
			wantStatusCode int
			assertFunc     func(*testing.T)
		}{
			{
				name:           "EmptyID",
				id:             "",
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "NonIntID",
				id:             "qwerty",
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "UserBoardNotFound",
				id:             "123",
				wantStatusCode: http.StatusUnauthorized,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "UserNotAdmin",
				id:             "4",
				wantStatusCode: http.StatusUnauthorized,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "Success",
				id:             "1",
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T) {
					var boardID int
					err := dbConn.QueryRow(
						"SELECT boardID FROM app.user_board WHERE boardID = 1",
					).Scan(&boardID)
					if err != sql.ErrNoRows {
						t.Error("user_board row was not deleted")
					}
					err = dbConn.QueryRow(
						"SELECT id FROM app.board WHERE id = 1",
					).Scan(&boardID)
					if err != sql.ErrNoRows {
						t.Error("board row was not deleted")
					}
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				req, err := http.NewRequest(
					http.MethodDelete, "?id="+c.id, nil,
				)
				if err != nil {
					t.Fatal(err)
				}
				addBearerAuth(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2I" +
						"xMjMifQ.Y8_6K50EHUEJlJf4X21fNCFhYWhVIqN3Tw1niz8XwZc",
				)(req)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				if err = assert.Equal(
					c.wantStatusCode, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				c.assertFunc(t)
			})
		}
	})
}
