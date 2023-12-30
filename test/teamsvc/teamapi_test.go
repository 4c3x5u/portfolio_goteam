//go:build itest

package test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang-jwt/jwt/v4"

	"github.com/kxplxn/goteam/internal/teamsvc/teamapi"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db/teamtbl"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/test"
)

func TestTeamAPI(t *testing.T) {
	handler := teamapi.NewGetHandler(
		cookie.NewAuthDecoder(test.JWTKey),
		teamtbl.NewRetriever(test.DB()),
		teamtbl.NewInserter(test.DB()),
		cookie.NewInviteEncoder(test.JWTKey, 1*time.Hour),
		log.New(),
	)

	t.Run("GET", func(t *testing.T) {
		for _, c := range []struct {
			name       string
			authFunc   func(*http.Request)
			wantStatus int
			assertFunc func(*testing.T, *http.Response)
		}{
			{
				name:       "NoAuth",
				authFunc:   func(r *http.Request) {},
				wantStatus: http.StatusUnauthorized,
				assertFunc: func(*testing.T, *http.Response) {},
			},
			{
				name:       "InvalidAuth",
				authFunc:   test.AddAuthCookie("asdfasdf"),
				wantStatus: http.StatusUnauthorized,
				assertFunc: func(*testing.T, *http.Response) {},
			},
			{
				name: "NotAdmin",
				authFunc: test.AddAuthCookie(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc0FkbWluIjpmYWx" +
						"zZX0.Uz6JmqHbxSrzyKAIktxRW4Y_0ldqi_bEcNkYfvIIM8I",
				),
				wantStatus: http.StatusUnauthorized,
				assertFunc: func(*testing.T, *http.Response) {},
			},
			{
				name: "Created",
				authFunc: test.AddAuthCookie(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc0FkbWluIjp0cn" +
						"VlLCJ0ZWFtSUQiOiJkNWRjYTliYy1iYzk4LTQ3YjQtYjhiNy05Z" +
						"jAxODEzZGE1NzEiLCJ1c2VybmFtZSI6Im5ld3VzZXIifQ.lCjQi" +
						"rzU_3yxOi2bNXRLuyxgzbUnEftITcIFMz2jCb8",
				),
				wantStatus: http.StatusCreated,
				assertFunc: func(t *testing.T, resp *http.Response) {
					wantMembers := []string{"newuser"}
					wantBoardLen := 1
					wantBoardName := "New Board"

					// assert on response body
					var respBody teamapi.GetResp
					err := json.NewDecoder(resp.Body).Decode(&respBody)
					if err != nil {
						t.Fatal(err)
					}
					assert.AllEqual(t.Error, respBody.Members, wantMembers)
					assert.Equal(t.Error, len(respBody.Boards), wantBoardLen)
					assert.Equal(t.Error, respBody.Boards[0].Name, wantBoardName)

					// asssert on db
					out, err := test.DB().GetItem(
						context.Background(),
						&dynamodb.GetItemInput{
							TableName: &tableName,
							Key: map[string]types.AttributeValue{
								"ID": &types.AttributeValueMemberS{
									Value: respBody.ID,
								},
							},
						},
					)
					if err != nil {
						t.Fatal(err)
					}
					var team teamtbl.Team
					err = attributevalue.UnmarshalMap(out.Item, &team)
					if err != nil {
						t.Fatal(err)
					}
					assert.AllEqual(t.Error, team.Members, wantMembers)
					assert.Equal(t.Error, len(team.Boards), wantBoardLen)
					assert.Equal(t.Error, team.Boards[0].Name, wantBoardName)
				},
			},
			{
				name:       "OK",
				authFunc:   test.AddAuthCookie(test.T1AdminToken),
				wantStatus: http.StatusOK,
				assertFunc: func(t *testing.T, resp *http.Response) {
					wantResp := teamapi.GetResp{
						ID:      "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
						Members: []string{"team1Admin", "team1Member"},
						Boards: []teamtbl.Board{
							{
								ID:   "91536664-9749-4dbb-a470-6e52aa353ae4",
								Name: "Team 1 Board 1",
							},
							{
								ID:   "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
								Name: "New Board Name",
							},
							{
								ID:   "1559a33c-54c5-42c8-8e5f-fe096f7760fa",
								Name: "Team 1 Board 3",
							},
						},
					}

					var respBody teamapi.GetResp
					err := json.NewDecoder(resp.Body).Decode(&respBody)
					if err != nil {
						t.Fatal(err)
					}

					assert.Equal(t.Error, respBody.ID, wantResp.ID)
					assert.AllEqual(t.Error,
						respBody.Members, wantResp.Members,
					)
					assert.Equal(t.Error,
						len(respBody.Boards), len(wantResp.Boards),
					)
					for i, b := range wantResp.Boards {
						assert.Equal(t.Error, respBody.Boards[i].ID, b.ID)
						assert.Equal(t.Error, respBody.Boards[i].Name, b.Name)
					}

					ckInv := resp.Cookies()[0]
					assert.Equal(t.Error, ckInv.Name, "invite-token")
					assert.True(t.Error,
						ckInv.Expires.After(time.Now().Add(59*time.Minute)))
					assert.True(t.Error,
						ckInv.Expires.Before(time.Now().Add(61*time.Minute)))
					claims := jwt.MapClaims{}
					if _, err := jwt.ParseWithClaims(
						ckInv.Value, &claims, func(token *jwt.Token) (any, error) {
							return test.JWTKey, nil
						},
					); err != nil {
						t.Error(err)
					}
					assert.Equal(t.Error,
						claims["teamID"].(string),
						"afeadc4a-68b0-4c33-9e83-4648d20ff26a")
					assert.True(t.Error,
						claims["exp"].(float64) >
							float64(time.Now().Add(59*time.Minute).Unix()),
					)
					assert.True(t.Error,
						claims["exp"].(float64) >
							float64(time.Now().Add(59*time.Minute).Unix()),
					)
					assert.True(t.Error,
						claims["exp"].(float64) <
							float64(time.Now().Add(61*time.Minute).Unix()),
					)
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				r := httptest.NewRequest(http.MethodGet, "/team", nil)
				c.authFunc(r)
				w := httptest.NewRecorder()

				handler.Handle(w, r, "")

				resp := w.Result()
				assert.Equal(t.Error, resp.StatusCode, c.wantStatus)
				c.assertFunc(t, resp)
			})
		}
	})
}
