//go:build utest

package team

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/teamtable"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
)

func TestGetHandler(t *testing.T) {
	authDecoder := &cookie.FakeDecoder[cookie.Auth]{}
	retriever := &db.FakeRetriever[teamtable.Team]{}
	inserter := &db.FakeInserter[teamtable.Team]{}
	log := &pkgLog.FakeErrorer{}
	sut := NewGetHandler(authDecoder, retriever, inserter, log)

	wantTeam := teamtable.Team{
		ID:      "teamid",
		Members: []string{"memberone", "membertwo"},
		Boards: []teamtable.Board{
			{ID: "board1", Name: "boardone"},
			{ID: "board2", Name: "boardtwo"},
		},
	}

	for _, c := range []struct {
		name          string
		auth          string
		errDecodeAuth error
		authDecoded   cookie.Auth
		errRetrieve   error
		team          teamtable.Team
		errInsert     error
		wantStatus    int
		assertFunc    func(*testing.T, *http.Response, string)
	}{
		{
			name:          "NoAuth",
			auth:          "",
			errDecodeAuth: nil,
			authDecoded:   cookie.Auth{},
			errRetrieve:   nil,
			team:          teamtable.Team{},
			errInsert:     nil,
			wantStatus:    http.StatusUnauthorized,
			assertFunc:    func(*testing.T, *http.Response, string) {},
		},
		{
			name:          "InvalidAuth",
			auth:          "nonempty",
			errDecodeAuth: errors.New("decode auth failed"),
			authDecoded:   cookie.Auth{},
			errRetrieve:   nil,
			team:          teamtable.Team{},
			errInsert:     nil,
			wantStatus:    http.StatusUnauthorized,
			assertFunc:    func(*testing.T, *http.Response, string) {},
		},
		{
			name:          "ErrRetrieve",
			auth:          "nonempty",
			errDecodeAuth: nil,
			authDecoded:   cookie.Auth{},
			errRetrieve:   errors.New("retrieve failed"),
			team:          teamtable.Team{},
			errInsert:     nil,
			wantStatus:    http.StatusInternalServerError,
			assertFunc:    assert.OnLoggedErr("retrieve failed"),
		},
		{
			name:          "OK",
			auth:          "nonempty",
			errDecodeAuth: nil,
			authDecoded:   cookie.Auth{},
			errRetrieve:   nil,
			team:          wantTeam,
			errInsert:     nil,
			wantStatus:    http.StatusOK,
			assertFunc: func(t *testing.T, res *http.Response, _ string) {
				var team teamtable.Team
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
		// check comments in implementation for explanation
		{
			name:          "NotAdmin",
			auth:          "nonempty",
			errDecodeAuth: nil,
			authDecoded:   cookie.Auth{IsAdmin: false},
			errRetrieve:   db.ErrNoItem,
			team:          teamtable.Team{},
			errInsert:     nil,
			wantStatus:    http.StatusUnauthorized,
			assertFunc:    func(t *testing.T, res *http.Response, _ string) {},
		},
		{
			name:          "ErrInsert",
			auth:          "nonempty",
			errDecodeAuth: nil,
			authDecoded:   cookie.Auth{IsAdmin: true},
			errRetrieve:   db.ErrNoItem,
			team:          teamtable.Team{},
			errInsert:     errors.New("insert failed"),
			wantStatus:    http.StatusInternalServerError,
			assertFunc:    assert.OnLoggedErr("insert failed"),
		},
		{
			name:          "Created",
			auth:          "nonempty",
			errDecodeAuth: nil,
			authDecoded:   cookie.Auth{IsAdmin: true, Username: "newuser"},
			errRetrieve:   db.ErrNoItem,
			team:          teamtable.Team{},
			errInsert:     nil,
			wantStatus:    http.StatusCreated,
			assertFunc: func(t *testing.T, res *http.Response, _ string) {
				var team teamtable.Team
				if err := json.NewDecoder(res.Body).Decode(&team); err != nil {
					t.Fatal(err)
				}

				assert.AllEqual(t.Error, team.Members, []string{"newuser"})
				assert.Equal(t.Error, len(team.Boards), 1)
				assert.Equal(t.Error, team.Boards[0].Name, "New Board")
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			authDecoder.Err = c.errDecodeAuth
			authDecoder.Res = c.authDecoded
			retriever.Err = c.errRetrieve
			retriever.Res = c.team
			inserter.Err = c.errInsert

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
