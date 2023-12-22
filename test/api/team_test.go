//go:build itest

package api

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	teamAPI "github.com/kxplxn/goteam/internal/api/team"
	"github.com/kxplxn/goteam/pkg/assert"
	teamTable "github.com/kxplxn/goteam/pkg/db/team"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

func TestTeamAPI(t *testing.T) {
	handler := teamAPI.NewGetHandler(
		token.DecodeAuth,
		teamTable.NewRetriever(svcDynamo),
		teamTable.NewInserter(svcDynamo),
		pkgLog.New(),
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
				authFunc:   addCookieAuth("asdfasdf"),
				wantStatus: http.StatusUnauthorized,
				assertFunc: func(*testing.T, *http.Response) {},
			},
			{
				name:       "OK",
				authFunc:   addCookieAuth(tkTeam1Member),
				wantStatus: http.StatusOK,
				assertFunc: func(t *testing.T, r *http.Response) {
					wantResp := teamAPI.GetResp{
						ID:      "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
						Members: []string{"team1Admin", "team1Member"},
						Boards: []teamTable.Board{
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

					var resp teamAPI.GetResp
					err := json.NewDecoder(r.Body).Decode(&resp)
					if err != nil {
						t.Fatal(err)
					}

					assert.Equal(t.Error, resp.ID, wantResp.ID)
					assert.AllEqual(t.Error, resp.Members, wantResp.Members)
					assert.Equal(t.Error, len(resp.Boards), len(wantResp.Boards))
					for i, b := range wantResp.Boards {
						assert.Equal(t.Error, resp.Boards[i].ID, b.ID)
						assert.Equal(t.Error, resp.Boards[i].Name, b.Name)
					}
				},
			},
			// refer to comments in implementation for the below tests
			{
				name: "NotAdmin",
				authFunc: addCookieAuth(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc0FkbWluIjpmYWx" +
						"zZX0.Uz6JmqHbxSrzyKAIktxRW4Y_0ldqi_bEcNkYfvIIM8I",
				),
				wantStatus: http.StatusUnauthorized,
				assertFunc: func(t *testing.T, r *http.Response) {},
			},
			{
				name: "Created",
				authFunc: addCookieAuth(
					"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc0FkbWluIjp0cn" +
						"VlLCJ0ZWFtSUQiOiJkNWRjYTliYy1iYzk4LTQ3YjQtYjhiNy05Z" +
						"jAxODEzZGE1NzEiLCJ1c2VybmFtZSI6Im5ld3VzZXIifQ.lCjQi" +
						"rzU_3yxOi2bNXRLuyxgzbUnEftITcIFMz2jCb8",
				),
				wantStatus: http.StatusCreated,
				assertFunc: func(t *testing.T, r *http.Response) {
					wantMembers := []string{"newuser"}
					wantBoardLen := 1
					wantBoardName := "New Board"

					// assert on response
					var resp teamAPI.GetResp
					err := json.NewDecoder(r.Body).Decode(&resp)
					if err != nil {
						t.Fatal(err)
					}
					assert.AllEqual(t.Error, resp.Members, wantMembers)
					assert.Equal(t.Error, len(resp.Boards), wantBoardLen)
					assert.Equal(t.Error, resp.Boards[0].Name, wantBoardName)

					// asssert on db
					out, err := svcDynamo.GetItem(
						context.Background(),
						&dynamodb.GetItemInput{
							TableName: &teamTableName,
							Key: map[string]types.AttributeValue{
								"ID": &types.AttributeValueMemberS{
									Value: resp.ID,
								},
							},
						},
					)
					if err != nil {
						t.Fatal(err)
					}
					var team teamTable.Team
					err = attributevalue.UnmarshalMap(out.Item, &team)
					if err != nil {
						t.Fatal(err)
					}
					assert.AllEqual(t.Error, team.Members, wantMembers)
					assert.Equal(t.Error, len(team.Boards), wantBoardLen)
					assert.Equal(t.Error, team.Boards[0].Name, wantBoardName)
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				r := httptest.NewRequest(http.MethodGet, "/team", nil)
				c.authFunc(r)
				w := httptest.NewRecorder()
				handler.Handle(w, r, "")

				res := w.Result()

				assert.Equal(t.Error, res.StatusCode, c.wantStatus)

				c.assertFunc(t, res)
			})
		}
	})
}
