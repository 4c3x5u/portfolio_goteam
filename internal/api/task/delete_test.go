//go:build utest

package task

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/dbaccess"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// TestDeleteHandler tests the Handle method of DeleteHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDeleteHandler(t *testing.T) {
	decodeAuth := &token.FakeDecode[token.Auth]{}
	decodeState := &token.FakeDecode[token.State]{}
	taskDeleter := &dbaccess.FakeDeleter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDeleteHandler(
		decodeAuth.Func,
		decodeState.Func,
		taskDeleter,
		log,
	)

	for _, c := range []struct {
		name           string
		authToken      string
		errDecodeAuth  error
		auth           token.Auth
		stateToken     string
		errDecodeState error
		state          token.State
		errDeleteTask  error
		wantStatus     int
		assertFunc     func(*testing.T, *http.Response, string)
	}{
		{
			name:           "NoAuth",
			authToken:      "",
			errDecodeAuth:  nil,
			auth:           token.Auth{},
			state:          token.State{},
			stateToken:     "",
			errDecodeState: nil,
			errDeleteTask:  nil,
			wantStatus:     http.StatusUnauthorized,
			assertFunc:     assert.OnResErr("Auth token not found."),
		},
		{
			name:           "ErrEncodeAuth",
			authToken:      "nonempty",
			errDecodeAuth:  errors.New("encode auth failed"),
			auth:           token.Auth{},
			stateToken:     "",
			errDecodeState: nil,
			state:          token.State{},
			errDeleteTask:  nil,
			wantStatus:     http.StatusUnauthorized,
			assertFunc:     assert.OnResErr("Invalid auth token."),
		},
		{
			name:           "NotAdmin",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: false},
			stateToken:     "",
			errDecodeState: nil,
			state:          token.State{},
			errDeleteTask:  nil,
			wantStatus:     http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only board admins can delete tasks.",
			),
		},
		{
			name:           "NoState",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: true},
			stateToken:     "",
			errDecodeState: nil,
			state:          token.State{},
			errDeleteTask:  nil,
			wantStatus:     http.StatusBadRequest,
			assertFunc:     assert.OnResErr("State token not found."),
		},
		{
			name:           "ErrDecodeState",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: errors.New("encode state failed"),
			state:          token.State{},
			errDeleteTask:  nil,
			wantStatus:     http.StatusBadRequest,
			assertFunc:     assert.OnResErr("Invalid state token."),
		},
		// if the ID is not found in state, it's invalid
		{
			name:           "IDInvalid",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			state:          token.State{},
			errDeleteTask:  nil,
			wantStatus:     http.StatusBadRequest,
			assertFunc:     assert.OnResErr("Invalid task ID."),
		},
		{
			name:           "ErrDeleteTask",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			state: token.State{
				Boards: []token.Board{{Columns: []token.Column{{
					Tasks: []token.Task{{ID: "foo"}},
				}}}},
			},
			errDeleteTask: sql.ErrNoRows,
			wantStatus:    http.StatusInternalServerError,
			assertFunc:    assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:           "Success",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			state: token.State{
				Boards: []token.Board{{Columns: []token.Column{{
					Tasks: []token.Task{{ID: "foo"}},
				}}}},
			},
			errDeleteTask: nil,
			wantStatus:    http.StatusOK,
			assertFunc:    func(*testing.T, *http.Response, string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Decoded = c.auth
			decodeAuth.Err = c.errDecodeAuth
			decodeState.Decoded = c.state
			decodeState.Err = c.errDecodeState
			taskDeleter.Err = c.errDeleteTask

			r := httptest.NewRequest("", "/?id=foo", nil)
			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  "auth-token",
					Value: c.authToken,
				})
			}
			if c.stateToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  "state-token",
					Value: c.authToken,
				})
			}

			w := httptest.NewRecorder()

			sut.Handle(w, r, "")
			res := w.Result()

			assert.Equal(t.Error, res.StatusCode, c.wantStatus)

			c.assertFunc(t, res, log.InMessage)
		})
	}
}
