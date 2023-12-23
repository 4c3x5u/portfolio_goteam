//go:build utest

package board

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/pkg/api"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/teamtable"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/validator"
)

func TestPostHandler(t *testing.T) {
	decodeAuth := &cookie.FakeDecoder[cookie.Auth]{}
	decodeState := &cookie.FakeDecoder[cookie.State]{}
	nameValidator := &api.FakeStringValidator{}
	inserter := &db.FakeInserterDualKey[teamtable.Board]{}
	encodeState := &cookie.FakeEncoder[cookie.State]{}
	log := &log.FakeErrorer{}
	sut := NewPostHandler(
		decodeAuth,
		decodeState,
		nameValidator,
		inserter,
		encodeState,
		log,
	)

	for _, c := range []struct {
		name            string
		authToken       string
		errDecodeAuth   error
		authDecoded     cookie.Auth
		stateToken      string
		errDecodeState  error
		stateDecoded    cookie.State
		errValidateName error
		boardUpdaterErr error
		outStateToken   http.Cookie
		errEncodeState  error
		wantStatusCode  int
		assertFunc      func(*testing.T, *http.Response, string)
	}{
		{
			name:            "NoAuth",
			authToken:       "",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{},
			stateToken:      "",
			errDecodeState:  nil,
			stateDecoded:    cookie.State{},
			errValidateName: nil,
			boardUpdaterErr: nil,
			outStateToken:   http.Cookie{},
			errEncodeState:  nil,
			wantStatusCode:  http.StatusUnauthorized,
			assertFunc:      assert.OnResErr("Auth token not found."),
		},
		{
			name:            "InvalidAuth",
			authToken:       "nonempty",
			errDecodeAuth:   cookie.ErrInvalid,
			authDecoded:     cookie.Auth{},
			stateToken:      "",
			errDecodeState:  nil,
			stateDecoded:    cookie.State{},
			errValidateName: nil,
			boardUpdaterErr: nil,
			outStateToken:   http.Cookie{},
			errEncodeState:  nil,
			wantStatusCode:  http.StatusUnauthorized,
			assertFunc:      assert.OnResErr("Invalid auth token."),
		},
		{
			name:            "NotAdmin",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: false},
			stateToken:      "",
			errDecodeState:  nil,
			stateDecoded:    cookie.State{},
			errValidateName: nil,
			boardUpdaterErr: nil,
			outStateToken:   http.Cookie{},
			errEncodeState:  nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can edit boards.",
			),
		},
		{
			name:            "NoState",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			stateToken:      "",
			errDecodeState:  nil,
			stateDecoded:    cookie.State{},
			errValidateName: nil,
			boardUpdaterErr: nil,
			outStateToken:   http.Cookie{},
			errEncodeState:  nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc:      assert.OnResErr("State token not found."),
		},
		{
			name:            "InvalidState",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			stateToken:      "nonempty",
			errDecodeState:  cookie.ErrInvalid,
			stateDecoded:    cookie.State{},
			errValidateName: nil,
			boardUpdaterErr: nil,
			outStateToken:   http.Cookie{},
			errEncodeState:  nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc:      assert.OnResErr("Invalid state token."),
		},
		{
			name:            "LimitReached",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			stateToken:      "nonempty",
			errDecodeState:  nil,
			stateDecoded:    cookie.State{Boards: []cookie.Board{{}, {}, {}}},
			errValidateName: nil,
			boardUpdaterErr: nil,
			outStateToken:   http.Cookie{},
			errEncodeState:  nil,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"You have already created the maximum amount of boards " +
					"allowed per team. Please delete one of your boards to " +
					"create a new one.",
			),
		},
		{
			name:           "NameEmpty",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: cookie.State{Boards: []cookie.Board{
				{ID: "c193d6ba-ebfe-45fe-80d9-00b545690b4b"},
			}},
			errValidateName: validator.ErrEmpty,
			boardUpdaterErr: nil,
			outStateToken:   http.Cookie{},
			errEncodeState:  nil,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc:      assert.OnResErr("Board name cannot be empty."),
		},
		{
			name:           "NameTooLong",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: cookie.State{Boards: []cookie.Board{
				{ID: "c193d6ba-ebfe-45fe-80d9-00b545690b4b"},
			}},
			errValidateName: validator.ErrTooLong,
			boardUpdaterErr: nil,
			outStateToken:   http.Cookie{},
			errEncodeState:  nil,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Board name cannot be longer than 35 characters.",
			),
		},
		{
			name:           "BoardNotFound",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: cookie.State{Boards: []cookie.Board{
				{ID: "c193d6ba-ebfe-45fe-80d9-00b545690b4b"},
			}},
			errValidateName: nil,
			boardUpdaterErr: db.ErrLimitReached,
			outStateToken:   http.Cookie{},
			errEncodeState:  nil,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"You have already created the maximum amount of boards " +
					"allowed per team. Please delete one of your boards to " +
					"create a new one.",
			),
		},
		{
			name:           "BoardUpdaterErr",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: cookie.State{Boards: []cookie.Board{
				{ID: "c193d6ba-ebfe-45fe-80d9-00b545690b4b"},
			}},
			errValidateName: nil,
			boardUpdaterErr: errors.New("update board failed"),
			outStateToken:   http.Cookie{},
			errEncodeState:  nil,
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr("update board failed"),
		},
		{
			name:           "ErrEncodeState",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: cookie.State{Boards: []cookie.Board{
				{ID: "c193d6ba-ebfe-45fe-80d9-00b545690b4b"},
			}},
			errValidateName: nil,
			boardUpdaterErr: nil,
			outStateToken:   http.Cookie{},
			errEncodeState:  errors.New("encode state failed"),
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr("encode state failed"),
		},
		{
			name:           "Success",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: cookie.State{Boards: []cookie.Board{
				{ID: "c193d6ba-ebfe-45fe-80d9-00b545690b4b"},
			}},
			errValidateName: nil,
			boardUpdaterErr: nil,
			outStateToken:   http.Cookie{Name: "foo", Value: "bar"},
			errEncodeState:  nil,
			wantStatusCode:  http.StatusOK,
			assertFunc: func(t *testing.T, resp *http.Response, _ string) {
				ck := resp.Cookies()[0]
				assert.Equal(t.Error, ck.Name, "foo")
				assert.Equal(t.Error, ck.Value, "bar")
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Err = c.errDecodeAuth
			decodeAuth.Res = c.authDecoded
			decodeState.Err = c.errDecodeState
			decodeState.Res = c.stateDecoded
			nameValidator.Err = c.errValidateName
			inserter.Err = c.boardUpdaterErr
			encodeState.Res = c.outStateToken
			encodeState.Err = c.errEncodeState

			w := httptest.NewRecorder()
			r := httptest.NewRequest("", "/", strings.NewReader(`{
                "id": "c193d6ba-ebfe-45fe-80d9-00b545690b4b"
            }`))

			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  "auth-token",
					Value: c.authToken,
				})
			}
			if c.stateToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  "state-token",
					Value: c.stateToken,
				})
			}

			sut.Handle(w, r, "")
			res := w.Result()

			assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

			c.assertFunc(t, res, log.InMessage)
		})
	}
}
