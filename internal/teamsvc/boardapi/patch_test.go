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

func TestPatchHandler(t *testing.T) {
	decodeAuth := &cookie.FakeDecoder[cookie.Auth]{}
	decodeState := &cookie.FakeDecoder[cookie.State]{}
	idValidator := &api.FakeStringValidator{}
	nameValidator := &api.FakeStringValidator{}
	updater := &db.FakeUpdaterDualKey[teamtbl.Board]{}
	log := &log.FakeErrorer{}
	sut := NewPatchHandler(
		decodeAuth,
		decodeState,
		idValidator,
		nameValidator,
		updater,
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
		errValidateID   error
		errValidateName error
		errUpdateBoard  error
		wantStatus      int
		assertFunc      func(*testing.T, *http.Response, []any)
	}{
		{
			name:            "NoAuth",
			authToken:       "",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{},
			stateToken:      "",
			errDecodeState:  nil,
			stateDecoded:    cookie.State{},
			errValidateID:   nil,
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusUnauthorized,
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
			errValidateID:   nil,
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusUnauthorized,
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
			errValidateID:   nil,
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusForbidden,
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
			errValidateID:   nil,
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusForbidden,
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
			errValidateID:   nil,
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusForbidden,
			assertFunc:      assert.OnResErr("Invalid state token."),
		},
		{
			name:            "IDEmpty",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			stateToken:      "nonempty",
			errDecodeState:  nil,
			stateDecoded:    cookie.State{},
			errValidateID:   validator.ErrEmpty,
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc:      assert.OnResErr("Board ID cannot be empty."),
		},
		{
			name:            "IDNotUUID",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			stateToken:      "nonempty",
			errDecodeState:  nil,
			stateDecoded:    cookie.State{},
			errValidateID:   validator.ErrWrongFormat,
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc:      assert.OnResErr("Board ID must be a UUID."),
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
			errValidateID:   nil,
			errValidateName: validator.ErrEmpty,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusBadRequest,
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
			errValidateID:   nil,
			errValidateName: validator.ErrTooLong,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Board name cannot be longer than 35 characters.",
			),
		},
		{
			name:            "NoAccess",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			stateToken:      "nonempty",
			errDecodeState:  nil,
			stateDecoded:    cookie.State{},
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
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
			errUpdateBoard:  db.ErrNoItem,
			wantStatus:      http.StatusNotFound,
			assertFunc:      assert.OnResErr("Board not found."),
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
			errUpdateBoard:  errors.New("update board failed"),
			wantStatus:      http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr("update board failed"),
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
			errUpdateBoard:  nil,
			wantStatus:      http.StatusOK,
			assertFunc:      func(*testing.T, *http.Response, []any) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Err = c.errDecodeAuth
			decodeAuth.Res = c.authDecoded
			decodeState.Err = c.errDecodeState
			decodeState.Res = c.stateDecoded
			idValidator.Err = c.errValidateID
			nameValidator.Err = c.errValidateName
			updater.Err = c.errUpdateBoard

			w := httptest.NewRecorder()
			r := httptest.NewRequest("", "/", strings.NewReader(`{
                "id": "c193d6ba-ebfe-45fe-80d9-00b545690b4b"
            }`))

			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  cookie.AuthName,
					Value: c.authToken,
				})
			}
			if c.stateToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  cookie.StateName,
					Value: c.stateToken,
				})
			}

			sut.Handle(w, r, "")
			res := w.Result()

			assert.Equal(t.Error, res.StatusCode, c.wantStatus)

			c.assertFunc(t, res, log.Args)
		})
	}
}
