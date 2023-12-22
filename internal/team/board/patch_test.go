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
	"github.com/kxplxn/goteam/pkg/db"
	teamTable "github.com/kxplxn/goteam/pkg/db/teamtable"
	boardTable "github.com/kxplxn/goteam/pkg/legacydb/board"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
	"github.com/kxplxn/goteam/pkg/validator"
)

// TestPATCHHandler tests the Handle method of PATCHHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	decodeAuth := &token.FakeDecode[token.Auth]{}
	decodeState := &token.FakeDecode[token.State]{}
	idValidator := &api.FakeStringValidator{}
	nameValidator := &api.FakeStringValidator{}
	updater := &db.FakeUpdaterDualKey[teamTable.Board]{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPatchHandler(
		decodeAuth.Func,
		decodeState.Func,
		idValidator,
		nameValidator,
		updater,
		log,
	)

	for _, c := range []struct {
		name            string
		authToken       string
		errDecodeAuth   error
		authDecoded     token.Auth
		stateToken      string
		errDecodeState  error
		stateDecoded    token.State
		errValidateID   error
		errValidateName error
		board           boardTable.Record
		boardUpdaterErr error
		wantStatusCode  int
		assertFunc      func(*testing.T, *http.Response, string)
	}{
		{
			name:            "NoAuth",
			authToken:       "",
			errDecodeAuth:   nil,
			authDecoded:     token.Auth{},
			stateToken:      "",
			errDecodeState:  nil,
			stateDecoded:    token.State{},
			errValidateID:   nil,
			errValidateName: nil,
			board:           boardTable.Record{},
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusUnauthorized,
			assertFunc:      assert.OnResErr("Auth token not found."),
		},
		{
			name:            "InvalidAuth",
			authToken:       "nonempty",
			errDecodeAuth:   token.ErrInvalid,
			authDecoded:     token.Auth{},
			stateToken:      "",
			errDecodeState:  nil,
			stateDecoded:    token.State{},
			errValidateID:   nil,
			errValidateName: nil,
			board:           boardTable.Record{},
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusUnauthorized,
			assertFunc:      assert.OnResErr("Invalid auth token."),
		},
		{
			name:            "NotAdmin",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     token.Auth{IsAdmin: false},
			stateToken:      "",
			errDecodeState:  nil,
			stateDecoded:    token.State{},
			errValidateID:   nil,
			errValidateName: nil,
			board:           boardTable.Record{},
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can edit boards.",
			),
		},
		{
			name:            "NoState",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     token.Auth{IsAdmin: true},
			stateToken:      "",
			errDecodeState:  nil,
			stateDecoded:    token.State{},
			errValidateID:   nil,
			errValidateName: nil,
			board:           boardTable.Record{},
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc:      assert.OnResErr("State token not found."),
		},
		{
			name:            "InvalidState",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     token.Auth{IsAdmin: true},
			stateToken:      "nonempty",
			errDecodeState:  token.ErrInvalid,
			stateDecoded:    token.State{},
			errValidateID:   nil,
			errValidateName: nil,
			board:           boardTable.Record{},
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc:      assert.OnResErr("Invalid state token."),
		},
		{
			name:            "IDEmpty",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     token.Auth{IsAdmin: true},
			stateToken:      "nonempty",
			errDecodeState:  nil,
			stateDecoded:    token.State{},
			errValidateID:   validator.ErrEmpty,
			errValidateName: nil,
			board:           boardTable.Record{},
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc:      assert.OnResErr("Board ID cannot be empty."),
		},
		{
			name:            "IDNotUUID",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     token.Auth{IsAdmin: true},
			stateToken:      "nonempty",
			errDecodeState:  nil,
			stateDecoded:    token.State{},
			errValidateID:   validator.ErrWrongFormat,
			errValidateName: nil,
			board:           boardTable.Record{},
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc:      assert.OnResErr("Board ID must be a UUID."),
		},
		{
			name:           "NameEmpty",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: token.State{Boards: []token.Board{
				{ID: "c193d6ba-ebfe-45fe-80d9-00b545690b4b"},
			}},
			errValidateID:   nil,
			errValidateName: validator.ErrEmpty,
			board:           boardTable.Record{},
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc:      assert.OnResErr("Board name cannot be empty."),
		},
		{
			name:           "NameTooLong",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: token.State{Boards: []token.Board{
				{ID: "c193d6ba-ebfe-45fe-80d9-00b545690b4b"},
			}},
			errValidateID:   nil,
			errValidateName: validator.ErrTooLong,
			board:           boardTable.Record{},
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Board name cannot be longer than 35 characters.",
			),
		},
		{
			name:            "NoAccess",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			authDecoded:     token.Auth{IsAdmin: true},
			stateToken:      "nonempty",
			errDecodeState:  nil,
			stateDecoded:    token.State{},
			errValidateName: nil,
			board:           boardTable.Record{},
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:           "BoardNotFound",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: token.State{Boards: []token.Board{
				{ID: "c193d6ba-ebfe-45fe-80d9-00b545690b4b"},
			}},
			errValidateName: nil,
			board:           boardTable.Record{},
			boardUpdaterErr: db.ErrNoItem,
			wantStatusCode:  http.StatusNotFound,
			assertFunc:      assert.OnResErr("Board not found."),
		},
		{
			name:           "BoardUpdaterErr",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: token.State{Boards: []token.Board{
				{ID: "c193d6ba-ebfe-45fe-80d9-00b545690b4b"},
			}},
			errValidateName: nil,
			board:           boardTable.Record{},
			boardUpdaterErr: errors.New("update board failed"),
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr("update board failed"),
		},
		{
			name:           "Success",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: token.State{Boards: []token.Board{
				{ID: "c193d6ba-ebfe-45fe-80d9-00b545690b4b"},
			}},
			errValidateName: nil,
			board:           boardTable.Record{},
			boardUpdaterErr: nil,
			wantStatusCode:  http.StatusOK,
			assertFunc:      func(*testing.T, *http.Response, string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Err = c.errDecodeAuth
			decodeAuth.Res = c.authDecoded
			decodeState.Err = c.errDecodeState
			decodeState.Res = c.stateDecoded
			idValidator.Err = c.errValidateID
			nameValidator.Err = c.errValidateName
			updater.Err = c.boardUpdaterErr

			w := httptest.NewRecorder()
			r := httptest.NewRequest("", "/", strings.NewReader(`{
                "id": "c193d6ba-ebfe-45fe-80d9-00b545690b4b"
            }`))

			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  token.AuthName,
					Value: c.authToken,
				})
			}
			if c.stateToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  token.StateName,
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
