//go:build itest

package itest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/api"
	"server/api/board"
	"server/assert"
	"server/auth"
	boardTable "server/dbaccess/board"
	userboardTable "server/dbaccess/userboard"
	pkgLog "server/log"
)

func TestBoard(t *testing.T) {
	// Create board API handler.
	log := pkgLog.New()
	sut := board.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodPost: board.NewPOSTHandler(
				board.NewNameValidator(),
				userboardTable.NewCounter(db),
				boardTable.NewInserter(db),
				log,
			),
			http.MethodDelete: board.NewDELETEHandler(
				board.NewIDValidator(),
				userboardTable.NewSelector(db),
				boardTable.NewDeleter(db),
				log,
			),
			http.MethodPatch: board.NewPATCHHandler(
				board.NewIDValidator(),
				board.NewNameValidator(),
				boardTable.NewSelector(db),
				userboardTable.NewSelector(db),
				boardTable.NewUpdater(db),
				log,
			),
		},
	)

	const (
		bob123AuthToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ" +
			"ib2IxMjMifQ.Y8_6K50EHUEJlJf4X21fNCFhYWhVIqN3Tw1niz8XwZc"
		bob124AuthToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ" +
			"ib2IxMjQifQ.LqENrj9APUHgQ3X0HRN6-IFMIg6nyo0_n74KfoxA0qI"
	)

	// used in various test cases to authenticate the request sent
	addBearerAuth := func(token string) func(*http.Request) {
		return func(req *http.Request) {
			req.Header.Add("Authorization", "Bearer "+token)
		}
	}

	// Used in status 400 error cases to assert on the error message.
	assertOnErrMsg := func(
		wantErrMsg string,
	) func(*testing.T, *httptest.ResponseRecorder) {
		return func(t *testing.T, w *httptest.ResponseRecorder) {
			resBody := board.ResBody{}
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
					http.MethodPost, http.MethodDelete, http.MethodPatch,
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

	t.Run(http.MethodPost, func(t *testing.T) {
		for _, c := range []struct {
			name           string
			authFunc       func(*http.Request)
			boardName      string
			wantStatusCode int
			assertFunc     func(*testing.T, *httptest.ResponseRecorder)
		}{
			{
				name:           "EmptyBoardName",
				authFunc:       addBearerAuth(bob123AuthToken),
				boardName:      "",
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assertOnErrMsg("Board name cannot be empty."),
			},
			{
				name:           "TooLongBoardName",
				authFunc:       addBearerAuth(bob123AuthToken),
				boardName:      "A Board Whose Name Is Just Too Long!",
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg(
					"Board name cannot be longer than 35 characters.",
				)},
			{
				name:           "TooManyBoards",
				authFunc:       addBearerAuth(bob123AuthToken),
				boardName:      "bob123's new board",
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg(
					"You have already created the maximum amount of boards " +
						"allowed per user. Please delete one of your boards " +
						"to create a new one.",
				)},
			{
				name:           "Success",
				authFunc:       addBearerAuth(bob124AuthToken),
				boardName:      "bob124's new board",
				wantStatusCode: http.StatusOK,
				assertFunc: func(*testing.T, *httptest.ResponseRecorder) {
					// Assert that bob124 is assigned to the board as admin.
					var boardCount int
					err := db.QueryRow(
						"SELECT COUNT(*) FROM app.user_board " +
							"WHERE boardID = 5 " +
							"AND username = 'bob124' " +
							"AND isAdmin = TRUE",
					).Scan(&boardCount)
					if err != nil {
						t.Error(err)
					}
					if err = assert.Equal(1, boardCount); err != nil {
						t.Error(err)
					}

					// Assert that 4 columns are created for this board.
					for order := 1; order < 5; order++ {
						var columnCount int
						err = db.QueryRow(
							`SELECT COUNT(*) FROM app."column" `+
								`WHERE boardID = 5 AND "order" = $1`,
							order,
						).Scan(&columnCount)
						if err != nil {
							t.Error(err)
						}
						if err = assert.Equal(1, columnCount); err != nil {
							t.Error(err)
						}
					}
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				reqBody, err := json.Marshal(map[string]string{
					"name": c.boardName,
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
				c.authFunc(req)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				if err = assert.Equal(
					c.wantStatusCode, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				// Run case-specific assertions.
				c.assertFunc(t, w)
			})
		}
	})

	t.Run(http.MethodDelete, func(t *testing.T) {
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
				wantStatusCode: http.StatusForbidden,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "UserNotAdmin",
				id:             "4",
				wantStatusCode: http.StatusForbidden,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "Success",
				id:             "1",
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T) {
					// Assert that all user_board records with this board ID are
					// deleted.
					var userBoardCount int
					err := db.QueryRow(
						"SELECT COUNT(*) FROM app.user_board WHERE boardID = 1",
					).Scan(&userBoardCount)
					if err != nil {
						t.Fatal(err)
					}
					if err = assert.Equal(0, userBoardCount); err != nil {
						t.Error(err)
					}

					// Assert that the board record with this ID is deleted.
					var boardCount int
					err = db.QueryRow(
						"SELECT COUNT(*) FROM app.board WHERE id = 1",
					).Scan(&boardCount)
					if err != nil {
						t.Fatal(err)
					}
					if err = assert.Equal(0, boardCount); err != nil {
						t.Error(err)
					}

					// Assert that all column records with this board ID are
					// deleted.
					var columnCount int
					err = db.QueryRow(
						`SELECT COUNT(*) FROM app."column" WHERE boardID = 1`,
					).Scan(&columnCount)
					if err != nil {
						t.Fatal(err)
					}
					if err = assert.Equal(0, columnCount); err != nil {
						t.Error(err)
					}

					// Assert that all task records associated with each
					// column record is deleted.
					for columnID := 1; columnID < 5; columnID++ {
						var count int
						err = db.QueryRow(
							`SELECT COUNT(*) FROM app.task WHERE columnID = $1`,
							columnID,
						).Scan(&count)
						if err != nil {
							t.Fatal(err)
						}
						if err = assert.Equal(0, count); err != nil {
							t.Error(err)
						}
					}

					// Assert that all subtask records associated with each
					// task record is deleted.
					for taskID := 1; taskID < 5; taskID++ {
						var count int
						err = db.QueryRow(
							"SELECT COUNT(*) FROM app.subtask "+
								"WHERE taskID = $1",
							taskID,
						).Scan(&count)
						if err != nil {
							t.Fatal(err)
						}
						if err = assert.Equal(0, count); err != nil {
							t.Error(err)
						}
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
				addBearerAuth(bob123AuthToken)(req)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				if err = assert.Equal(
					c.wantStatusCode, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				// Run case-specific assertions.
				c.assertFunc(t)
			})
		}
	})

	t.Run(http.MethodPatch, func(t *testing.T) {
		for _, c := range []struct {
			name       string
			id         string
			boardName  string
			authFunc   func(*http.Request)
			statusCode int
			assertFunc func(*testing.T, *httptest.ResponseRecorder)
		}{
			{
				name:       "IDEmpty",
				id:         "",
				boardName:  "",
				authFunc:   addBearerAuth(bob123AuthToken),
				statusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg("Board ID cannot be empty."),
			},
			{
				name:       "IDNotInt",
				id:         "A",
				boardName:  "",
				authFunc:   addBearerAuth(bob123AuthToken),
				statusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg("Board ID must be an integer."),
			},
			{
				name:       "BoardNameEmpty",
				id:         "2",
				boardName:  "",
				authFunc:   addBearerAuth(bob123AuthToken),
				statusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg("Board name cannot be empty."),
			},
			{
				name:       "BoardNameTooLong",
				id:         "2",
				boardName:  "A Board Whose Name Is Just Too Long!",
				authFunc:   addBearerAuth(bob123AuthToken),
				statusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg(
					"Board name cannot be longer than 35 characters.",
				),
			},
			{
				name:       "BoardNotFound",
				id:         "1001",
				boardName:  "New Board Name",
				authFunc:   addBearerAuth(bob123AuthToken),
				statusCode: http.StatusNotFound,
				assertFunc: assertOnErrMsg("Board not found."),
			},
			{
				name:       "UserBoardNotFound",
				id:         "3",
				boardName:  "New Board Name",
				authFunc:   addBearerAuth(bob124AuthToken),
				statusCode: http.StatusForbidden,
				assertFunc: assertOnErrMsg(
					"You do not have access to this board.",
				),
			},
			{
				name:       "UserNotAdmin",
				id:         "4",
				boardName:  "New Board Name",
				authFunc:   addBearerAuth(bob123AuthToken),
				statusCode: http.StatusForbidden,
				assertFunc: assertOnErrMsg(
					"Only board admins can edit the board.",
				),
			},
			{
				name:       "Success",
				id:         "2",
				boardName:  "New Board Name",
				authFunc:   addBearerAuth(bob123AuthToken),
				statusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *httptest.ResponseRecorder) {
					var boardName string
					err := db.QueryRow(
						"SELECT name FROM app.board WHERE id = 2",
					).Scan(&boardName)
					if err != nil {
						t.Fatal(err)
					}
					if boardName != "New Board Name" {
						t.Error("Board name was not updated.")
					}
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				reqBody, err := json.Marshal(map[string]string{
					"name": c.boardName,
				})
				if err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest(
					http.MethodPatch, "?id="+c.id, bytes.NewReader(reqBody),
				)
				if err != nil {
					t.Fatal(err)
				}
				c.authFunc(req)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				if err = assert.Equal(
					c.statusCode, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				c.assertFunc(t, w)
			})
		}
	})
}
