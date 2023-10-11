//go:build itest

package itest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/api"
	taskAPI "server/api/task"
	"server/assert"
	"server/auth"
	columnTable "server/dbaccess/column"
	taskTable "server/dbaccess/task"
	userboardTable "server/dbaccess/userboard"
	pkgLog "server/log"
)

// TestTaskHandler tests the http.Handler for the task API route and asserts
// that it behaves correctly during various execution paths.
func TestTaskHandler(t *testing.T) {
	idValidator := taskAPI.NewIDValidator()
	titleValidator := taskAPI.NewTitleValidator()
	log := pkgLog.New()
	sut := api.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodPost: taskAPI.NewPOSTHandler(
				titleValidator,
				titleValidator,
				columnTable.NewSelector(db),
				userboardTable.NewSelector(db),
				taskTable.NewInserter(db),
				log,
			),
			http.MethodPatch: taskAPI.NewPATCHHandler(
				idValidator,
				titleValidator,
				titleValidator,
				taskTable.NewSelector(db),
				columnTable.NewSelector(db),
				userboardTable.NewSelector(db),
				taskTable.NewUpdater(db),
				log,
			),
			http.MethodDelete: taskAPI.NewDELETEHandler(idValidator, log),
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
				t.Run(http.MethodPost, func(t *testing.T) {
					req, err := http.NewRequest(http.MethodPost, "", nil)
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
			})
		}
	})

	t.Run(http.MethodPost, func(t *testing.T) {
		for _, c := range []struct {
			name           string
			reqBody        map[string]any
			wantStatusCode int
			assertFunc     func(*testing.T, *http.Response, string)
		}{
			{
				name: "TaskTitleEmpty",
				reqBody: map[string]any{
					"title":       "",
					"description": "",
					"column":      0,
					"subtasks":    []string{},
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Task title cannot be empty."),
			},
			{
				name: "TaskTitleTooLong",
				reqBody: map[string]any{
					"title": "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasd" +
						"qweasd",
					"description": "",
					"column":      0,
					"subtasks":    []string{},
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Task title cannot be longer than 50 characters.",
				),
			},
			{
				name: "SubtaskTitleEmpty",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "",
					"column":      0,
					"subtasks":    []string{""},
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Subtask title cannot be empty.",
				),
			},
			{
				name: "SubtaskTitleTooLong",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "",
					"column":      0,
					"subtasks": []string{
						"asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd",
					},
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Subtask title cannot be longer than 50 characters.",
				),
			},
			{
				name: "ColumnNotFound",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "",
					"column":      1001,
					"subtasks":    []string{"Some Subtask"},
				},
				wantStatusCode: http.StatusNotFound,
				assertFunc:     assert.OnResErr("Column not found."),
			},
			{
				name: "NoAccess",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "",
					"column":      8,
					"subtasks":    []string{"Some Subtask"},
				},
				wantStatusCode: http.StatusUnauthorized,
				assertFunc: assert.OnResErr(
					"You do not have access to this board.",
				),
			},
			{
				name: "NotAdmin",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "",
					"column":      9,
					"subtasks":    []string{"Some Subtask"},
				},
				wantStatusCode: http.StatusUnauthorized,
				assertFunc: assert.OnResErr(
					"Only board admins can create tasks.",
				),
			},
			{
				name: "Success",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "Do something. Then, do something else.",
					"column":      10,
					"subtasks": []string{
						"Some Subtask", "Some Other Subtask",
					},
				},
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ string) {
					// A task with the order of 1 and 2 already exists in the given
					// column. Therefore, the order of the newly created task must
					// be 3.
					var taskID, taskOrder int
					if err := db.QueryRow(
						`SELECT id, "order" FROM app.task `+
							`WHERE columnID = $1 AND title = $2`,
						10,
						"Some Task",
					).Scan(&taskID, &taskOrder); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(3, taskOrder); err != nil {
						t.Error(err)
					}

					// The order of the subtasks must be set respective to their
					// sequential order.
					for i, subtaskTitle := range []string{
						"Some Subtask", "Some Other Subtask",
					} {
						wantOrder := i + 1
						var subtaskOrder int
						if err := db.QueryRow(
							`SELECT "order" FROM app.subtask `+
								`WHERE taskID = $1 AND title = $2`,
							taskID,
							subtaskTitle,
						).Scan(&subtaskOrder); err != nil {
							t.Error(err)
						}
						if err := assert.Equal(
							wantOrder, subtaskOrder,
						); err != nil {
							t.Error(err)
						}
					}
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				reqBodyBytes, err := json.Marshal(c.reqBody)
				if err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest(
					http.MethodPost, "", bytes.NewReader(reqBodyBytes),
				)
				if err != nil {
					t.Fatal(err)
				}

				addBearerAuth(jwtBob123)(req)

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

	t.Run(http.MethodPatch, func(t *testing.T) {
		for _, c := range []struct {
			name           string
			taskID         string
			reqBody        map[string]any
			wantStatusCode int
			assertFunc     func(*testing.T, *http.Response, string)
		}{
			{
				name:   "TaskIDEmpty",
				taskID: "",
				reqBody: map[string]any{
					"title":       "",
					"description": "",
					"subtasks":    []map[string]any{},
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Task ID cannot be empty."),
			},
			{
				name:   "TaskIDNotInt",
				taskID: "A",
				reqBody: map[string]any{
					"title":       "",
					"description": "",
					"subtasks":    []map[string]any{},
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Task ID must be an integer."),
			},
			{
				name:   "TaskTitleEmpty",
				taskID: "0",
				reqBody: map[string]any{
					"title":       "",
					"description": "",
					"subtasks":    []map[string]any{},
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Task title cannot be empty."),
			},
			{
				name:   "TaskTitleTooLong",
				taskID: "0",
				reqBody: map[string]any{
					"title": "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasd" +
						"qweasd",
					"description": "",
					"subtasks":    []map[string]any{},
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Task title cannot be longer than 50 characters.",
				),
			},
			{
				name:   "SubtaskTitleEmpty",
				taskID: "0",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "",
					"subtasks": []map[string]any{
						{"title": ""},
					},
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Subtask title cannot be empty.",
				),
			},
			{
				name:   "SubtaskTitleTooLong",
				taskID: "0",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "",
					"subtasks": []map[string]any{{
						"title": "asdqweasdqweasdqweasdqweasdqweasdqweasdqwea" +
							"sdqweasd",
					}},
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Subtask title cannot be longer than 50 characters.",
				),
			},
			{
				name:   "TaskNotFound",
				taskID: "1001",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "",
					"subtasks":    []map[string]any{{"title": "Some Subtask"}},
				},
				wantStatusCode: http.StatusNotFound,
				assertFunc:     assert.OnResErr("Task not found."),
			},
			{
				name:   "SourceNoAccess",
				taskID: "7",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "",
					"subtasks": []map[string]any{{
						"title": "Some Subtask",
					}},
				},
				wantStatusCode: http.StatusUnauthorized,
				assertFunc: assert.OnResErr(
					"You do not have access to this board.",
				),
			},
			{
				name:   "SourceNotAdmin",
				taskID: "8",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "",
					"subtasks": []map[string]any{{
						"title": "Some Subtask",
					}},
				},
				wantStatusCode: http.StatusUnauthorized,
				assertFunc: assert.OnResErr(
					"Only board admins can edit tasks.",
				),
			},
			{
				name:   "Success",
				taskID: "9",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "Some Description",
					"subtasks": []map[string]any{
						{
							"title": "Some Subtask",
							"order": 1,
							"done":  false,
						},
						{
							"title": "Some Other Subtask",
							"order": 2,
							"done":  true,
						},
					},
				},
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ string) {
					var (
						title       string
						description *string
						columnID    int
						order       int
					)
					if err := db.QueryRow(
						`SELECT title, description, columnID, "order" `+
							`FROM app.task WHERE id = 9`,
					).Scan(&title, &description, &columnID, &order); err != nil {
						t.Error(err)
					}

					if err := assert.Equal("Some Task", title); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(
						"Some Description", *description,
					); err != nil {
						t.Error(err)
					}

					rows, err := db.Query(
						`SELECT title, "order", isDone FROM app.subtask WHERE taskID = 9`,
					)
					if err != nil {
						t.Fatal(err)
					}
					var subtasks []taskTable.Subtask
					for rows.Next() {
						var subtask taskTable.Subtask
						if err := rows.Scan(
							&subtask.Title, &subtask.Order, &subtask.IsDone,
						); err != nil {
							t.Fatal(err)
						}
						subtasks = append(subtasks, subtask)
					}
					if err = assert.Equal(2, len(subtasks)); err != nil {
						t.Error(err)
					}
					if err = assert.Equal("Some Subtask", subtasks[0].Title); err != nil {
						t.Error(err)
					}
					if err = assert.Equal(1, subtasks[0].Order); err != nil {
						t.Error(err)
					}
					if err = assert.Equal(false, subtasks[0].IsDone); err != nil {
						t.Error(err)
					}
					if err = assert.Equal("Some Other Subtask", subtasks[1].Title); err != nil {
						t.Error(err)
					}
					if err = assert.Equal(2, subtasks[1].Order); err != nil {
						t.Error(err)
					}
					if err = assert.Equal(true, subtasks[1].IsDone); err != nil {
						t.Error(err)
					}
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				reqBodyBytes, err := json.Marshal(c.reqBody)
				if err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest(
					http.MethodPatch,
					"?id="+c.taskID,
					bytes.NewReader(reqBodyBytes),
				)
				if err != nil {
					t.Fatal(err)
				}

				addBearerAuth(jwtBob123)(req)

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

	t.Run(http.MethodDelete, func(t *testing.T) {
		t.Run("IDEmpty", func(t *testing.T) {
			wantStatusCode := http.StatusBadRequest
			wantErrMsg := "Task ID cannot be empty."

			r, err := http.NewRequest(http.MethodDelete, "?id=", nil)
			if err != nil {
				t.Fatal(err)
			}
			addBearerAuth(jwtBob123)(r)
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, r)
			res := w.Result()

			if err = assert.Equal(wantStatusCode, res.StatusCode); err != nil {
				t.Error(err)
			}

			assert.OnResErr(wantErrMsg)(t, res, "")
		})

		t.Run("IDNotInt", func(t *testing.T) {
			id := "A"
			wantStatusCode := http.StatusBadRequest
			wantErrMsg := "Task ID must be an integer."

			r, err := http.NewRequest(http.MethodDelete, "?id="+id, nil)
			if err != nil {
				t.Fatal(err)
			}
			addBearerAuth(jwtBob123)(r)
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, r)
			res := w.Result()

			if err = assert.Equal(wantStatusCode, res.StatusCode); err != nil {
				t.Error(err)
			}

			assert.OnResErr(wantErrMsg)(t, res, "")
		})
	})
}
