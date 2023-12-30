//go:build utest

package teamapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/teamtbl"
	"github.com/kxplxn/goteam/pkg/log"
)

func TestGetHandler(t *testing.T) {
	authDecoder := &cookie.FakeDecoder[cookie.Auth]{}
	teamRetriever := &db.FakeRetriever[teamtbl.Team]{}
	teamInserter := &db.FakeInserter[teamtbl.Team]{}
	teamUpdater := &db.FakeUpdater[teamtbl.Team]{}
	inviteEncoder := &cookie.FakeEncoder[cookie.Invite]{}
	log := &log.FakeErrorer{}
	sut := NewGetHandler(
		authDecoder,
		teamRetriever,
		teamInserter,
		teamUpdater,
		inviteEncoder,
		log,
	)

	wantTeam := teamtbl.Team{
		ID:      "teamid",
		Members: []string{"memberone", "membertwo"},
		Boards: []teamtbl.Board{
			{ID: "board1", Name: "boardone"},
			{ID: "board2", Name: "boardtwo"},
		},
	}

	for _, c := range []struct {
		name            string
		auth            string
		errDecodeAuth   error
		authDecoded     cookie.Auth
		errRetrieve     error
		team            teamtbl.Team
		errInsert       error
		errUpdate       error
		errEncodeInvite error
		inviteEncoded   http.Cookie
		wantStatus      int
		assertFunc      func(*testing.T, *http.Response, []any)
	}{
		{
			name:            "NoAuth",
			auth:            "",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{},
			errRetrieve:     nil,
			team:            teamtbl.Team{},
			errInsert:       nil,
			errUpdate:       nil,
			errEncodeInvite: nil,
			inviteEncoded:   http.Cookie{},
			wantStatus:      http.StatusUnauthorized,
			assertFunc:      func(*testing.T, *http.Response, []any) {},
		},
		{
			name:            "InvalidAuth",
			auth:            "nonempty",
			errDecodeAuth:   errors.New("decode auth failed"),
			authDecoded:     cookie.Auth{},
			errRetrieve:     nil,
			team:            teamtbl.Team{},
			errInsert:       nil,
			errUpdate:       nil,
			errEncodeInvite: nil,
			inviteEncoded:   http.Cookie{},
			wantStatus:      http.StatusUnauthorized,
			assertFunc:      func(*testing.T, *http.Response, []any) {},
		},
		{
			name:            "ErrRetrieve",
			auth:            "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{},
			errRetrieve:     errors.New("retrieve failed"),
			team:            teamtbl.Team{},
			errInsert:       nil,
			errUpdate:       nil,
			errEncodeInvite: nil,
			inviteEncoded:   http.Cookie{},
			wantStatus:      http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr("retrieve failed"),
		},
		// check comments in implementation for explanation
		{
			name:            "NotAdmin",
			auth:            "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: false},
			errRetrieve:     db.ErrNoItem,
			team:            teamtbl.Team{},
			errInsert:       nil,
			errUpdate:       nil,
			errEncodeInvite: nil,
			inviteEncoded:   http.Cookie{},
			wantStatus:      http.StatusUnauthorized,
			assertFunc:      func(*testing.T, *http.Response, []any) {},
		},
		{
			name:            "ErrInsert",
			auth:            "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errRetrieve:     db.ErrNoItem,
			team:            teamtbl.Team{},
			errInsert:       errors.New("insert failed"),
			errUpdate:       nil,
			errEncodeInvite: nil,
			inviteEncoded:   http.Cookie{},
			wantStatus:      http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr("insert failed"),
		},
		{
			name:            "ErrUpdate",
			auth:            "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: false},
			errRetrieve:     nil,
			team:            teamtbl.Team{},
			errInsert:       nil,
			errUpdate:       errors.New("update failed"),
			errEncodeInvite: nil,
			inviteEncoded:   http.Cookie{},
			wantStatus:      http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr("update failed"),
		},
		{
			name:            "ErrEncodeInvite",
			auth:            "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true, Username: "memberone"},
			errRetrieve:     nil,
			team:            teamtbl.Team{Members: []string{"memberone"}},
			errInsert:       nil,
			errUpdate:       nil,
			errEncodeInvite: errors.New("encode invite failed"),
			inviteEncoded:   http.Cookie{},
			wantStatus:      http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr("encode invite failed"),
		},
		{
			name:            "OK",
			auth:            "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true, Username: "memberone"},
			errRetrieve:     nil,
			team:            wantTeam,
			errInsert:       nil,
			errUpdate:       nil,
			errEncodeInvite: nil,
			inviteEncoded:   http.Cookie{Name: "invite-token", Value: "aksdfj"},
			wantStatus:      http.StatusOK,
			assertFunc: func(t *testing.T, resp *http.Response, _ []any) {
				var team teamtbl.Team
				if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t.Error, team.ID, wantTeam.ID)
				assert.AllEqual(t.Error, team.Members, wantTeam.Members)
				for i, b := range wantTeam.Boards {
					assert.Equal(t.Error, team.Boards[i].ID, b.ID)
					assert.Equal(t.Error, team.Boards[i].Name, b.Name)
				}

				assert.Equal(t.Error, len(resp.Cookies()), 0)
			},
		},
		{
			name:            "OKNewTeam",
			auth:            "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true, Username: "newuser"},
			errRetrieve:     db.ErrNoItem,
			team:            teamtbl.Team{},
			errInsert:       nil,
			errUpdate:       nil,
			errEncodeInvite: nil,
			inviteEncoded:   http.Cookie{Name: "invite-token", Value: "aksdfj"},
			wantStatus:      http.StatusCreated,
			assertFunc: func(t *testing.T, resp *http.Response, _ []any) {
				var team teamtbl.Team

				if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
					t.Fatal(err)
				}

				assert.AllEqual(t.Error, team.Members, []string{"newuser"})
				assert.Equal(t.Error, len(team.Boards), 1)
				assert.Equal(t.Error, team.Boards[0].Name, "New Board")

				ckInv := resp.Cookies()[0]
				assert.Equal(t.Error, ckInv.Name, "invite-token")
				assert.Equal(t.Error, ckInv.Value, "aksdfj")
			},
		},
		{
			name:            "OKNewMember",
			auth:            "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: false, Username: "newuser"},
			errRetrieve:     nil,
			team:            wantTeam,
			errInsert:       nil,
			errUpdate:       nil,
			errEncodeInvite: nil,
			inviteEncoded:   http.Cookie{},
			wantStatus:      http.StatusOK,
			assertFunc: func(t *testing.T, resp *http.Response, _ []any) {
				var team teamtbl.Team
				if err := json.NewDecoder(resp.Body).Decode(&team); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t.Error, team.ID, wantTeam.ID)
				assert.AllEqual(t.Error,
					team.Members, append(wantTeam.Members, "newuser"),
				)
				for i, b := range wantTeam.Boards {
					assert.Equal(t.Error, team.Boards[i].ID, b.ID)
					assert.Equal(t.Error, team.Boards[i].Name, b.Name)
				}

				assert.Equal(t.Error, len(resp.Cookies()), 0)
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			authDecoder.Err = c.errDecodeAuth
			authDecoder.Res = c.authDecoded
			teamRetriever.Err = c.errRetrieve
			teamRetriever.Res = c.team
			teamInserter.Err = c.errInsert
			teamUpdater.Err = c.errUpdate
			inviteEncoder.Err = c.errEncodeInvite
			inviteEncoder.Res = c.inviteEncoded
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/", nil)
			if c.auth != "" {
				r.AddCookie(&http.Cookie{Name: "auth-token", Value: c.auth})
			}

			sut.Handle(w, r, "")

			resp := w.Result()
			assert.Equal(t.Error, resp.StatusCode, c.wantStatus)
			c.assertFunc(t, resp, log.Args)
		})
	}
}
