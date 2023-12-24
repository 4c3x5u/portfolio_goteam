//go:build itest

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/internal/tasksvc/taskapi"
	"github.com/kxplxn/goteam/pkg/api"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db/tasktbl"
	"github.com/kxplxn/goteam/pkg/log"
)

func TestTaskAPI(t *testing.T) {
	titleValidator := taskapi.NewTitleValidator()
	log := log.New()

	sut := api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: taskapi.NewPostHandler(
			authDecoder,
			stateDecoder,
			titleValidator,
			titleValidator,
			taskapi.NewColNoValidator(),
			tasktbl.NewInserter(db),
			stateEncoder,
			log,
		),
		http.MethodPatch: taskapi.NewPatchHandler(
			authDecoder,
			stateDecoder,
			titleValidator,
			titleValidator,
			tasktbl.NewUpdater(db),
			log,
		),
		http.MethodDelete: taskapi.NewDeleteHandler(
			authDecoder,
			stateDecoder,
			tasktbl.NewDeleter(db),
			stateEncoder,
			log,
		),
	})

	t.Run("POST", func(t *testing.T) {
		for _, c := range []struct {
			name           string
			reqBody        string
			authFunc       func(*http.Request)
			wantStatusCode int
			assertFunc     func(*testing.T, *http.Response, []any)
		}{
			{
				name:           "NoAuth",
				reqBody:        `{}`,
				authFunc:       func(*http.Request) {},
				wantStatusCode: http.StatusUnauthorized,
				assertFunc:     assert.OnResErr("Auth token not found."),
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
					addCookieState(tkTeam1State)(r)
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
					addCookieState(tkTeam1State)(r)
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
					addCookieState(tkTeam1State)(r)
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
					addCookieState(tkTeam1State)(r)
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
					addCookieState(tkTeam1State)(r)
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
					addCookieState(tkTeam1State)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Subtask title cannot be longer than 50 characters.",
				),
			},
			{
				name: "OK",
				reqBody: `{
                    "board":       "91536664-9749-4dbb-a470-6e52aa353ae4",
					"description": "Do something. Then, do something else.",
                    "column":      1,
					"title":       "Some Task",
					"subtasks":    ["Some Subtask", "Some Other Subtask"]
				}`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ []any) {
					keyEx := expression.Key("BoardID").Equal(expression.Value(
						"91536664-9749-4dbb-a470-6e52aa353ae4",
					))
					expr, err := expression.NewBuilder().
						WithKeyCondition(keyEx).Build()
					assert.Nil(t.Fatal, err)
					time.Sleep(1 * time.Second)

					// try a few times as it takes a while for the task to
					// appear in the database for some reason
					// FIXME: find the root cause and eliminate it if possible
					//        OR just make the switch to a DynamoDB instance on
					//        Docker for local testing...
					var taskFound bool
					for i := 0; i < 3; i++ {
						out, err := db.Query(
							context.Background(),
							&dynamodb.QueryInput{
								TableName: &taskTableName,
								IndexName: aws.String(
									"BoardID_index",
								),
								ExpressionAttributeNames:  expr.Names(),
								ExpressionAttributeValues: expr.Values(),
								KeyConditionExpression:    expr.KeyCondition(),
							},
						)
						assert.Nil(t.Fatal, err)

						for _, av := range out.Items {
							var task tasktbl.Task
							err = attributevalue.UnmarshalMap(av, &task)
							assert.Nil(t.Fatal, err)

							wantDescr := "Do something. Then, do something " +
								"else."
							if task.Description == wantDescr &&
								task.ColumnNumber == 1 &&
								task.Title == "Some Task" &&
								task.Subtasks[0].Title == "Some Subtask" &&
								task.Subtasks[0].IsDone == false &&
								task.Subtasks[1].Title == "Some Other Subtas"+
									"k" &&
								task.Subtasks[1].IsDone == false {

								taskFound = true
								break
							}
						}

						time.Sleep(2 * time.Second)
					}

					assert.True(t.Fatal, taskFound)
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPost, "/tasks/task", strings.NewReader(c.reqBody),
				)
				c.authFunc(req)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)

				res := w.Result()
				assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)
				c.assertFunc(t, res, []any{})
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
			assertFunc     func(*testing.T, *http.Response, []any)
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
					addCookieState(tkTeam1State)(r)
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
					addCookieState(tkTeam1State)(r)
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
					addCookieState(tkTeam1State)(r)
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
					addCookieState(tkTeam1State)(r)
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
					addCookieState(tkTeam1State)(r)
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
					addCookieState(tkTeam1State)(r)
				},
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ []any) {
					out, err := db.GetItem(
						context.Background(),
						&dynamodb.GetItemInput{
							TableName: &taskTableName,
							Key: map[string]types.AttributeValue{
								"TeamID": &types.AttributeValueMemberS{
									Value: "afeadc4a-68b0-4c33-9e83-4648d20ff" +
										"26a",
								},
								"ID": &types.AttributeValueMemberS{
									Value: "e0021a56-6a1e-4007-b773-395d3991f" +
										"b7e",
								},
							},
						},
					)
					assert.Nil(t.Fatal, err)

					var task tasktbl.Task
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
					"/tasks/task?id="+c.taskID,
					strings.NewReader(c.reqBody),
				)
				c.authFunc(req)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)

				res := w.Result()
				assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)
				c.assertFunc(t, res, []any{})
			})
		}
	})

	t.Run("DELETE", func(t *testing.T) {
		for _, c := range []struct {
			name           string
			id             string
			authFunc       func(*http.Request)
			wantStatusCode int
			assertFunc     func(*testing.T, *http.Response, []any)
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
					addCookieState(tkTeam1State)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Invalid task ID."),
			},
			{
				name: "OK",
				id:   "9dd9c982-8d1c-49ac-a412-3b01ba74b634",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ []any) {
					out, err := db.GetItem(
						context.Background(),
						&dynamodb.GetItemInput{
							TableName: &taskTableName,
							Key: map[string]types.AttributeValue{
								"TeamID": &types.AttributeValueMemberS{
									Value: "91536664-9749-4dbb-a470-6e52aa353" +
										"ae4",
								},
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
					http.MethodDelete, "/tasks/task?id="+c.id, nil,
				)
				c.authFunc(r)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, r)

				res := w.Result()
				assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)
				c.assertFunc(t, res, []any{})
			})
		}
	})
}
