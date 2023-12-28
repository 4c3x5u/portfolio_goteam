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
	idValidator := &api.FakeStringValidator{}
	nameValidator := &api.FakeStringValidator{}
	updater := &db.FakeUpdaterDualKey[teamtbl.Board]{}
	log := &log.FakeErrorer{}
	sut := NewPatchHandler(
		decodeAuth,
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
			errValidateID:   nil,
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusUnauthorized,
			assertFunc:      assert.OnRespErr("Auth token not found."),
		},
		{
			name:            "InvalidAuth",
			authToken:       "nonempty",
			errDecodeAuth:   cookie.ErrInvalid,
			authDecoded:     cookie.Auth{},
			errValidateID:   nil,
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusUnauthorized,
			assertFunc:      assert.OnRespErr("Invalid auth token."),
		},
		{
			name:            "NotAdmin",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: false},
			errValidateID:   nil,
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusForbidden,
			assertFunc: assert.OnRespErr(
				"Only team admins can edit boards.",
			),
		},
		{
			name:            "IDEmpty",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errValidateID:   validator.ErrEmpty,
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc:      assert.OnRespErr("Board ID cannot be empty."),
		},
		{
			name:            "IDNotUUID",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errValidateID:   validator.ErrWrongFormat,
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc:      assert.OnRespErr("Board ID must be a UUID."),
		},
		{
			name:            "NameEmpty",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errValidateID:   nil,
			errValidateName: validator.ErrEmpty,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc:      assert.OnRespErr("Board name cannot be empty."),
		},
		{
			name:            "NameTooLong",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errValidateID:   nil,
			errValidateName: validator.ErrTooLong,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Board name cannot be longer than 35 characters.",
			),
		},
		{
			name:            "BoardNotFound",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errValidateName: nil,
			errUpdateBoard:  db.ErrNoItem,
			wantStatus:      http.StatusNotFound,
			assertFunc:      assert.OnRespErr("Board not found."),
		},
		{
			name:            "BoardUpdaterErr",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errValidateName: nil,
			errUpdateBoard:  errors.New("update board failed"),
			wantStatus:      http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr("update board failed"),
		},
		{
			name:            "Success",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     cookie.Auth{IsAdmin: true},
			errValidateName: nil,
			errUpdateBoard:  nil,
			wantStatus:      http.StatusOK,
			assertFunc:      func(*testing.T, *http.Response, []any) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Err = c.errDecodeAuth
			decodeAuth.Res = c.authDecoded
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

			sut.Handle(w, r, "")

			resp := w.Result()
			assert.Equal(t.Error, resp.StatusCode, c.wantStatus)
			c.assertFunc(t, resp, log.Args)
		})
	}
}
