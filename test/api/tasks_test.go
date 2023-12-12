//go:build itest

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/internal/api"
	tasksAPI "github.com/kxplxn/goteam/internal/api/tasks"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/auth"
	taskTable "github.com/kxplxn/goteam/pkg/db/task"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

func TestTasksAPI(t *testing.T) {
	log := pkgLog.New()
	sut := api.NewHandler(
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodPatch: tasksAPI.NewPATCHHandler(
				token.DecodeAuth,
				token.DecodeState,
				tasksAPI.NewColNoValidator(),
				taskTable.NewMultiUpdater(svcDynamo),
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
				t.Run(http.MethodPatch, func(t *testing.T) {
					req := httptest.NewRequest(http.MethodPatch, "/", nil)
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
			})
		}
	})
	t.Run("PATCH", func(t *testing.T) {
		for _, c := range []struct {
			name       string
			reqBody    string
			authFunc   func(*http.Request)
			statusCode int
			assertFunc func(*testing.T, *http.Response, string)
		}{
			{
				name:       "NoAuth",
				reqBody:    `[]`,
				authFunc:   func(*http.Request) {},
				statusCode: http.StatusUnauthorized,
				assertFunc: assert.OnResErr("Auth token not found."),
			},
			{
				name:       "InvalidAuth",
				reqBody:    `[]`,
				authFunc:   addCookieAuth("asdfjkahsd"),
				statusCode: http.StatusUnauthorized,
				assertFunc: assert.OnResErr("Invalid auth token."),
			},
			{
				name:       "NotAdmin",
				reqBody:    `[]`,
				authFunc:   addCookieAuth(tkTeam1Member),
				statusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr("Only team admins can edit tasks."),
			},
			{
				name:       "NoState",
				reqBody:    `[]`,
				authFunc:   addCookieAuth(tkTeam1Admin),
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("State token not found."),
			},
			{
				name:    "InvalidState",
				reqBody: `[]`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState("askdjfhasdlfk")(r)
				},
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Invalid state token."),
			},
			{
				name:    "InvalidTaskID",
				reqBody: `[{"id": "0"}]`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Invalid task ID."),
			},
			{
				name: "Success",
				reqBody: `[{
                    "id": "c684a6a0-404d-46fa-9fa5-1497f9874567", 
                    "title": "task 5",
                    "order": 2,
                    "subtasks": [],
                    "board": "f0c5d521-ccb5-47cc-ba40-313ddb901165",
                    "column": 2
                }]`,
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkStateTeam1)(r)
				},
				statusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ string) {
					out, err := svcDynamo.GetItem(
						context.Background(),
						&dynamodb.GetItemInput{
							TableName: &taskTableName,
							Key: map[string]types.AttributeValue{
								"TeamID": &types.AttributeValueMemberS{
									Value: "afeadc4a-68b0-4c33-9e83-4648d20ff" +
										"26a",
								},
								"ID": &types.AttributeValueMemberS{
									Value: "c684a6a0-404d-46fa-9fa5-1497f9874" +
										"567",
								},
							},
						},
					)
					assert.Nil(t.Fatal, err)

					var task taskTable.Task
					assert.Nil(t.Fatal, attributevalue.UnmarshalMap(
						out.Item, &task,
					))

					assert.Equal(t.Error,
						task.ID, "c684a6a0-404d-46fa-9fa5-1497f9874567",
					)
					assert.Equal(t.Error, task.Title, "task 5")
					assert.Equal(t.Error, task.Order, 2)
					assert.Equal(t.Error, len(task.Subtasks), 0)
					assert.Equal(t.Error,
						task.BoardID, "f0c5d521-ccb5-47cc-ba40-313ddb901165",
					)
					assert.Equal(t.Error, task.ColumnNumber, 2)
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				req := httptest.NewRequest(
					http.MethodPatch, "/tasks", strings.NewReader(c.reqBody),
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
