//go:build itest

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/internal/api"
	boardAPI "github.com/kxplxn/goteam/internal/api/board"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/auth"
	boardTable "github.com/kxplxn/goteam/pkg/dbaccess/board"
	teamTable "github.com/kxplxn/goteam/pkg/dbaccess/team"
	userTable "github.com/kxplxn/goteam/pkg/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
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
						req := httptest.NewRequest(method, "/board", nil)
						c.authFunc(req)
						w := httptest.NewRecorder()

						sut.ServeHTTP(w, req)
						res := w.Result()

						assert.Equal(t.Error,
							res.StatusCode, http.StatusUnauthorized,
						)
						assert.Equal(t.Error,
							res.Header.Values("WWW-Authenticate")[0], "Bearer",
						)
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

					assert.Equal(t.Error, resp.User.Username, "team2Admin")

					assert.True(t.Error, resp.User.IsAdmin)

					assert.Equal(t.Error, resp.Team.ID, 2)
					assert.Equal(t.Error,
						resp.Team.InviteCode,
						"66ca0ddf-5f62-4713-bcc9-36cb0954eb7b",
					)
					assert.Equal(t.Error, len(resp.TeamMembers), 2)

					member := resp.TeamMembers[0]
					assert.Equal(t.Error, member.Username, "team2Admin")
					assert.True(t.Error, member.IsAdmin)

					// When ID is empty, a new board will be created for user.
					assert.Equal(t.Error, len(resp.Boards), 1)
					board := resp.Boards[0]
					assert.Equal(t.Error, board.ID, 5)
					assert.Equal(t.Error, board.Name, "New Board")
					assert.Equal(t.Error, resp.ActiveBoard.ID, 5)
					assert.Equal(t.Error, resp.ActiveBoard.Name, "New Board")

					for i, wantColumn := range []boardAPI.Column{
						{ID: 12, Order: 1},
						{ID: 13, Order: 2},
						{ID: 14, Order: 3},
						{ID: 15, Order: 4},
					} {
						column := resp.ActiveBoard.Columns[i]

						assert.Equal(t.Error, len(column.Tasks), 0)
						assert.Equal(t.Error, column.ID, wantColumn.ID)
						assert.Equal(t.Error, column.Order, wantColumn.Order)
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

					assert.Equal(t.Error, resp.User.Username, "team2Member")
					assert.True(t.Error, !resp.User.IsAdmin)

					assert.Equal(t.Error, resp.Team.ID, 2)
					assert.Equal(t.Error,
						resp.Team.InviteCode,
						"66ca0ddf-5f62-4713-bcc9-36cb0954eb7b",
					)

					assert.Equal(t.Error, len(resp.TeamMembers), 2)

					member := resp.TeamMembers[0]
					assert.Equal(t.Error, member.Username, "team2Admin")
					assert.True(t.Error, member.IsAdmin)

					// When ID is empty, a new board will be created for user.
					assert.Equal(t.Error, len(resp.Boards), 1)
					board := resp.Boards[0]
					assert.Equal(t.Error, board.ID, 5)
					assert.Equal(t.Error, board.Name, "New Board")
					assert.Equal(t.Error, resp.ActiveBoard.ID, 5)
					assert.Equal(t.Error, resp.ActiveBoard.Name, "New Board")

					for i, wantColumn := range []boardAPI.Column{
						{ID: 12, Order: 1},
						{ID: 13, Order: 2},
						{ID: 14, Order: 3},
						{ID: 15, Order: 4},
					} {
						column := resp.ActiveBoard.Columns[i]

						assert.Equal(t.Error, len(column.Tasks), 0)
						assert.Equal(t.Error, wantColumn.ID, column.ID)
						assert.Equal(t.Error, column.Order, wantColumn.Order)
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

					assert.Equal(t.Error, resp.User.Username, "team1Member")
					assert.True(t.Error, !resp.User.IsAdmin)

					assert.Equal(t.Error, resp.Team.ID, 1)
					assert.Equal(t.Error,
						resp.Team.InviteCode,
						"afeadc4a-68b0-4c33-9e83-4648d20ff26a",
					)

					for i, wantMember := range []boardAPI.TeamMember{
						{Username: "team1Admin", IsAdmin: true},
						{Username: "team1Member", IsAdmin: false},
					} {
						member := resp.TeamMembers[i]
						assert.Equal(t.Error,
							member.Username, wantMember.Username,
						)
						assert.Equal(t.Error,
							member.IsAdmin, wantMember.IsAdmin,
						)
					}

					for i, wantBoard := range []boardAPI.Board{
						{ID: 1, Name: "Team 1 Board 1"},
						{ID: 2, Name: "Team 1 Board 2"},
						{ID: 3, Name: "Team 1 Board 3"},
					} {
						board := resp.Boards[i]
						assert.Equal(t.Error, board.ID, wantBoard.ID)
						assert.Equal(t.Error, board.Name, wantBoard.Name)
					}

					assert.Equal(t.Error, resp.ActiveBoard.ID, 2)
					assert.Equal(t.Error,
						resp.ActiveBoard.Name, "Team 1 Board 2",
					)

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

						assert.Equal(t.Error, column.ID, wantColumn.ID)
						assert.Equal(t.Error, column.Order, wantColumn.Order)

						for j, wantTask := range wantColumn.Tasks {
							task := column.Tasks[j]

							assert.Equal(t.Error, task.ID, wantTask.ID)
							assert.Equal(t.Error, wantTask.Title, task.Title)
							assert.Equal(t.Error,
								task.Description, wantTask.Description,
							)
							assert.Equal(t.Error, task.Order, wantTask.Order)

							for k, wantSubtask := range wantTask.Subtasks {
								subtask := task.Subtasks[k]

								assert.Equal(t.Error,
									subtask.ID, wantSubtask.ID,
								)
								assert.Equal(t.Error,
									subtask.Title, wantSubtask.Title,
								)
								assert.Equal(t.Error,
									subtask.Order,
									wantSubtask.Order,
								)
								assert.Equal(t.Error,
									subtask.IsDone,
									wantSubtask.IsDone,
								)
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

				assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

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
					assert.Equal(t.Error, count, 1)

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
						assert.Equal(t.Error, columnCount, 1)
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

				assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

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
					assert.Equal(t.Error, boardCount, 0)

					var count int
					err = db.QueryRow(
						"SELECT COUNT(*) FROM app.board WHERE id = 4",
					).Scan(&count)
					if err != nil {
						t.Error(err)
					}
					assert.Equal(t.Error, count, 0)

					// Assert that all column records with this board ID are
					// deleted.
					var columnCount int
					err = db.QueryRow(
						`SELECT COUNT(*) FROM app."column" WHERE boardID = 4`,
					).Scan(&columnCount)
					if err != nil {
						t.Fatal(err)
					}
					assert.Equal(t.Error, columnCount, 0)

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
						assert.Equal(t.Error, count, 0)
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
						assert.Equal(t.Error, count, 0)
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

				assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

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

				assert.Equal(t.Error, res.StatusCode, c.statusCode)

				c.assertFunc(t, res, "")
			})
		}
	})
}
