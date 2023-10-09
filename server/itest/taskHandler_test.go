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
				titleValidator,
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
			task           map[string]any
			wantStatusCode int
			assertFunc     func(*testing.T, *http.Response, string)
		}{
			{
				name: "TaskTitleEmpty",
				task: map[string]any{
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
				task: map[string]any{
					"title": "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqwe" +
						"asd",
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
				task: map[string]any{
					"title":       "Some Task",
					"description": "",
					"column":      0,
					"subtasks":    []string{""},
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Subtask title cannot be empty."),
			},
			{
				name: "SubtaskTitleTooLong",
				task: map[string]any{
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
				task: map[string]any{
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
				task: map[string]any{
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
				task: map[string]any{
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
				task: map[string]any{
					"title":       "Some Task",
					"description": "Do something. Then, do something else.",
					"column":      10,
					"subtasks":    []string{"Some Subtask", "Some Other Subtask"},
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
				task, err := json.Marshal(c.task)
				if err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest(
					http.MethodPost, "", bytes.NewReader(task),
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
			name       string
			taskTitle  string
			wantErrMsg string
		}{
			{
				name:       "TaskTitleEmpty",
				taskTitle:  "",
				wantErrMsg: "Task title cannot be empty.",
			},
			{
				name: "TaskTitleTooLong",
				taskTitle: "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqwe" +
					"asd",
				wantErrMsg: "Task title cannot be longer than 50 characters.",
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				wantStatusCode := http.StatusBadRequest

				task, err := json.Marshal(map[string]any{
					"title":       c.taskTitle,
					"description": "",
					"column":      0,
					"subtasks":    []string{},
				})
				if err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest(
					http.MethodPatch, "", bytes.NewReader(task),
				)
				if err != nil {
					t.Fatal(err)
				}

				addBearerAuth(jwtBob123)(req)

				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				if err = assert.Equal(
					wantStatusCode, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				assert.OnResErr(c.wantErrMsg)(t, res, "")
			})
		}
	})
}
