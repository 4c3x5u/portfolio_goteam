package team

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
	teamTable "github.com/kxplxn/goteam/pkg/db/team"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

func TestGetHandler(t *testing.T) {
	decodeAuth := &token.FakeDecode[token.Auth]{}
	retriever := &db.FakeRetriever[teamTable.Team]{}
	log := &pkgLog.FakeErrorer{}
	sut := NewGetHandler(decodeAuth.Func, retriever, log)

	wantTeam := teamTable.Team{
		ID:      "teamid",
		Members: []string{"memberone", "membertwo"},
		Boards: []teamTable.Board{
			{ID: "board1", Name: "boardone"},
			{ID: "board2", Name: "boardtwo"},
		},
	}

	for _, c := range []struct {
		name          string
		auth          string
		errDecodeAuth error
		authDecoded   token.Auth
		errRetrieve   error
		team          teamTable.Team
		wantStatus    int
		assertFunc    func(*testing.T, *http.Response, string)
	}{
		{
			name:          "NoAuth",
			auth:          "",
			errDecodeAuth: nil,
			authDecoded:   token.Auth{},
			errRetrieve:   nil,
			team:          teamTable.Team{},
			wantStatus:    http.StatusUnauthorized,
			assertFunc:    func(*testing.T, *http.Response, string) {},
		},
		{
			name:          "InvalidAuth",
			auth:          "nonempty",
			errDecodeAuth: errors.New("decode auth failed"),
			authDecoded:   token.Auth{},
			errRetrieve:   nil,
			team:          teamTable.Team{},
			wantStatus:    http.StatusUnauthorized,
			assertFunc:    func(*testing.T, *http.Response, string) {},
		},
		{
			name:          "ErrRetrieve",
			auth:          "nonempty",
			errDecodeAuth: nil,
			authDecoded:   token.Auth{},
			errRetrieve:   errors.New("retrieve failed"),
			team:          teamTable.Team{},
			wantStatus:    http.StatusInternalServerError,
			assertFunc:    assert.OnLoggedErr("retrieve failed"),
		},
		{
			name:          "OK",
			auth:          "nonempty",
			errDecodeAuth: nil,
			authDecoded:   token.Auth{},
			errRetrieve:   nil,
			team:          wantTeam,
			wantStatus:    http.StatusOK,
			assertFunc: func(t *testing.T, res *http.Response, _ string) {
				var team teamTable.Team
				if err := json.NewDecoder(res.Body).Decode(&team); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t.Error, team.ID, wantTeam.ID)
				assert.AllEqual(t.Error, team.Members, wantTeam.Members)
				for i, b := range wantTeam.Boards {
					assert.Equal(t.Error, team.Boards[i].ID, b.ID)
					assert.Equal(t.Error, team.Boards[i].Name, b.Name)
				}
			},
		},
		// if the team wasn't found and since we trust the JWT, we check whether
		// it tells us that the user is admin to determine whether to create a
		// new team
		{
			name:          "NotAdmin",
			auth:          "nonempty",
			errDecodeAuth: nil,
			authDecoded:   token.Auth{IsAdmin: false},
			errRetrieve:   db.ErrNoItem,
			team:          wantTeam,
			wantStatus:    http.StatusUnauthorized,
			assertFunc:    func(t *testing.T, res *http.Response, _ string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Err = c.errDecodeAuth
			decodeAuth.Res = c.authDecoded
			retriever.Err = c.errRetrieve
			retriever.Res = c.team

			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if c.auth != "" {
				r.AddCookie(&http.Cookie{Name: "auth-token", Value: c.auth})
			}
			w := httptest.NewRecorder()

			sut.Handle(w, r, "")

			res := w.Result()

			assert.Equal(t.Error, res.StatusCode, c.wantStatus)

			c.assertFunc(t, res, log.InMessage)
		})
	}
}
