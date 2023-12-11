//go:build utest

package task

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// TestDeleteHandler tests the Handle method of DeleteHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDeleteHandler(t *testing.T) {
	decodeAuth := &token.FakeDecode[token.Auth]{}
	decodeState := &token.FakeDecode[token.State]{}
	taskDeleter := &db.FakeDeleter{}
	encodeState := &token.FakeEncode[token.State]{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDeleteHandler(
		decodeAuth.Func,
		decodeState.Func,
		taskDeleter,
		encodeState.Func,
		log,
	)

	for _, c := range []struct {
		name           string
		authToken      string
		errDecodeAuth  error
		auth           token.Auth
		inState        string
		errDecodeState error
		inStateDecoded token.State
		errDeleteTask  error
		errEncodeState error
		outState       string
		wantStatus     int
		assertFunc     func(*testing.T, *http.Response, string)
	}{
		{
			name:           "NoAuth",
			authToken:      "",
			errDecodeAuth:  nil,
			auth:           token.Auth{},
			inStateDecoded: token.State{},
			inState:        "",
			errDecodeState: nil,
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       "",
			wantStatus:     http.StatusUnauthorized,
			assertFunc:     assert.OnResErr("Auth token not found."),
		},
		{
			name:           "ErrDecodeAuth",
			authToken:      "nonempty",
			errDecodeAuth:  errors.New("decode auth failed"),
			auth:           token.Auth{},
			inState:        "",
			errDecodeState: nil,
			inStateDecoded: token.State{},
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       "",
			wantStatus:     http.StatusUnauthorized,
			assertFunc:     assert.OnResErr("Invalid auth token."),
		},
		{
			name:           "NotAdmin",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: false},
			inState:        "",
			errDecodeState: nil,
			inStateDecoded: token.State{},
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       "",
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
			inState:        "",
			errDecodeState: nil,
			inStateDecoded: token.State{},
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       "",
			wantStatus:     http.StatusBadRequest,
			assertFunc:     assert.OnResErr("State token not found."),
		},
		{
			name:           "ErrDecodeState",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: true},
			inState:        "nonempty",
			errDecodeState: errors.New("encode state failed"),
			inStateDecoded: token.State{},
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       "",
			wantStatus:     http.StatusBadRequest,
			assertFunc:     assert.OnResErr("Invalid state token."),
		},
		// if the ID is not found in state, it's invalid
		{
			name:           "IDInvalid",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: true},
			inState:        "nonempty",
			errDecodeState: nil,
			inStateDecoded: token.State{},
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       "",
			wantStatus:     http.StatusBadRequest,
			assertFunc:     assert.OnResErr("Invalid task ID."),
		},
		{
			name:           "NotFound",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: true},
			inState:        "nonempty",
			errDecodeState: nil,
			inStateDecoded: token.State{
				Boards: []token.Board{{Columns: []token.Column{{
					Tasks: []token.Task{{ID: "foo"}},
				}}}},
			},
			errDeleteTask:  db.ErrNoItem,
			errEncodeState: nil,
			outState:       "",
			wantStatus:     http.StatusNotFound,
			assertFunc:     assert.OnResErr("Task not found."),
		},
		{
			name:           "ErrDeleteTask",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: true},
			inState:        "nonempty",
			errDecodeState: nil,
			inStateDecoded: token.State{
				Boards: []token.Board{{Columns: []token.Column{{
					Tasks: []token.Task{{ID: "foo"}},
				}}}},
			},
			errDeleteTask:  errors.New("delete task failed"),
			errEncodeState: nil,
			outState:       "",
			wantStatus:     http.StatusInternalServerError,
			assertFunc:     assert.OnLoggedErr("delete task failed"),
		},
		{
			name:           "ErrEncodeState",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: true},
			inState:        "nonempty",
			errDecodeState: nil,
			inStateDecoded: token.State{
				Boards: []token.Board{{Columns: []token.Column{{
					Tasks: []token.Task{{ID: "foo"}},
				}}}},
			},
			errDeleteTask:  nil,
			errEncodeState: errors.New("encode state failed"),
			outState:       "",
			wantStatus:     http.StatusInternalServerError,
			assertFunc:     assert.OnLoggedErr("encode state failed"),
		},
		{
			name:           "Success",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			auth:           token.Auth{IsAdmin: true},
			inState:        "nonempty",
			errDecodeState: nil,
			inStateDecoded: token.State{
				Boards: []token.Board{{Columns: []token.Column{{
					Tasks: []token.Task{{ID: "foo"}},
				}}}},
			},
			errDeleteTask:  nil,
			errEncodeState: nil,
			outState:       "asdfkljhadfskjsdfah",
			wantStatus:     http.StatusOK,
			assertFunc: func(t *testing.T, r *http.Response, _ string) {
				ckState := r.Cookies()[0]
				assert.Equal(t.Error, ckState.Value, "asdfkljhadfskjsdfah")
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Decoded = c.auth
			decodeAuth.Err = c.errDecodeAuth
			decodeState.Decoded = c.inStateDecoded
			decodeState.Err = c.errDecodeState
			taskDeleter.Err = c.errDeleteTask
			encodeState.Encoded = c.outState
			encodeState.Err = c.errEncodeState

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
			res := w.Result()

			assert.Equal(t.Error, res.StatusCode, c.wantStatus)

			c.assertFunc(t, res, log.InMessage)
		})
	}
}
