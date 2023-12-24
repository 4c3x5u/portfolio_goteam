//go:build utest

package boardapi

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

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDELETEHandler(t *testing.T) {
	authDecoder := &cookie.FakeDecoder[cookie.Auth]{}
	stateDecoder := &cookie.FakeDecoder[cookie.State]{}
	deleter := &db.FakeDeleterDualKey{}
	stateEncoder := &cookie.FakeEncoder[cookie.State]{}
	log := &log.FakeErrorer{}
	sut := NewDeleteHandler(
		authDecoder, stateDecoder, deleter, stateEncoder, log,
	)

	for _, c := range []struct {
		name           string
		boardID        string
		authToken      string
		errDecodeAuth  error
		authDecoded    cookie.Auth
		stateToken     string
		errDecodeState error
		stateDecoded   cookie.State
		deleteBoardErr error
		errEncodeState error
		outState       http.Cookie
		wantStatusCode int
		assertFunc     func(*testing.T, *http.Response, []any)
	}{
		{
			name:           "NoAuth",
			boardID:        "",
			authToken:      "",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{},
			stateToken:     "",
			errDecodeState: nil,
			stateDecoded:   cookie.State{},
			deleteBoardErr: nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatusCode: http.StatusUnauthorized,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "InvalidAuth",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  cookie.ErrInvalid,
			authDecoded:    cookie.Auth{},
			stateToken:     "",
			errDecodeState: nil,
			stateDecoded:   cookie.State{},
			deleteBoardErr: nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatusCode: http.StatusUnauthorized,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "NotAdmin",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: false},
			stateToken:     "",
			errDecodeState: nil,
			stateDecoded:   cookie.State{},
			deleteBoardErr: nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatusCode: http.StatusForbidden,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "NoState",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			stateToken:     "",
			errDecodeState: nil,
			stateDecoded:   cookie.State{},
			deleteBoardErr: nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatusCode: http.StatusForbidden,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "InvalidState",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: cookie.ErrInvalid,
			stateDecoded:   cookie.State{},
			deleteBoardErr: nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatusCode: http.StatusForbidden,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "EmptyID",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded:   cookie.State{},
			deleteBoardErr: nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "InvalidID",
			boardID:        "adksfjahsd",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded:   cookie.State{},
			deleteBoardErr: nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "NoAccess",
			boardID:        "66c16e54-c14f-4481-ada6-404bca897fb0",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded:   cookie.State{Boards: []cookie.Board{{ID: "adsjkhf"}}},
			deleteBoardErr: nil,
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatusCode: http.StatusForbidden,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "DeleteErr",
			boardID:        "66c16e54-c14f-4481-ada6-404bca897fb0",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true, TeamID: "1"},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: cookie.State{Boards: []cookie.Board{
				{ID: "66c16e54-c14f-4481-ada6-404bca897fb0"},
			}},
			deleteBoardErr: errors.New("delete board failed"),
			errEncodeState: nil,
			outState:       http.Cookie{},
			wantStatusCode: http.StatusInternalServerError,
			assertFunc:     assert.OnLoggedErr("delete board failed"),
		},
		{
			name:           "ErrEncodeState",
			boardID:        "66c16e54-c14f-4481-ada6-404bca897fb0",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true, TeamID: "1"},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: cookie.State{Boards: []cookie.Board{
				{ID: "66c16e54-c14f-4481-ada6-404bca897fb0"},
			}},
			deleteBoardErr: nil,
			errEncodeState: errors.New("encode state failed"),
			outState:       http.Cookie{},
			wantStatusCode: http.StatusInternalServerError,
			assertFunc:     assert.OnLoggedErr("encode state failed"),
		},
		{
			name:           "Success",
			boardID:        "66c16e54-c14f-4481-ada6-404bca897fb0",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true, TeamID: "1"},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: cookie.State{Boards: []cookie.Board{
				{ID: "66c16e54-c14f-4481-ada6-404bca897fb0"},
			}},
			deleteBoardErr: nil,
			errEncodeState: nil,
			outState:       http.Cookie{Name: "foo", Value: "bar"},
			wantStatusCode: http.StatusOK,
			assertFunc: func(t *testing.T, resp *http.Response, _ []any) {
				ckState := resp.Cookies()[0]
				assert.Equal(t.Error, ckState.Name, "foo")
				assert.Equal(t.Error, ckState.Value, "bar")
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			authDecoder.Err = c.errDecodeAuth
			authDecoder.Res = c.authDecoded
			stateDecoder.Err = c.errDecodeState
			stateDecoder.Res = c.stateDecoded
			deleter.Err = c.deleteBoardErr
			stateEncoder.Err = c.errEncodeState
			stateEncoder.Res = c.outState
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/?id="+c.boardID, nil)
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

			sut.Handle(w, r, "")

			resp := w.Result()
			assert.Equal(t.Error, resp.StatusCode, c.wantStatusCode)
			c.assertFunc(t, resp, log.Args)
		})
	}
}
