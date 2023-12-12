//go:build itest

package api

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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
		t.Run("NoAuth", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/team", nil)
			w := httptest.NewRecorder()
			handler.Handle(w, r, "")

			res := w.Result()

			assert.Equal(t.Error, res.StatusCode, http.StatusUnauthorized)
		})

		t.Run("InvalidAuth", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/team", nil)
			r.AddCookie(&http.Cookie{Name: "auth-token", Value: "asdfasdf"})
			w := httptest.NewRecorder()
			handler.Handle(w, r, "")

			res := w.Result()

			assert.Equal(t.Error, res.StatusCode, http.StatusUnauthorized)
		})

		t.Run("OK", func(t *testing.T) {
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
						Name: "Team 1 Board 2",
					},
					{
						ID:   "1559a33c-54c5-42c8-8e5f-fe096f7760fa",
						Name: "Team 1 Board 3",
					},
				},
			}
			r := httptest.NewRequest(http.MethodGet, "/team", nil)
			addCookieAuth(tkTeam1Member)(r)
			w := httptest.NewRecorder()

			handler.Handle(w, r, "")

			res := w.Result()

			assert.Equal(t.Error, res.StatusCode, http.StatusOK)

			var resp teamAPI.GetResp
			if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t.Error, resp.ID, wantResp.ID)
			for i, m := range wantResp.Members {
				assert.Equal(t.Error, resp.Members[i], m)
			}
			for i, b := range wantResp.Boards {
				assert.Equal(t.Error, resp.Boards[i].ID, b.ID)
				assert.Equal(t.Error, resp.Boards[i].Name, b.Name)
			}
		})
	})

	// if no team was found and since we trust the token, this is our sign from
	// the user endpoint that we should create a new team for the user

	// however, we should not create a new team if the user is not admin - this
	// case should never happen but we should still test it
	t.Run("NotAdmin", func(t *testing.T) {
		authNotAdmin := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc0FkbWluIjp" +
			"mYWxzZX0.Uz6JmqHbxSrzyKAIktxRW4Y_0ldqi_bEcNkYfvIIM8I"
		r := httptest.NewRequest(http.MethodGet, "/team", nil)
		addCookieAuth(authNotAdmin)(r)
		w := httptest.NewRecorder()

		handler.Handle(w, r, "")

		res := w.Result()
		assert.Equal(t.Error, res.StatusCode, http.StatusUnauthorized)
	})

	t.Run("Created", func(t *testing.T) {
		authAdmin := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc0FkbWluIjp0cn" +
			"VlLCJ0ZWFtSUQiOiJkNWRjYTliYy1iYzk4LTQ3YjQtYjhiNy05ZjAxODEzZGE1Nz" +
			"EiLCJ1c2VybmFtZSI6Im5ld3VzZXIifQ.lCjQirzU_3yxOi2bNXRLuyxgzbUnEft" +
			"ITcIFMz2jCb8"
		wantResp := teamAPI.GetResp{
			ID:      "d5dca9bc-bc98-47b4-b8b7-9f01813da571",
			Members: []string{"newuser"},
			Boards:  []teamTable.Board{{Name: "New Board"}},
		}
		r := httptest.NewRequest(http.MethodGet, "/team", nil)
		addCookieAuth(authAdmin)(r)
		w := httptest.NewRecorder()

		handler.Handle(w, r, "")

		res := w.Result()
		var resp teamAPI.GetResp
		if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
			t.Fatal(err)
		}

		assert.Equal(t.Error, resp.ID, wantResp.ID)
		assert.AllEqual(t.Error, resp.Members, wantResp.Members)
		assert.Equal(t.Error, len(resp.Boards), 1)
		assert.Equal(t.Error, resp.Boards[0].Name, wantResp.Boards[0].Name)
	})
}
