package api

import (
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

		// t.Run("Found", func(t *testing.T) {
		// 	wantResp := teamAPI.GetResp{
		// 		ID:      "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		// 		Members: []string{"team1Admin", "team1Member"},
		// 		Boards: []teamAPI.Board{
		// 			{
		// 				ID:   "91536664-9749-4dbb-a470-6e52aa353ae4",
		// 				Name: "Team 1 Board 1",
		// 			},
		// 			{
		// 				ID:   "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
		// 				Name: "Team 1 Board 2",
		// 			},
		// 			{
		// 				ID:   "1559a33c-54c5-42c8-8e5f-fe096f7760fa",
		// 				Name: "Team 1 Board 3",
		// 			},
		// 		},
		// 	}
		// 	r := httptest.NewRequest(http.MethodGet, "/team", nil)
		// 	addCookieAuth(tkTeam1Member)(r)
		// 	w := httptest.NewRecorder()

		// 	handler.ServeHTTP(w, r)

		// 	res := w.Result()
		// 	var resp teamAPI.GetResp
		// 	if err := json.NewDecoder(res.Body).Decode(&resp); err != nil {
		// 		t.Fatal(err)
		// 	}

		// 	assert.Equal(t.Error, resp.ID, wantResp.ID)
		// 	for i, m := range wantResp.Members {
		// 		assert.Equal(t.Error, resp.Members[i], m)
		// 	}
		// 	for i, b := range wantResp.Boards {
		// 		assert.Equal(t.Error, resp.Boards[i].ID, b.ID)
		// 		assert.Equal(t.Error, resp.Boards[i].Name, b.Name)
		// 	}
		// })
	})
}
