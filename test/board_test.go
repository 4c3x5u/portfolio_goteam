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

	boardAPI "github.com/kxplxn/goteam/internal/team/board"
	"github.com/kxplxn/goteam/pkg/api"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db/teamtable"
	"github.com/kxplxn/goteam/pkg/log"
)

func TestBoardAPI(t *testing.T) {
	nameValidator := boardAPI.NewNameValidator()
	log := log.New()
	sut := api.NewHandler(map[string]api.MethodHandler{
		http.MethodPost: boardAPI.NewPostHandler(
			authDecoder,
			stateDecoder,
			nameValidator,
			teamtable.NewBoardInserter(db),
			stateEncoder,
			log,
		),
		http.MethodDelete: boardAPI.NewDeleteHandler(
			authDecoder,
			stateDecoder,
			teamtable.NewBoardDeleter(db),
			stateEncoder,
			log,
		),
		http.MethodPatch: boardAPI.NewPatchHandler(
			authDecoder,
			stateDecoder,
			boardAPI.NewIDValidator(),
			nameValidator,
			teamtable.NewBoardUpdater(db),
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
				assertFunc: assert.OnResErr("Auth token not found."),
			},
			{
				name:       "InvalidAuth",
				boardName:  "",
				authFunc:   addCookieAuth("asdkfjahsaksdfjhas"),
				wantStatus: http.StatusUnauthorized,
				assertFunc: assert.OnResErr("Invalid auth token."),
			},
			{
				name:       "NotAdmin",
				boardName:  "",
				authFunc:   addCookieAuth(tkTeam4Member),
				wantStatus: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can edit boards.",
				),
			},
			{
				name:       "NoState",
				boardName:  "",
				authFunc:   addCookieAuth(tkTeam4Admin),
				wantStatus: http.StatusForbidden,
				assertFunc: assert.OnResErr("State token not found."),
			},
			{
				name:      "InvalidState",
				boardName: "",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam4Admin)(r)
					addCookieState("asdkljfhaskldfjhasdklf")(r)
				},
				wantStatus: http.StatusForbidden,
				assertFunc: assert.OnResErr("Invalid state token."),
			},
			{
				name: "EmptyBoardName",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam4Admin)(r)
					addCookieState(tkTeam4State)(r)
				},
				boardName:  "",
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board name cannot be empty."),
			},
			{
				name: "TooLongBoardName",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam4Admin)(r)
					addCookieState(tkTeam4State)(r)
				},
				boardName:  "A Board Whose Name Is Just Too Long!",
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Board name cannot be longer than 35 characters.",
				),
			},
			{
				name: "LimitReached",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				boardName:  "bob123's new board",
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"You have already created the maximum amount of boards " +
						"allowed per team. Please delete one of your boards " +
						"to create a new one.",
				),
			},
			{
				name: "OK",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam4Admin)(r)
					addCookieState(tkTeam4State)(r)
				},
				boardName:  "Team 4 Board 1",
				wantStatus: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ []any) {
					out, err := db.GetItem(
						context.Background(), &dynamodb.GetItemInput{
							TableName: &teamTableName,
							Key: map[string]types.AttributeValue{
								"ID": &types.AttributeValueMemberS{
									Value: "3c3ec4ea-a850-4fc5-aab0-24e9e7223" +
										"bbc",
								},
							},
						},
					)
					assert.Nil(t.Fatal, err)

					var team *teamtable.Team
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
				res := w.Result()

				assert.Equal(t.Error, res.StatusCode, c.wantStatus)

				// run case-specific assertions
				c.assertFunc(t, res, []any{})
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
				assertFunc: assert.OnResErr("Auth token not found."),
			},
			{
				name:       "InvalidAuth",
				boardID:    "",
				boardName:  "",
				authFunc:   addCookieAuth("asdkfjahsaksdfjhas"),
				wantStatus: http.StatusUnauthorized,
				assertFunc: assert.OnResErr("Invalid auth token."),
			},
			{
				name:       "NotAdmin",
				boardID:    "",
				boardName:  "",
				authFunc:   addCookieAuth(tkTeam1Member),
				wantStatus: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can edit boards.",
				),
			},
			{
				name:       "NoState",
				boardID:    "",
				boardName:  "",
				authFunc:   addCookieAuth(tkTeam1Admin),
				wantStatus: http.StatusForbidden,
				assertFunc: assert.OnResErr("State token not found."),
			},
			{
				name:      "InvalidState",
				boardID:   "",
				boardName: "",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState("asdkljfhaskldfjhasdklf")(r)
				},
				wantStatus: http.StatusForbidden,
				assertFunc: assert.OnResErr("Invalid state token."),
			},
			{
				name:      "IDEmpty",
				boardID:   "",
				boardName: "",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board ID cannot be empty."),
			},
			{
				name:      "IDNotUUID",
				boardID:   "askdfjhas",
				boardName: "",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board ID must be a UUID."),
			},
			{
				name:      "BoardNameEmpty",
				boardID:   "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
				boardName: "",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board name cannot be empty."),
			},
			{
				name:      "BoardNameTooLong",
				boardID:   "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
				boardName: "A Board Whose Name Is Just Too Long!",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				wantStatus: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Board name cannot be longer than 35 characters.",
				),
			},
			{
				name:      "ErrNoAccess",
				boardID:   "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
				boardName: "New Board Name",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam3State)(r)
				},
				wantStatus: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"You do not have access to this board.",
				),
			},
			{
				name:      "OK",
				boardID:   "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
				boardName: "New Board Name",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				wantStatus: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ []any) {
					out, err := db.GetItem(
						context.Background(), &dynamodb.GetItemInput{
							TableName: &teamTableName,
							Key: map[string]types.AttributeValue{
								"ID": &types.AttributeValueMemberS{
									Value: "afeadc4a-68b0-4c33-9e83-4648d20ff" +
										"26a",
								},
							},
						},
					)
					assert.Nil(t.Fatal, err)

					var team *teamtable.Team
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
				res := w.Result()

				assert.Equal(t.Error, res.StatusCode, c.wantStatus)

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
			assertFunc     func(*testing.T)
		}{
			{
				name:           "NotAdmin",
				id:             "f0c5d521-ccb5-47cc-ba40-313ddb901165",
				authFunc:       addCookieAuth(tkTeam1Member),
				wantStatusCode: http.StatusForbidden,
				assertFunc:     func(*testing.T) {},
			},
			{
				name:           "NoState",
				id:             "f0c5d521-ccb5-47cc-ba40-313ddb901165",
				authFunc:       addCookieAuth(tkTeam3Admin),
				wantStatusCode: http.StatusForbidden,
				assertFunc:     func(*testing.T) {},
			},
			{
				name: "NoAccess",
				id:   "f0c5d521-ccb5-47cc-ba40-313ddb901165",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam3Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				wantStatusCode: http.StatusForbidden,
				assertFunc:     func(*testing.T) {},
			},
			{
				name: "EmptyID",
				id:   "",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam3Admin)(r)
					addCookieState(tkTeam3State)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     func(*testing.T) {},
			},
			{
				name: "InvalidID",
				id:   "qwerty",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam3Admin)(r)
					addCookieState(tkTeam3State)(r)
				},
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     func(*testing.T) {},
			},
			{
				name: "Success",
				id:   "f0c5d521-ccb5-47cc-ba40-313ddb901165",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam3Admin)(r)
					addCookieState(tkTeam3State)(r)
				},
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T) {
					out, err := db.GetItem(context.Background(),
						&dynamodb.GetItemInput{
							TableName: &teamTableName,
							Key: map[string]types.AttributeValue{
								"ID": &types.AttributeValueMemberS{
									Value: "74c80ae5-64f3-4298-a8ff-48f8f920c" +
										"7d4",
								},
							},
						},
					)
					assert.Nil(t.Fatal, err)

					var team teamtable.Team
					err = attributevalue.UnmarshalMap(out.Item, &team)
					assert.Nil(t.Fatal, err)

					assert.Equal(t.Error, len(team.Boards), 0)
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				r := httptest.NewRequest(
					http.MethodDelete, "/team/board?id="+c.id, nil,
				)
				c.authFunc(r)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, r)
				res := w.Result()

				assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

				// run case-specific assertions
				c.assertFunc(t)
			})
		}
	})
}
