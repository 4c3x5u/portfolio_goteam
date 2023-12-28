//go:build itest

package test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/internal/teamsvc/boardapi"
	"github.com/kxplxn/goteam/pkg/api"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db/teamtbl"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/test"
)

func TestBoardAPI(t *testing.T) {
	authDecoder := cookie.NewAuthDecoder(test.JWTKey)
	nameValidator := boardapi.NewNameValidator()
	log := log.New()
	sut := api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: boardapi.NewPostHandler(
			authDecoder,
			nameValidator,
			teamtbl.NewBoardInserter(test.DB()),
			log,
		),
		http.MethodDelete: boardapi.NewDeleteHandler(
			authDecoder,
			teamtbl.NewBoardDeleter(test.DB()),
			log,
		),
		http.MethodPatch: boardapi.NewPatchHandler(
			authDecoder,
			boardapi.NewIDValidator(),
			nameValidator,
			teamtbl.NewBoardUpdater(test.DB()),
			log,
		),
	})

	t.Run("POST", func(t *testing.T) {
		for _, c := range []struct {
			name       string
			authFunc   func(*http.Request)
			boardName  string
			wantStatus int
			assertFunc func(*testing.T, *http.Response, []any)
		}{
			{
				name:       "NoAuth",
				boardName:  "",
				authFunc:   func(*http.Request) {},
				wantStatus: http.StatusUnauthorized,
				assertFunc: assert.OnRespErr("Auth token not found."),
			},
			{
				name:       "InvalidAuth",
				boardName:  "",
				authFunc:   test.AddAuthCookie("asdkfjahsaksdfjhas"),
				wantStatus: http.StatusUnauthorized,
				assertFunc: assert.OnRespErr("Invalid auth token."),
			},
			{
				name:       "NotAdmin",
				boardName:  "",
				authFunc:   test.AddAuthCookie(test.T4MemberToken),
				wantStatus: http.StatusForbidden,
				assertFunc: assert.OnRespErr(
					"Only team admins can edit boards.",
				),
			},
			{
				name: "EmptyBoardName",
				authFunc: func(r *http.Request) {
					test.AddAuthCookie(test.T4AdminToken)(r)
					test.AddStateCookie(test.T4StateToken)(r)
				},
				boardName:  "",
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnRespErr("Board name cannot be empty."),
			},
			{
				name: "TooLongBoardName",
				authFunc: func(r *http.Request) {
					test.AddAuthCookie(test.T4AdminToken)(r)
					test.AddStateCookie(test.T4StateToken)(r)
				},
				boardName:  "A Board Whose Name Is Just Too Long!",
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnRespErr(
					"Board name cannot be longer than 35 characters.",
				),
			},
			{
				name: "LimitReached",
				authFunc: func(r *http.Request) {
					test.AddAuthCookie(test.T1AdminToken)(r)
					test.AddStateCookie(test.T1StateToken)(r)
				},
				boardName:  "bob123's new board",
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnRespErr(
					"You have already created the maximum amount of boards " +
						"allowed per team. Please delete one of your boards " +
						"to create a new one.",
				),
			},
			{
				name: "OK",
				authFunc: func(r *http.Request) {
					test.AddAuthCookie(test.T4AdminToken)(r)
					test.AddStateCookie(test.T4StateToken)(r)
				},
				boardName:  "Team 4 Board 1",
				wantStatus: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ []any) {
					out, err := test.DB().GetItem(
						context.Background(), &dynamodb.GetItemInput{
							TableName: &tableName,
							Key: map[string]types.AttributeValue{
								"ID": &types.AttributeValueMemberS{
									Value: "3c3ec4ea-a850-4fc5-aab0-24e9e7223" +
										"bbc",
								},
							},
						},
					)
					assert.Nil(t.Fatal, err)

					var team *teamtbl.Team
					err = attributevalue.UnmarshalMap(out.Item, &team)
					assert.Nil(t.Fatal, err)

					var found bool
					for _, b := range team.Boards {
						if b.Name == "Team 4 Board 1" {
							found = true
							break
						}
					}
					assert.True(t.Error, found)
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest(
					http.MethodPost, "/team/board", strings.NewReader(`{
                        "name": "`+c.boardName+`"
                    }`),
				)
				c.authFunc(r)

				sut.ServeHTTP(w, r)

				resp := w.Result()
				assert.Equal(t.Error, resp.StatusCode, c.wantStatus)
				c.assertFunc(t, resp, []any{})
			})
		}
	})

	t.Run("PATCH", func(t *testing.T) {
		for _, c := range []struct {
			name       string
			boardID    string
			boardName  string
			authFunc   func(*http.Request)
			wantStatus int
			assertFunc func(*testing.T, *http.Response, []any)
		}{
			{
				name:       "NoAuth",
				boardID:    "",
				boardName:  "",
				authFunc:   func(*http.Request) {},
				wantStatus: http.StatusUnauthorized,
				assertFunc: assert.OnRespErr("Auth token not found."),
			},
			{
				name:       "InvalidAuth",
				boardID:    "",
				boardName:  "",
				authFunc:   test.AddAuthCookie("asdkfjahsaksdfjhas"),
				wantStatus: http.StatusUnauthorized,
				assertFunc: assert.OnRespErr("Invalid auth token."),
			},
			{
				name:       "NotAdmin",
				boardID:    "",
				boardName:  "",
				authFunc:   test.AddAuthCookie(test.T1MemberToken),
				wantStatus: http.StatusForbidden,
				assertFunc: assert.OnRespErr(
					"Only team admins can edit boards.",
				),
			},
			{
				name:       "IDEmpty",
				boardID:    "",
				boardName:  "",
				authFunc:   test.AddAuthCookie(test.T1AdminToken),
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnRespErr("Board ID cannot be empty."),
			},
			{
				name:       "IDNotUUID",
				boardID:    "askdfjhas",
				boardName:  "",
				authFunc:   test.AddAuthCookie(test.T1AdminToken),
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnRespErr("Board ID must be a UUID."),
			},
			{
				name:       "BoardNameEmpty",
				boardID:    "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
				boardName:  "",
				authFunc:   test.AddAuthCookie(test.T1AdminToken),
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnRespErr("Board name cannot be empty."),
			},
			{
				name:       "BoardNameTooLong",
				boardID:    "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
				boardName:  "A Board Whose Name Is Just Too Long!",
				authFunc:   test.AddAuthCookie(test.T1AdminToken),
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnRespErr(
					"Board name cannot be longer than 35 characters.",
				),
			},
			{
				name:       "OK",
				boardID:    "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
				boardName:  "New Board Name",
				authFunc:   test.AddAuthCookie(test.T1AdminToken),
				wantStatus: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ []any) {
					out, err := test.DB().GetItem(
						context.Background(), &dynamodb.GetItemInput{
							TableName: &tableName,
							Key: map[string]types.AttributeValue{
								"ID": &types.AttributeValueMemberS{
									Value: "afeadc4a-68b0-4c33-9e83-4648d20ff" +
										"26a",
								},
							},
						},
					)
					assert.Nil(t.Fatal, err)

					var team *teamtbl.Team
					err = attributevalue.UnmarshalMap(out.Item, &team)
					assert.Nil(t.Fatal, err)

					var found bool
					for _, b := range team.Boards {
						if b.ID == "fdb82637-f6a5-4d55-9dc3-9f60061e632f" {
							assert.Equal(t.Error, b.Name, "New Board Name")
							found = true
							break
						}
					}
					if !found {
						t.Error("board not found for team")
					}
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest(
					http.MethodPatch, "/team/board", strings.NewReader(`{
                        "id": "`+c.boardID+`",
                        "name": "`+c.boardName+`"
                    }`),
				)
				c.authFunc(r)

				sut.ServeHTTP(w, r)

				resp := w.Result()
				assert.Equal(t.Error, resp.StatusCode, c.wantStatus)
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
			assertFunc     func(*testing.T)
		}{
			{
				name:           "NotAdmin",
				id:             "f0c5d521-ccb5-47cc-ba40-313ddb901165",
				authFunc:       test.AddAuthCookie(test.T1MemberToken),
				wantStatusCode: http.StatusForbidden,
				assertFunc:     func(*testing.T) {},
			},
			{
				name: "EmptyID",
				id:   "",
				authFunc: func(r *http.Request) {
					test.AddAuthCookie(test.T3AdminToken)(r)
					test.AddStateCookie(test.T3StateToken)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     func(*testing.T) {},
			},
			{
				name: "InvalidID",
				id:   "qwerty",
				authFunc: func(r *http.Request) {
					test.AddAuthCookie(test.T3AdminToken)(r)
					test.AddStateCookie(test.T3StateToken)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     func(*testing.T) {},
			},
			{
				name: "Success",
				id:   "f0c5d521-ccb5-47cc-ba40-313ddb901165",
				authFunc: func(r *http.Request) {
					test.AddAuthCookie(test.T3AdminToken)(r)
					test.AddStateCookie(test.T3StateToken)(r)
				},
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T) {
					out, err := test.DB().GetItem(context.Background(),
						&dynamodb.GetItemInput{
							TableName: &tableName,
							Key: map[string]types.AttributeValue{
								"ID": &types.AttributeValueMemberS{
									Value: "74c80ae5-64f3-4298-a8ff-48f8f920c" +
										"7d4",
								},
							},
						},
					)
					assert.Nil(t.Fatal, err)

					var team teamtbl.Team
					err = attributevalue.UnmarshalMap(out.Item, &team)
					assert.Nil(t.Fatal, err)

					assert.Equal(t.Error, len(team.Boards), 0)
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				w := httptest.NewRecorder()
				r := httptest.NewRequest(
					http.MethodDelete, "/team/board?id="+c.id, nil,
				)
				c.authFunc(r)

				sut.ServeHTTP(w, r)

				resp := w.Result()
				assert.Equal(t.Error, resp.StatusCode, c.wantStatusCode)
				c.assertFunc(t)
			})
		}
	})
}
