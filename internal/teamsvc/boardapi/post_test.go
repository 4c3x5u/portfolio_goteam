//go:build utest

package boardapi

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
	"github.com/kxplxn/goteam/pkg/db/teamtbl"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/validator"
)

func TestPostHandler(t *testing.T) {
	decodeAuth := &cookie.FakeDecoder[cookie.Auth]{}
	nameValidator := &api.FakeStringValidator{}
	inserter := &db.FakeInserterDualKey[teamtbl.Board]{}
	log := &log.FakeErrorer{}
	sut := NewPostHandler(decodeAuth, nameValidator, inserter, log)

	for _, c := range []struct {
		name            string
		authToken       string
		errDecodeAuth   error
		authDecoded     cookie.Auth
		errValidateName error
		boardUpdaterErr error
		wantStatusCode  int
		assertFunc      func(*testing.T, *http.Response, []any)
	}{
		{
			name:            "NoAuth",
			authToken:       "",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{},
			errValidateName: nil,
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusUnauthorized,
			assertFunc:      assert.OnRespErr("Auth token not found."),
		},
		{
			name:            "InvalidAuth",
			authToken:       "nonempty",
			errDecodeAuth:   cookie.ErrInvalid,
			authDecoded:     cookie.Auth{},
			errValidateName: nil,
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusUnauthorized,
			assertFunc:      assert.OnRespErr("Invalid auth token."),
		},
		{
			name:            "NotAdmin",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: false},
			errValidateName: nil,
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc: assert.OnRespErr(
				"Only team admins can edit boards.",
			),
		},
		{
			name:            "NameEmpty",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errValidateName: validator.ErrEmpty,
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc:      assert.OnRespErr("Board name cannot be empty."),
		},
		{
			name:            "NameTooLong",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errValidateName: validator.ErrTooLong,
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Board name cannot be longer than 35 characters.",
			),
		},
		{
			name:            "ErrLimitReached",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errValidateName: nil,
			boardUpdaterErr: db.ErrLimitReached,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"You have already created the maximum amount of boards " +
					"allowed per team. Please delete one of your boards to " +
					"create a new one.",
			),
		},
		{
			name:            "BoardUpdaterErr",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errValidateName: nil,
			boardUpdaterErr: errors.New("update board failed"),
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr("update board failed"),
		},
		{
			name:            "Success",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errValidateName: nil,
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusOK,
			assertFunc:      func(*testing.T, *http.Response, []any) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Err = c.errDecodeAuth
			decodeAuth.Res = c.authDecoded
			nameValidator.Err = c.errValidateName
			inserter.Err = c.boardUpdaterErr
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

			sut.Handle(w, r, "")

			resp := w.Result()
			assert.Equal(t.Error, resp.StatusCode, c.wantStatusCode)
			c.assertFunc(t, resp, log.Args)
		})
	}
}
