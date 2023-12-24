//go:build utest

package task

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/log"
)

// TestDeleteHandler tests the Handle method of DeleteHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDeleteHandler(t *testing.T) {
	authDecoder := &cookie.FakeDecoder[cookie.Auth]{}
	stateDecoder := &cookie.FakeDecoder[cookie.State]{}
	taskDeleter := &db.FakeDeleterDualKey{}
	stateEncoder := &cookie.FakeEncoder[cookie.State]{}
	log := &log.FakeErrorer{}
	sut := NewDeleteHandler(
		authDecoder,
		stateDecoder,
		taskDeleter,
		stateEncoder,
		log,
	)

	for _, c := range []struct {
		name           string
		authToken      string
		errDecodeAuth  error
		auth           cookie.Auth
		inState        string
		errDecodeState error
		inStateDecoded cookie.State
		errDeleteTask  error
		errEncodeState error
		outState       http.Cookie
		wantStatus     int
		assertFunc     func(*testing.T, *http.Response, []any)
	}{
		{
			name:           "NoAuth",
			authToken:      "",
			errDecodeAuth:  nil,
			auth:           cookie.Auth{},
			inStateDecoded: cookie.State{},
			inState:        "",
			errDecodeState: nil,
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatus:     http.StatusUnauthorized,
			assertFunc:     assert.OnResErr("Auth token not found."),
		},
		{
			name:           "ErrDecodeAuth",
			authToken:      "nonempty",
			errDecodeAuth:  errors.New("decode auth failed"),
			auth:           cookie.Auth{},
			inState:        "",
			errDecodeState: nil,
			inStateDecoded: cookie.State{},
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatus:     http.StatusUnauthorized,
			assertFunc:     assert.OnResErr("Invalid auth token."),
		},
		{
			name:           "NotAdmin",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           cookie.Auth{IsAdmin: false},
			inState:        "",
			errDecodeState: nil,
			inStateDecoded: cookie.State{},
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatus:     http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can delete tasks.",
			),
		},
		{
			name:           "NoState",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           cookie.Auth{IsAdmin: true},
			inState:        "",
			errDecodeState: nil,
			inStateDecoded: cookie.State{},
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatus:     http.StatusBadRequest,
			assertFunc:     assert.OnResErr("State token not found."),
		},
		{
			name:           "ErrDecodeState",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           cookie.Auth{IsAdmin: true},
			inState:        "nonempty",
			errDecodeState: errors.New("encode state failed"),
			inStateDecoded: cookie.State{},
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatus:     http.StatusBadRequest,
			assertFunc:     assert.OnResErr("Invalid state token."),
		},
		// if the ID is not found in state, it's invalid
		{
			name:           "IDInvalid",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           cookie.Auth{IsAdmin: true},
			inState:        "nonempty",
			errDecodeState: nil,
			inStateDecoded: cookie.State{},
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatus:     http.StatusBadRequest,
			assertFunc:     assert.OnResErr("Invalid task ID."),
		},
		{
			name:           "NotFound",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           cookie.Auth{IsAdmin: true},
			inState:        "nonempty",
			errDecodeState: nil,
			inStateDecoded: cookie.State{
				Boards: []cookie.Board{{Columns: []cookie.Column{{
					Tasks: []cookie.Task{{ID: "foo"}},
				}}}},
			},
			errDeleteTask:  db.ErrNoItem,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatus:     http.StatusNotFound,
			assertFunc:     assert.OnResErr("Task not found."),
		},
		{
			name:           "ErrDeleteTask",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           cookie.Auth{IsAdmin: true},
			inState:        "nonempty",
			errDecodeState: nil,
			inStateDecoded: cookie.State{
				Boards: []cookie.Board{{Columns: []cookie.Column{{
					Tasks: []cookie.Task{{ID: "foo"}},
				}}}},
			},
			errDeleteTask:  errors.New("delete task failed"),
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatus:     http.StatusInternalServerError,
			assertFunc:     assert.OnLoggedErr("delete task failed"),
		},
		{
			name:           "ErrEncodeState",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           cookie.Auth{IsAdmin: true},
			inState:        "nonempty",
			errDecodeState: nil,
			inStateDecoded: cookie.State{
				Boards: []cookie.Board{{Columns: []cookie.Column{{
					Tasks: []cookie.Task{{ID: "foo"}},
				}}}},
			},
			errDeleteTask:  nil,
			errEncodeState: errors.New("encode state failed"),
			outState:       http.Cookie{},
			wantStatus:     http.StatusInternalServerError,
			assertFunc:     assert.OnLoggedErr("encode state failed"),
		},
		{
			name:           "Success",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           cookie.Auth{IsAdmin: true},
			inState:        "nonempty",
			errDecodeState: nil,
			inStateDecoded: cookie.State{
				Boards: []cookie.Board{{Columns: []cookie.Column{{
					Tasks: []cookie.Task{{ID: "foo"}},
				}}}},
			},
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       http.Cookie{Name: "foo", Value: "bar"},
			wantStatus:     http.StatusOK,
			assertFunc: func(t *testing.T, resp *http.Response, _ []any) {
				ckState := resp.Cookies()[0]
				assert.Equal(t.Error, ckState.Name, "foo")
				assert.Equal(t.Error, ckState.Value, "bar")
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			authDecoder.Res = c.auth
			authDecoder.Err = c.errDecodeAuth
			stateDecoder.Res = c.inStateDecoded
			stateDecoder.Err = c.errDecodeState
			taskDeleter.Err = c.errDeleteTask
			stateEncoder.Res = c.outState
			stateEncoder.Err = c.errEncodeState

			r := httptest.NewRequest("", "/?id=foo", nil)
			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  "auth-token",
					Value: c.authToken,
				})
			}
			if c.inState != "" {
				r.AddCookie(&http.Cookie{
					Name:  "state-token",
					Value: c.authToken,
				})
			}

			w := httptest.NewRecorder()

			sut.Handle(w, r, "")
			resp := w.Result()

			assert.Equal(t.Error, resp.StatusCode, c.wantStatus)

			c.assertFunc(t, resp, log.Args)
		})
	}
}
