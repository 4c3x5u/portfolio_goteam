//go:build itest

package itest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	boardAPI "github.com/kxplxn/goteam/server/api/board"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/auth"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	teamTable "github.com/kxplxn/goteam/server/dbaccess/team"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestBoardHandler tests the http.Handler for the board API route and asserts
// that it behaves correctly during various execution paths.
func TestBoardHandler(t *testing.T) {
	// Create board API handler.
	userSelector := userTable.NewSelector(db)
	nameValidator := boardAPI.NewNameValidator()
	idValidator := boardAPI.NewIDValidator()
	boardSelector := boardTable.NewSelector(db)
	log := pkgLog.New()
	sut := api.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodGet: boardAPI.NewGETHandler(
				userSelector,
				idValidator,
				boardTable.NewRecursiveSelector(db),
				teamTable.NewSelector(db),
				userTable.NewSelectorByTeamID(db),
				boardTable.NewSelectorByTeamID(db),
				log,
			),
			http.MethodPost: boardAPI.NewPOSTHandler(
				userSelector,
				nameValidator,
				boardTable.NewCounter(db),
				boardTable.NewInserter(db),
				log,
			),
			http.MethodDelete: boardAPI.NewDELETEHandler(
				userSelector,
				idValidator,
				boardSelector,
				boardTable.NewDeleter(db),
				log,
			),
			http.MethodPatch: boardAPI.NewPATCHHandler(
				userSelector,
				idValidator,
				nameValidator,
				boardSelector,
				boardTable.NewUpdater(db),
				log,
			),
		},
	)

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
					http.MethodGet,
					http.MethodPost,
					http.MethodDelete,
					http.MethodPatch,
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

	t.Run("GET", func(t *testing.T) {
		for _, c := range []struct {
			name           string
			authFunc       func(*http.Request)
			boardID        string
			wantStatusCode int
			assertFunc     func(*testing.T, *http.Response, string)
		}{
			{
				name:           "InvalidID",
				authFunc:       addBearerAuth(jwtTeam1Member),
				boardID:        "foo",
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     func(*testing.T, *http.Response, string) {},
			},
			{
				name:           "NotFound",
				authFunc:       addBearerAuth(jwtTeam1Member),
				boardID:        "1001",
				wantStatusCode: http.StatusNotFound,
				assertFunc:     func(*testing.T, *http.Response, string) {},
			},
			{
				name:           "WrongTeam",
				authFunc:       addBearerAuth(jwtTeam2Member),
				boardID:        "2",
				wantStatusCode: http.StatusForbidden,
				assertFunc:     func(*testing.T, *http.Response, string) {},
			},
			{
				name:           "OK",
				authFunc:       addBearerAuth(jwtTeam1Member),
				boardID:        "2",
				wantStatusCode: http.StatusOK,
				assertFunc: func(*testing.T, *http.Response, string) {
					// TODO: assert on response body
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				req, err := http.NewRequest(
					http.MethodGet, "?id="+c.boardID, nil,
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

				c.assertFunc(t, res, "")
			})
		}
	})

	t.Run("POST", func(t *testing.T) {
		for _, c := range []struct {
			name           string
			authFunc       func(*http.Request)
			boardName      string
			wantStatusCode int
			assertFunc     func(*testing.T, *http.Response, string)
		}{
			{
				name:           "NotAdmin",
				authFunc:       addBearerAuth(jwtTeam2Member),
				boardName:      "Team 2 Board 2",
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can create boards.",
				),
			},
			{
				name:           "EmptyBoardName",
				authFunc:       addBearerAuth(jwtTeam1Admin),
				boardName:      "",
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Board name cannot be empty."),
			},
			{
				name:           "TooLongBoardName",
				authFunc:       addBearerAuth(jwtTeam1Admin),
				boardName:      "A Board Whose Name Is Just Too Long!",
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Board name cannot be longer than 35 characters.",
				),
			},
			{
				name:           "TooManyBoards",
				authFunc:       addBearerAuth(jwtTeam1Admin),
				boardName:      "bob123's new board",
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"You have already created the maximum amount of boards " +
						"allowed per user. Please delete one of your boards " +
						"to create a new one.",
				),
			},
			{
				name:           "Success",
				authFunc:       addBearerAuth(jwtTeam2Admin),
				boardName:      "Team 2 Board 2",
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ string) {
					// Assert that bob124 is assigned to the board as admin.
					var count int
					err := db.QueryRow(
						"SELECT COUNT(*) boardID FROM app.board " +
							"WHERE teamID = 2",
					).Scan(&count)
					if err != nil {
						t.Error(err)
					}
					if err = assert.Equal(1, count); err != nil {
						t.Error(err)
					}

					var boardID int
					err = db.QueryRow(
						"SELECT id FROM app.board WHERE teamID = 2",
					).Scan(&boardID)
					if err != nil {
						t.Error(err)
					}

					// Assert that 4 columns are created for this board.
					for order := 1; order < 5; order++ {
						var columnCount int
						err = db.QueryRow(
							`SELECT COUNT(*) FROM app."column" `+
								`WHERE boardID = $1 AND "order" = $2`,
							boardID,
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
				c.assertFunc(t, res, "")
			})
		}
	})

	t.Run("DELETE", func(t *testing.T) {
		for _, c := range []struct {
			name           string
			id             string
			authFunc       func(*http.Request)
			wantStatusCode int
			assertFunc     func(*testing.T)
		}{
			{
				name:           "NotAdmin",
				id:             "1",
				authFunc:       addBearerAuth(jwtTeam1Member),
				wantStatusCode: http.StatusForbidden,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "EmptyID",
				id:             "",
				authFunc:       addBearerAuth(jwtTeam3Admin),
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "NonIntID",
				id:             "qwerty",
				authFunc:       addBearerAuth(jwtTeam3Admin),
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "Success",
				id:             "4",
				authFunc:       addBearerAuth(jwtTeam3Admin),
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T) {
					// Assert that all user_board records with this board ID are
					// deleted.
					// Assert that the board record with this ID is deleted.
					var boardCount int
					err := db.QueryRow(
						"SELECT COUNT(*) FROM app.board WHERE teamID = 3",
					).Scan(&boardCount)
					if err != nil {
						t.Fatal(err)
					}
					if err = assert.Equal(0, boardCount); err != nil {
						t.Error(err)
					}

					var count int
					err = db.QueryRow(
						"SELECT COUNT(*) FROM app.board WHERE id = 4",
					).Scan(&count)
					if err != nil {
						t.Error(err)
					}
					if err = assert.Equal(0, count); err != nil {
						t.Error(err)
					}

					// Assert that all column records with this board ID are
					// deleted.
					var columnCount int
					err = db.QueryRow(
						`SELECT COUNT(*) FROM app."column" WHERE boardID = 4`,
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
				c.assertFunc(t)
			})
		}
	})

	t.Run("PATCH", func(t *testing.T) {
		for _, c := range []struct {
			name       string
			id         string
			boardName  string
			authFunc   func(*http.Request)
			statusCode int
			assertFunc func(*testing.T, *http.Response, string)
		}{
			{
				name:       "NotAdmin",
				id:         "1",
				boardName:  "New Board Name",
				authFunc:   addBearerAuth(jwtTeam1Member),
				statusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can edit the board.",
				),
			},
			{
				name:       "IDEmpty",
				id:         "",
				boardName:  "",
				authFunc:   addBearerAuth(jwtTeam1Admin),
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board ID cannot be empty."),
			},
			{
				name:       "IDNotInt",
				id:         "A",
				boardName:  "",
				authFunc:   addBearerAuth(jwtTeam1Admin),
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board ID must be an integer."),
			},
			{
				name:       "BoardNameEmpty",
				id:         "2",
				boardName:  "",
				authFunc:   addBearerAuth(jwtTeam1Admin),
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board name cannot be empty."),
			},
			{
				name:       "BoardNameTooLong",
				id:         "2",
				boardName:  "A Board Whose Name Is Just Too Long!",
				authFunc:   addBearerAuth(jwtTeam1Admin),
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Board name cannot be longer than 35 characters.",
				),
			},
			{
				name:       "BoardNotFound",
				id:         "1001",
				boardName:  "New Board Name",
				authFunc:   addBearerAuth(jwtTeam1Admin),
				statusCode: http.StatusNotFound,
				assertFunc: assert.OnResErr("Board not found."),
			},
			{
				name:       "Success",
				id:         "2",
				boardName:  "New Board Name",
				authFunc:   addBearerAuth(jwtTeam1Admin),
				statusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ string) {
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

				c.assertFunc(t, res, "")
			})
		}
	})
}
