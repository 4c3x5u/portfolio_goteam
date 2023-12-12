//go:build itest

package api

import (
	"context"
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
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

func TestTaskAPI(t *testing.T) {
	titleValidator := taskAPI.NewTitleValidator()
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
			http.MethodPatch: taskAPI.NewPatchHandler(
				token.DecodeAuth,
				token.DecodeState,
				titleValidator,
				titleValidator,
				dynamoTaskTable.NewUpdater(svcDynamo),
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
						req := httptest.NewRequest(httpMethod, "/", nil)
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
			reqBody        string
			authFunc       func(*http.Request)
			wantStatusCode int
			assertFunc     func(*testing.T, *http.Response, string)
		}{
			{
				name:   "NotAdmin",
				taskID: "e0021a56-6a1e-4007-b773-395d3991fb7e",
				reqBody: `{
					"title":       "Some Task",
					"description": "",
					"subtasks":    [{"title": "Some Subtask"}]
				}`,
				authFunc:       addCookieAuth(tkTeam1Member),
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can edit tasks.",
				),
			},
			{
				name:   "InvalidTaskID",
				taskID: "",
				reqBody: `{
					"title":       "",
					"description": "",
					"subtasks":    []
				}`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Invalid task ID."),
			},
			{
				name:   "TaskTitleEmpty",
				taskID: "e0021a56-6a1e-4007-b773-395d3991fb7e",
				reqBody: `{
					"title":       "",
					"description": "",
					"subtasks":    []
				}`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Task title cannot be empty."),
			},
			{
				name:   "TaskTitleTooLong",
				taskID: "e0021a56-6a1e-4007-b773-395d3991fb7e",
				reqBody: `{
					"title": "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd",
					"description": "",
					"subtasks":    []
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
				name:   "SubtaskTitleEmpty",
				taskID: "e0021a56-6a1e-4007-b773-395d3991fb7e",
				reqBody: `{
					"title":       "Some Task",
					"description": "",
					"subtasks":    [{"title": ""}]
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
				name:   "SubtaskTitleTooLong",
				taskID: "e0021a56-6a1e-4007-b773-395d3991fb7e",
				reqBody: `{
					"title":       "Some Task",
					"description": "",
					"subtasks": [{
						"title": "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd"
					}]
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
				name:   "OK",
				taskID: "e0021a56-6a1e-4007-b773-395d3991fb7e",
				reqBody: `{
					"title":       "Some Task",
					"description": "Some Description",
					"subtasks": [
						{
							"title": "Some Subtask",
							"done":  false
						},
						{
							"title": "Some Other Subtask",
							"done":  true
						}
                    ]
				}`,
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
									Value: "e0021a56-6a1e-4007-b773-395d3991f" +
										"b7e",
								},
							},
						},
					)
					assert.Nil(t.Fatal, err)

					var task dynamoTaskTable.Task
					err = attributevalue.UnmarshalMap(out.Item, &task)
					assert.Nil(t.Fatal, err)

					assert.Equal(t.Error, task.Title, "Some Task")
					assert.Equal(t.Error, task.Description, "Some Description")
					assert.Equal(t.Error,
						task.Subtasks[0].Title, "Some Subtask",
					)
					assert.True(t.Error, !task.Subtasks[0].IsDone)
					assert.Equal(t.Error,
						task.Subtasks[1].Title, "Some Other Subtask",
					)
					assert.True(t.Error, task.Subtasks[1].IsDone)
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPatch,
					"/task?id="+c.taskID,
					strings.NewReader(c.reqBody),
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
					"Only team admins can delete tasks.",
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
				id:   "9dd9c982-8d1c-49ac-a412-3b01ba74b634",
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
									Value: "9dd9c982-8d1c-49ac-a412-3b01ba74b" +
										"634",
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
				r := httptest.NewRequest(
					http.MethodDelete, "/task?id="+c.id, nil,
				)
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
