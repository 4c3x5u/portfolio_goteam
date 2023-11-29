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
	boardInserter := boardTable.NewInserter(db)
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
				boardInserter,
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
				boardInserter,
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
			{name: "HeaderInvalid", authFunc: addCookieAuth("asdfasldfkjasd")},
		} {
			t.Run(c.name, func(t *testing.T) {
				for _, method := range []string{
					http.MethodGet,
					http.MethodPost,
					http.MethodDelete,
					http.MethodPatch,
				} {
					t.Run(method, func(t *testing.T) {
						req := httptest.NewRequest(method, "/", nil)
						c.authFunc(req)
						w := httptest.NewRecorder()

						sut.ServeHTTP(w, req)
						res := w.Result()

						if err := assert.Equal(
							http.StatusUnauthorized, res.StatusCode,
						); err != nil {
							t.Error(err)
						}

						if err := assert.Equal(
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
				authFunc:       addCookieAuth(jwtTeam1Member),
				boardID:        "foo",
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     func(*testing.T, *http.Response, string) {},
			},
			{
				name:           "NotFound",
				authFunc:       addCookieAuth(jwtTeam1Member),
				boardID:        "1001",
				wantStatusCode: http.StatusNotFound,
				assertFunc:     func(*testing.T, *http.Response, string) {},
			},
			{
				name:           "WrongTeam",
				authFunc:       addCookieAuth(jwtTeam2Member),
				boardID:        "2",
				wantStatusCode: http.StatusForbidden,
				assertFunc:     func(*testing.T, *http.Response, string) {},
			},
			{
				name:           "OKIDEmptyAdmin",
				authFunc:       addCookieAuth(jwtTeam2Admin),
				boardID:        "",
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, r *http.Response, _ string) {
					var resp boardAPI.GETResp
					if err := json.NewDecoder(r.Body).Decode(
						&resp,
					); err != nil {
						t.Fatal(err)
					}

					if err := assert.Equal(
						"team2Admin", resp.User.Username,
					); err != nil {
						t.Error(err)
					}
					if err := assert.True(resp.User.IsAdmin); err != nil {
						t.Error(err)
					}

					if err := assert.Equal(2, resp.Team.ID); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(
						"66ca0ddf-5f62-4713-bcc9-36cb0954eb7b",
						resp.Team.InviteCode,
					); err != nil {
						t.Error(err)
					}

					if err := assert.Equal(
						2, len(resp.TeamMembers),
					); err != nil {
						t.Error(err)
					}
					member := resp.TeamMembers[0]
					if err := assert.Equal(
						member.Username, "team2Admin",
					); err != nil {
						t.Error(err)
					}
					if err := assert.True(member.IsAdmin); err != nil {
						t.Error(err)
					}

					// When ID is empty, a new board will be created for user.
					if err := assert.Equal(
						1, len(resp.Boards),
					); err != nil {
						t.Error(err)
					}
					board := resp.Boards[0]
					if err := assert.Equal(5, board.ID); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(
						"New Board", board.Name,
					); err != nil {
						t.Error(err)
					}

					if err := assert.Equal(5, resp.ActiveBoard.ID); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(
						"New Board", resp.ActiveBoard.Name,
					); err != nil {
						t.Error(err)
					}

					for i, wantColumn := range []boardAPI.Column{
						{ID: 12, Order: 1},
						{ID: 13, Order: 2},
						{ID: 14, Order: 3},
						{ID: 15, Order: 4},
					} {
						column := resp.ActiveBoard.Columns[i]

						if err := assert.Equal(
							0, len(column.Tasks),
						); err != nil {
							t.Error(err)
						}
						if err := assert.Equal(
							wantColumn.ID, column.ID,
						); err != nil {
							t.Error(err)
						}
						if err := assert.Equal(
							wantColumn.Order, column.Order,
						); err != nil {
							t.Error(err)
						}
					}
				},
			},
			// FIXME: Depends on the previous test case's succes - bad.
			{
				name:           "OKIDEmptyMember",
				authFunc:       addCookieAuth(jwtTeam2Member),
				boardID:        "",
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, r *http.Response, _ string) {
					var resp boardAPI.GETResp
					if err := json.NewDecoder(r.Body).Decode(
						&resp,
					); err != nil {
						t.Fatal(err)
					}

					if err := assert.Equal(
						"team2Member", resp.User.Username,
					); err != nil {
						t.Error(err)
					}
					if err := assert.True(!resp.User.IsAdmin); err != nil {
						t.Error(err)
					}

					if err := assert.Equal(2, resp.Team.ID); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(
						"66ca0ddf-5f62-4713-bcc9-36cb0954eb7b",
						resp.Team.InviteCode,
					); err != nil {
						t.Error(err)
					}

					if err := assert.Equal(
						2, len(resp.TeamMembers),
					); err != nil {
						t.Error(err)
					}
					member := resp.TeamMembers[0]
					if err := assert.Equal(
						member.Username, "team2Admin",
					); err != nil {
						t.Error(err)
					}
					if err := assert.True(member.IsAdmin); err != nil {
						t.Error(err)
					}

					// When ID is empty, a new board will be created for user.
					if err := assert.Equal(
						1, len(resp.Boards),
					); err != nil {
						t.Error(err)
					}
					board := resp.Boards[0]
					if err := assert.Equal(5, board.ID); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(
						"New Board", board.Name,
					); err != nil {
						t.Error(err)
					}

					if err := assert.Equal(5, resp.ActiveBoard.ID); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(
						"New Board", resp.ActiveBoard.Name,
					); err != nil {
						t.Error(err)
					}

					for i, wantColumn := range []boardAPI.Column{
						{ID: 12, Order: 1},
						{ID: 13, Order: 2},
						{ID: 14, Order: 3},
						{ID: 15, Order: 4},
					} {
						column := resp.ActiveBoard.Columns[i]

						if err := assert.Equal(
							0, len(column.Tasks),
						); err != nil {
							t.Error(err)
						}
						if err := assert.Equal(
							wantColumn.ID, column.ID,
						); err != nil {
							t.Error(err)
						}
						if err := assert.Equal(
							wantColumn.Order, column.Order,
						); err != nil {
							t.Error(err)
						}
					}
				},
			},
			{
				name:           "OK",
				authFunc:       addCookieAuth(jwtTeam1Member),
				boardID:        "2",
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, r *http.Response, _ string) {
					var resp boardAPI.GETResp
					if err := json.NewDecoder(r.Body).Decode(
						&resp,
					); err != nil {
						t.Fatal(err)
					}

					if err := assert.Equal(
						"team1Member", resp.User.Username,
					); err != nil {
						t.Error(err)
					}
					if err := assert.True(!resp.User.IsAdmin); err != nil {
						t.Error(err)
					}

					if err := assert.Equal(1, resp.Team.ID); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(
						"afeadc4a-68b0-4c33-9e83-4648d20ff26a",
						resp.Team.InviteCode,
					); err != nil {
						t.Error(err)
					}

					for i, wantMember := range []boardAPI.TeamMember{
						{Username: "team1Admin", IsAdmin: true},
						{Username: "team1Member", IsAdmin: false},
					} {
						member := resp.TeamMembers[i]
						if err := assert.Equal(
							wantMember.Username,
							member.Username,
						); err != nil {
							t.Error(err)
						}
						if err := assert.Equal(
							wantMember.IsAdmin,
							member.IsAdmin,
						); err != nil {
							t.Error(err)
						}
					}

					for i, wantBoard := range []boardAPI.Board{
						{ID: 1, Name: "Team 1 Board 1"},
						{ID: 2, Name: "Team 1 Board 2"},
						{ID: 3, Name: "Team 1 Board 3"},
					} {
						board := resp.Boards[i]
						if err := assert.Equal(
							wantBoard.ID,
							board.ID,
						); err != nil {
							t.Error(err)
						}
						if err := assert.Equal(
							wantBoard.Name,
							board.Name,
						); err != nil {
							t.Error(err)
						}
					}

					if err := assert.Equal(2, resp.ActiveBoard.ID); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(
						"Team 1 Board 2",
						resp.ActiveBoard.Name,
					); err != nil {
						t.Error(err)
					}

					for i, wantColumn := range []boardAPI.Column{
						{ID: 8, Order: 1, Tasks: []boardAPI.Task{
							{
								ID:          10,
								Title:       "task 10",
								Description: "desc",
								Order:       1,
								Subtasks: []boardAPI.Subtask{
									{
										ID:     8,
										Title:  "subtask 8",
										Order:  1,
										IsDone: false,
									},
									{
										ID:     9,
										Title:  "subtask 9",
										Order:  2,
										IsDone: true,
									},
								},
							},
						}},
						{ID: 9, Order: 2, Tasks: []boardAPI.Task{}},
						{ID: 10, Order: 3, Tasks: []boardAPI.Task{
							{
								ID:          11,
								Title:       "task 11",
								Description: "",
								Order:       1,
							},
						}},
						{ID: 11, Order: 4, Tasks: []boardAPI.Task{}},
					} {
						column := resp.ActiveBoard.Columns[i]

						if err := assert.Equal(
							wantColumn.ID,
							column.ID,
						); err != nil {
							t.Error(err)
						}
						if err := assert.Equal(
							wantColumn.Order,
							column.Order,
						); err != nil {
							t.Error(err)
						}

						for j, wantTask := range wantColumn.Tasks {
							task := column.Tasks[j]

							if err := assert.Equal(
								wantTask.ID,
								task.ID,
							); err != nil {
								t.Error(err)
							}
							if err := assert.Equal(
								wantTask.Title,
								task.Title,
							); err != nil {
								t.Error(err)
							}
							if err := assert.Equal(
								wantTask.Description,
								task.Description,
							); err != nil {
								t.Error(err)
							}
							if err := assert.Equal(
								wantTask.Order,
								task.Order,
							); err != nil {
								t.Error(err)
							}

							for k, wantSubtask := range wantTask.Subtasks {
								subtask := task.Subtasks[k]

								if err := assert.Equal(
									wantSubtask.ID,
									subtask.ID,
								); err != nil {
									t.Error(err)
								}
								if err := assert.Equal(
									wantSubtask.Title,
									subtask.Title,
								); err != nil {
									t.Error(err)
								}
								if err := assert.Equal(
									wantSubtask.Order,
									subtask.Order,
								); err != nil {
									t.Error(err)
								}
								if err := assert.Equal(
									wantSubtask.IsDone,
									subtask.IsDone,
								); err != nil {
									t.Error(err)
								}
							}
						}
					}
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodGet, "/?id="+c.boardID, nil,
				)
				c.authFunc(req)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				if err := assert.Equal(
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
				authFunc:       addCookieAuth(jwtTeam2Member),
				boardName:      "Team 2 Board 2",
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can create boards.",
				),
			},
			{
				name:           "EmptyBoardName",
				authFunc:       addCookieAuth(jwtTeam1Admin),
				boardName:      "",
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Board name cannot be empty."),
			},
			{
				name:           "TooLongBoardName",
				authFunc:       addCookieAuth(jwtTeam1Admin),
				boardName:      "A Board Whose Name Is Just Too Long!",
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Board name cannot be longer than 35 characters.",
				),
			},
			{
				name:           "TooManyBoards",
				authFunc:       addCookieAuth(jwtTeam1Admin),
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
				authFunc:       addCookieAuth(jwtTeam4Admin),
				boardName:      "Team 4 Board 1",
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ string) {
					// Assert that bob124 is assigned to the board as admin.
					var count int
					err := db.QueryRow(
						"SELECT COUNT(*) boardID FROM app.board " +
							"WHERE teamID = 4",
					).Scan(&count)
					if err != nil {
						t.Error(err)
					}
					if err = assert.Equal(1, count); err != nil {
						t.Error(err)
					}

					var boardID int
					err = db.QueryRow(
						"SELECT id FROM app.board WHERE teamID = 4",
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
				req := httptest.NewRequest(
					http.MethodPost, "/", bytes.NewReader(reqBody),
				)
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
				authFunc:       addCookieAuth(jwtTeam1Member),
				wantStatusCode: http.StatusForbidden,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "EmptyID",
				id:             "",
				authFunc:       addCookieAuth(jwtTeam3Admin),
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "NonIntID",
				id:             "qwerty",
				authFunc:       addCookieAuth(jwtTeam3Admin),
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "Success",
				id:             "4",
				authFunc:       addCookieAuth(jwtTeam3Admin),
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
				req := httptest.NewRequest(
					http.MethodDelete, "/?id="+c.id, nil,
				)
				c.authFunc(req)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				if err := assert.Equal(
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
				authFunc:   addCookieAuth(jwtTeam1Member),
				statusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can edit the board.",
				),
			},
			{
				name:       "IDEmpty",
				id:         "",
				boardName:  "",
				authFunc:   addCookieAuth(jwtTeam1Admin),
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board ID cannot be empty."),
			},
			{
				name:       "IDNotInt",
				id:         "A",
				boardName:  "",
				authFunc:   addCookieAuth(jwtTeam1Admin),
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board ID must be an integer."),
			},
			{
				name:       "BoardNameEmpty",
				id:         "2",
				boardName:  "",
				authFunc:   addCookieAuth(jwtTeam1Admin),
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board name cannot be empty."),
			},
			{
				name:       "BoardNameTooLong",
				id:         "2",
				boardName:  "A Board Whose Name Is Just Too Long!",
				authFunc:   addCookieAuth(jwtTeam1Admin),
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Board name cannot be longer than 35 characters.",
				),
			},
			{
				name:       "BoardNotFound",
				id:         "1001",
				boardName:  "New Board Name",
				authFunc:   addCookieAuth(jwtTeam1Admin),
				statusCode: http.StatusNotFound,
				assertFunc: assert.OnResErr("Board not found."),
			},
			{
				name:       "Success",
				id:         "2",
				boardName:  "New Board Name",
				authFunc:   addCookieAuth(jwtTeam1Admin),
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
				req := httptest.NewRequest(
					http.MethodPatch, "/?id="+c.id, bytes.NewReader(reqBody),
				)
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
