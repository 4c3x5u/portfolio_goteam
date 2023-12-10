//go:build itest

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/internal/api"
	taskAPI "github.com/kxplxn/goteam/internal/api/task"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/auth"
	dynamoTaskTable "github.com/kxplxn/goteam/pkg/db/task"
	boardTable "github.com/kxplxn/goteam/pkg/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/pkg/dbaccess/column"
	taskTable "github.com/kxplxn/goteam/pkg/dbaccess/task"
	userTable "github.com/kxplxn/goteam/pkg/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// TestTaskHandler tests the http.Handler for the task API route and asserts
// that it behaves correctly during various execution paths.
func TestTaskHandler(t *testing.T) {
	idValidator := taskAPI.NewIDValidator()
	titleValidator := taskAPI.NewTitleValidator()
	taskSelector := taskTable.NewSelector(db)
	columnSelector := columnTable.NewSelector(db)
	boardSelector := boardTable.NewSelector(db)
	userSelector := userTable.NewSelector(db)
	log := pkgLog.New()

	sut := api.NewHandler(
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodPost: taskAPI.NewPostHandler(
				token.DecodeAuth,
				token.DecodeState,
				titleValidator,
				titleValidator,
				taskAPI.NewColNoValidator(),
				dynamoTaskTable.NewInserter(svcDynamo),
				token.EncodeState,
				log,
			),
			http.MethodPatch: taskAPI.NewPATCHHandler(
				userSelector,
				idValidator,
				titleValidator,
				titleValidator,
				taskSelector,
				columnSelector,
				boardSelector,
				taskTable.NewUpdater(db),
				log,
			),
			http.MethodDelete: taskAPI.NewDeleteHandler(
				token.DecodeAuth,
				token.DecodeState,
				dynamoTaskTable.NewDeleter(svcDynamo),
				token.EncodeState,
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
				for _, httpMethod := range []string{
					http.MethodPost, http.MethodPatch, http.MethodDelete,
				} {
					t.Run(httpMethod, func(t *testing.T) {
						req, err := http.NewRequest(httpMethod, "", nil)
						if err != nil {
							t.Fatal(err)
						}
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

	t.Run("POST", func(t *testing.T) {
		for _, c := range []struct {
			name           string
			reqBody        string
			authFunc       func(*http.Request)
			wantStatusCode int
			assertFunc     func(*testing.T, *http.Response, string)
		}{
			{
				name:           "NoAuth",
				reqBody:        `{}`,
				authFunc:       func(*http.Request) {},
				wantStatusCode: http.StatusUnauthorized,
				assertFunc: assert.OnResErr(
					"Auth token not found.",
				),
			},
			{
				name:           "InvalidAuth",
				reqBody:        `{}`,
				authFunc:       addCookieAuth("asdfkjldfs"),
				wantStatusCode: http.StatusUnauthorized,
				assertFunc:     assert.OnResErr("Invalid auth token."),
			},
			{
				name:           "NotAdmin",
				reqBody:        `{}`,
				authFunc:       addCookieAuth(tkTeam1Member),
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can create tasks.",
				),
			},
			{
				name:           "NoState",
				reqBody:        `{}`,
				authFunc:       addCookieAuth(tkTeam1Admin),
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"State token not found.",
				),
			},
			{
				name:    "InvalidState",
				reqBody: `{}`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState("ksadjfhaskdf")(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Invalid state token."),
			},
			{
				name: "ColNoOutOfBounds",
				reqBody: `{
                    "column": 5,
                    "board":  "91536664-9749-4dbb-a470-6e52aa353ae4"
                }`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Column number out of bounds.",
				),
			},
			{
				name:    "NoAccess",
				reqBody: `{}`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"You do not have access to this board.",
				),
			},
			{
				name: "TaskTitleEmpty",
				reqBody: `{
                    "board":  "91536664-9749-4dbb-a470-6e52aa353ae4",
                    "column": 1,
					"title":  ""
				}`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Task title cannot be empty."),
			},
			{
				name: "TaskTitleTooLong",
				reqBody: `{
                    "board":  "91536664-9749-4dbb-a470-6e52aa353ae4",
                    "column": 1,
					"title":  "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd"
				}`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Task title cannot be longer than 50 characters.",
				),
			},
			{
				name: "SubtaskTitleEmpty",
				reqBody: `{
                    "board":    "91536664-9749-4dbb-a470-6e52aa353ae4",
                    "column":   1,
					"title":    "Some Task",
					"subtasks": [""]
				}`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Subtask title cannot be empty.",
				),
			},
			{
				name: "SubtaskTitleTooLong",
				reqBody: `{
                    "board":    "91536664-9749-4dbb-a470-6e52aa353ae4",
                    "column":   1,
					"title":    "Some Task",
					"subtasks": [
						"asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd"
                    ]
				}`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Subtask title cannot be longer than 50 characters.",
				),
			},
			{
				name: "Success",
				reqBody: `{
                    "board":       "91536664-9749-4dbb-a470-6e52aa353ae4",
					"description": "Do something. Then, do something else.",
                    "column":      1,
					"title":       "Some Task",
					"subtasks": [
						"Some Subtask", "Some Other Subtask"
                    ]
				}`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ string) {
					keyEx := expression.Key("BoardID").Equal(expression.Value(
						"91536664-9749-4dbb-a470-6e52aa353ae4",
					))
					expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
					if err != nil {
						t.Fatal(err)
					}

					out, err := svcDynamo.Query(
						context.Background(),
						&dynamodb.QueryInput{
							TableName: &taskTableName,
							IndexName: aws.
								String("BoardID_index"),
							ExpressionAttributeNames:  expr.Names(),
							ExpressionAttributeValues: expr.Values(),
							KeyConditionExpression:    expr.KeyCondition(),
						},
					)
					if err != nil {
						t.Fatal("failed to get tasks:", err)
					}

					var taskFound bool
					for _, av := range out.Items {
						var task dynamoTaskTable.Task
						err = attributevalue.UnmarshalMap(av, &task)
						if err != nil {
							t.Fatal(err)
						}

						wantDescr := "Do something. Then, do something else."
						if task.Description == wantDescr &&
							task.ColumnNumber == 1 &&
							task.Title == "Some Task" &&
							task.Subtasks[0].Title == "Some Subtask" &&
							task.Subtasks[0].IsDone == false &&
							task.Subtasks[1].Title == "Some Other Subtask" &&
							task.Subtasks[1].IsDone == false {

							taskFound = true
						}
					}

					assert.True(t.Fatal, taskFound)
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPost, "/task", strings.NewReader(c.reqBody),
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

	t.Run("PATCH", func(t *testing.T) {
		for _, c := range []struct {
			name           string
			taskID         string
			reqBody        map[string]any
			authFunc       func(*http.Request)
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
				authFunc:       addCookieAuth(jwtTeam1Admin),
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
				authFunc:       addCookieAuth(jwtTeam1Admin),
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
				authFunc:       addCookieAuth(jwtTeam1Admin),
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
				authFunc:       addCookieAuth(jwtTeam1Admin),
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
				authFunc:       addCookieAuth(jwtTeam1Admin),
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
				authFunc:       addCookieAuth(jwtTeam1Admin),
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
				authFunc:       addCookieAuth(jwtTeam1Admin),
				wantStatusCode: http.StatusNotFound,
				assertFunc:     assert.OnResErr("Task not found."),
			},
			{
				name:   "WrongTeam",
				taskID: "7",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "",
					"subtasks": []map[string]any{{
						"title": "Some Subtask",
					}},
				},
				authFunc:       addCookieAuth(jwtTeam2Admin),
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"You do not have access to this board.",
				),
			},
			{
				name:   "NotAdmin",
				taskID: "8",
				reqBody: map[string]any{
					"title":       "Some Task",
					"description": "",
					"subtasks": []map[string]any{{
						"title": "Some Subtask",
					}},
				},
				authFunc:       addCookieAuth(jwtTeam1Member),
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can edit tasks.",
				),
			},
			{
				name:   "OK",
				taskID: "8",
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
				authFunc:       addCookieAuth(jwtTeam1Admin),
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
							`FROM app.task WHERE id = 8`,
					).Scan(&title, &description, &columnID, &order); err != nil {
						t.Error(err)
					}

					assert.Equal(t.Error, title, "Some Task")
					assert.Equal(t.Error, *description, "Some Description")

					rows, err := db.Query(
						`SELECT title, "order", isDone FROM app.subtask WHERE taskID = 8`,
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
					assert.Equal(t.Error, len(subtasks), 2)
					assert.Equal(t.Error, subtasks[0].Title, "Some Subtask")
					assert.Equal(t.Error, subtasks[0].Order, 1)
					assert.True(t.Error, !subtasks[0].IsDone)
					assert.Equal(t.Error, "Some Other Subtask", subtasks[1].Title)
					assert.Equal(t.Error, 2, subtasks[1].Order)
					assert.True(t.Error, subtasks[1].IsDone)
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

				c.authFunc(req)

				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

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
			assertFunc     func(*testing.T, *http.Response, string)
		}{
			{
				name:           "NoAuth",
				id:             "",
				authFunc:       func(*http.Request) {},
				wantStatusCode: http.StatusUnauthorized,
				assertFunc:     assert.OnResErr("Auth token not found."),
			},
			{
				name:           "InvalidAuth",
				id:             "",
				authFunc:       addCookieAuth("asdfasdf"),
				wantStatusCode: http.StatusUnauthorized,
				assertFunc:     assert.OnResErr("Invalid auth token."),
			},
			{
				name:           "NotAdmin",
				id:             "",
				authFunc:       addCookieAuth(tkTeam1Member),
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only board admins can delete tasks.",
				),
			},
			{
				name: "InvalidID",
				id:   "1001",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Invalid task ID."),
			},
			{
				name: "OK",
				id:   "c684a6a0-404d-46fa-9fa5-1497f9874567",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ string) {
					out, err := svcDynamo.GetItem(
						context.Background(),
						&dynamodb.GetItemInput{
							TableName: &taskTableName,
							Key: map[string]types.AttributeValue{
								"ID": &types.AttributeValueMemberS{
									Value: "c684a6a0-404d-46fa-9fa5-1497f9874" +
										"567",
								},
							},
						},
					)
					assert.Nil(t.Fatal, err)
					assert.Equal(t.Fatal, len(out.Item), 0)
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				r := httptest.NewRequest(http.MethodDelete, "/task?id="+c.id, nil)
				c.authFunc(r)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, r)
				res := w.Result()

				assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

				c.assertFunc(t, res, "")
			})
		}
	})
}
