package team

import (
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
