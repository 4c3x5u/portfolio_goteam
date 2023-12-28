//go:build itest

package tasksvc

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
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db/tasktbl"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/test"
)

func TestTaskAPI(t *testing.T) {
	authDecoder := cookie.NewAuthDecoder(test.JWTKey)
	titleValidator := taskapi.NewTitleValidator()
	log := log.New()
	sut := api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: taskapi.NewPostHandler(
			authDecoder,
			titleValidator,
			titleValidator,
			taskapi.NewColNoValidator(),
			tasktbl.NewInserter(test.DB()),
			log,
		),
		http.MethodPatch: taskapi.NewPatchHandler(
			authDecoder,
			titleValidator,
			titleValidator,
			tasktbl.NewUpdater(test.DB()),
			log,
		),
		http.MethodDelete: taskapi.NewDeleteHandler(
			authDecoder,
			tasktbl.NewDeleter(test.DB()),
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
				assertFunc:     assert.OnRespErr("Auth token not found."),
			},
			{
				name:           "InvalidAuth",
				reqBody:        `{}`,
				authFunc:       test.AddAuthCookie("asdfkjldfs"),
				wantStatusCode: http.StatusUnauthorized,
				assertFunc:     assert.OnRespErr("Invalid auth token."),
			},
			{
				name:           "NotAdmin",
				reqBody:        `{}`,
				authFunc:       test.AddAuthCookie(test.T1MemberToken),
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnRespErr(
					"Only team admins can create tasks.",
				),
			},
			{
				name: "ColNoOutOfBounds",
				reqBody: `{
                    "column": 5,
                    "board":  "91536664-9749-4dbb-a470-6e52aa353ae4"
                }`,
				authFunc:       test.AddAuthCookie(test.T1AdminToken),
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnRespErr(
					"Column number out of bounds.",
				),
			},
			{
				name: "TaskTitleEmpty",
				reqBody: `{
                    "board":  "91536664-9749-4dbb-a470-6e52aa353ae4",
                    "column": 1,
					"title":  ""
				}`,
				authFunc:       test.AddAuthCookie(test.T1AdminToken),
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnRespErr("Task title cannot be empty."),
			},
			{
				name: "TaskTitleTooLong",
				reqBody: `{
                    "board":  "91536664-9749-4dbb-a470-6e52aa353ae4",
                    "column": 1,
					"title":  "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd"
				}`,
				authFunc:       test.AddAuthCookie(test.T1AdminToken),
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnRespErr(
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
				authFunc:       test.AddAuthCookie(test.T1AdminToken),
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnRespErr(
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
				authFunc:       test.AddAuthCookie(test.T1AdminToken),
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnRespErr(
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
					"subtasks":    ["Some Subtask", "Some Other Subtask"],
                    "order":       2
				}`,
				authFunc:       test.AddAuthCookie(test.T1AdminToken),
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
						out, err := test.DB().Query(
							context.Background(),
							&dynamodb.QueryInput{
								TableName: &tableName,
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
								task.ColNo == 1 &&
								task.Title == "Some Task" &&
								task.Subtasks[0].Title == "Some Subtask" &&
								task.Subtasks[0].IsDone == false &&
								task.Subtasks[1].Title == "Some Other Subtas"+
									"k" &&
								task.Subtasks[1].IsDone == false &&
								task.Order == 2 {

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
				w := httptest.NewRecorder()
				r := httptest.NewRequest(
					http.MethodPost, "/tasks/task", strings.NewReader(c.reqBody),
				)
				c.authFunc(r)

				sut.ServeHTTP(w, r)

				resp := w.Result()
				assert.Equal(t.Error, resp.StatusCode, c.wantStatusCode)
				c.assertFunc(t, resp, []any{})
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
				authFunc:       test.AddAuthCookie(test.T1MemberToken),
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnRespErr(
					"Only team admins can edit tasks.",
				),
			},
			{
				name:   "TaskTitleEmpty",
				taskID: "e0021a56-6a1e-4007-b773-395d3991fb7e",
				reqBody: `{
					"title":       "",
					"description": "",
					"subtasks":    []
				}`,
				authFunc:       test.AddAuthCookie(test.T1AdminToken),
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnRespErr("Task title cannot be empty."),
			},
			{
				name:   "TaskTitleTooLong",
				taskID: "e0021a56-6a1e-4007-b773-395d3991fb7e",
				reqBody: `{
					"title": "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd",
					"description": "",
					"subtasks":    []
				}`,
				authFunc:       test.AddAuthCookie(test.T1AdminToken),
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnRespErr(
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
				authFunc:       test.AddAuthCookie(test.T1AdminToken),
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnRespErr(
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
				authFunc:       test.AddAuthCookie(test.T1AdminToken),
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnRespErr(
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
				authFunc:       test.AddAuthCookie(test.T1AdminToken),
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ []any) {
					out, err := test.DB().GetItem(
						context.Background(),
						&dynamodb.GetItemInput{
							TableName: &tableName,
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
				w := httptest.NewRecorder()
				r := httptest.NewRequest(
					http.MethodPatch,
					"/tasks/task?id="+c.taskID,
					strings.NewReader(c.reqBody),
				)
				c.authFunc(r)

				sut.ServeHTTP(w, r)

				resp := w.Result()
				assert.Equal(t.Error, resp.StatusCode, c.wantStatusCode)
				c.assertFunc(t, resp, []any{})
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
				assertFunc:     assert.OnRespErr("Auth token not found."),
			},
			{
				name:           "InvalidAuth",
				id:             "",
				authFunc:       test.AddAuthCookie("asdfasdf"),
				wantStatusCode: http.StatusUnauthorized,
				assertFunc:     assert.OnRespErr("Invalid auth token."),
			},
			{
				name:           "NotAdmin",
				id:             "",
				authFunc:       test.AddAuthCookie(test.T1MemberToken),
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnRespErr(
					"Only team admins can delete tasks.",
				),
			},
			{
				name:           "OK",
				id:             "9dd9c982-8d1c-49ac-a412-3b01ba74b634",
				authFunc:       test.AddAuthCookie(test.T1AdminToken),
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ []any) {
					out, err := test.DB().GetItem(
						context.Background(),
						&dynamodb.GetItemInput{
							TableName: &tableName,
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
				w := httptest.NewRecorder()
				r := httptest.NewRequest(
					http.MethodDelete, "/tasks/task?id="+c.id, nil,
				)
				c.authFunc(r)

				sut.ServeHTTP(w, r)

				resp := w.Result()
				assert.Equal(t.Error, resp.StatusCode, c.wantStatusCode)
				c.assertFunc(t, resp, []any{})
			})
		}
	})
}
