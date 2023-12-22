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
	lgcBoardAPI "github.com/kxplxn/goteam/internal/api/board"
	boardAPI "github.com/kxplxn/goteam/internal/api/team/board"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/auth"
	teamTable "github.com/kxplxn/goteam/pkg/db/team"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

func TestBoardAPI(t *testing.T) {
	log := pkgLog.New()
	idValidator := lgcBoardAPI.NewIDValidator()
	nameValidator := lgcBoardAPI.NewNameValidator()
	sut := api.NewHandler(
		auth.NewJWTValidator(jwtKey), map[string]api.MethodHandler{
			http.MethodDelete: boardAPI.NewDeleteHandler(
				token.DecodeAuth,
				token.DecodeState,
				teamTable.NewBoardDeleter(svcDynamo),
				log,
			),
			http.MethodPatch: lgcBoardAPI.NewPatchHandler(
				token.DecodeAuth,
				token.DecodeState,
				idValidator,
				nameValidator,
				teamTable.NewBoardUpdater(svcDynamo),
				log,
			),
		},
	)

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
				name:       "NoAuth",
				id:         "",
				boardName:  "",
				authFunc:   func(*http.Request) {},
				statusCode: http.StatusUnauthorized,
				assertFunc: assert.OnResErr("Auth token not found."),
			},
			{
				name:       "InvalidAuth",
				id:         "",
				boardName:  "",
				authFunc:   addCookieAuth("asdkfjahsaksdfjhas"),
				statusCode: http.StatusUnauthorized,
				assertFunc: assert.OnResErr("Invalid auth token."),
			},
			{
				name:       "NotAdmin",
				id:         "",
				boardName:  "",
				authFunc:   addCookieAuth(tkTeam1Member),
				statusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can edit boards.",
				),
			},
			{
				name:       "NoState",
				id:         "",
				boardName:  "",
				authFunc:   addCookieAuth(tkTeam1Admin),
				statusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr("State token not found."),
			},
			{
				name:      "InvalidState",
				id:        "",
				boardName: "",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState("asdkljfhaskldfjhasdklf")(r)
				},
				statusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr("Invalid state token."),
			},
			{
				name:      "IDEmpty",
				id:        "",
				boardName: "",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board ID cannot be empty."),
			},
			{
				name:      "IDNotUUID",
				id:        "askdfjhas",
				boardName: "",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board ID must be a UUID."),
			},
			{
				name:      "BoardNameEmpty",
				id:        "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
				boardName: "",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Board name cannot be empty."),
			},
			{
				name:      "BoardNameTooLong",
				id:        "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
				boardName: "A Board Whose Name Is Just Too Long!",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Board name cannot be longer than 35 characters.",
				),
			},
			{
				name:      "ErrNoAccess",
				id:        "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
				boardName: "New Board Name",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam3State)(r)
				},
				statusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"You do not have access to this board.",
				),
			},
			{
				name:      "Success",
				id:        "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
				boardName: "New Board Name",
				authFunc: func(r *http.Request) {
					addCookieAuth(tkTeam1Admin)(r)
					addCookieState(tkTeam1State)(r)
				},
				statusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ string) {
					out, err := svcDynamo.GetItem(
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

					var team *teamTable.Team
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
                        "id": "`+c.id+`",
                        "name": "`+c.boardName+`"
                    }`),
				)
				c.authFunc(r)

				sut.ServeHTTP(w, r)
				res := w.Result()

				assert.Equal(t.Error, res.StatusCode, c.statusCode)

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
					out, err := svcDynamo.GetItem(context.Background(),
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

					var team teamTable.Team
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

				// Run case-specific assertions.
				c.assertFunc(t)
			})
		}
	})
}
