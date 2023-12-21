//go:build itest

package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/internal/api"
	boardAPI "github.com/kxplxn/goteam/internal/api/team/board"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/auth"
	teamTable "github.com/kxplxn/goteam/pkg/db/team"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

func TestBoardAPI(t *testing.T) {
	log := pkgLog.New()
	sut := api.NewHandler(
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodDelete: boardAPI.NewDeleteHandler(
				token.DecodeAuth,
				token.DecodeState,
				teamTable.NewBoardDeleter(svcDynamo, svcDynamo),
				log,
			),
		},
	)

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
