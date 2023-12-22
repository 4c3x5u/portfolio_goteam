//go:build utest

package board

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

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDELETEHandler(t *testing.T) {
	decodeAuth := &token.FakeDecode[token.Auth]{}
	decodeState := &token.FakeDecode[token.State]{}
	deleter := &db.FakeDeleterDualKey{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDeleteHandler(decodeAuth.Func, decodeState.Func, deleter, log)

	for _, c := range []struct {
		name           string
		boardID        string
		authToken      string
		errDecodeAuth  error
		authDecoded    token.Auth
		stateToken     string
		errDecodeState error
		stateDecoded   token.State
		deleteBoardErr error
		wantStatusCode int
		assertFunc     func(*testing.T, *http.Response, string)
	}{
		{
			name:           "NoAuth",
			boardID:        "",
			authToken:      "",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{},
			stateToken:     "",
			errDecodeState: nil,
			stateDecoded:   token.State{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusUnauthorized,
			assertFunc:     func(*testing.T, *http.Response, string) {},
		},
		{
			name:           "InvalidAuth",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  token.ErrInvalid,
			authDecoded:    token.Auth{},
			stateToken:     "",
			errDecodeState: nil,
			stateDecoded:   token.State{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusUnauthorized,
			assertFunc:     func(*testing.T, *http.Response, string) {},
		},
		{
			name:           "NotAdmin",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: false},
			stateToken:     "",
			errDecodeState: nil,
			stateDecoded:   token.State{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     func(*testing.T, *http.Response, string) {},
		},
		{
			name:           "NoState",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "",
			errDecodeState: nil,
			stateDecoded:   token.State{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     func(*testing.T, *http.Response, string) {},
		},
		{
			name:           "InvalidState",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: token.ErrInvalid,
			stateDecoded:   token.State{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     func(*testing.T, *http.Response, string) {},
		},
		{
			name:           "EmptyID",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded:   token.State{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     func(*testing.T, *http.Response, string) {},
		},
		{
			name:           "InvalidID",
			boardID:        "adksfjahsd",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded:   token.State{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     func(*testing.T, *http.Response, string) {},
		},
		{
			name:           "NoAccess",
			boardID:        "66c16e54-c14f-4481-ada6-404bca897fb0",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded:   token.State{Boards: []token.Board{{ID: "adsjkhf"}}},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     func(*testing.T, *http.Response, string) {},
		},
		{
			name:           "DeleteErr",
			boardID:        "66c16e54-c14f-4481-ada6-404bca897fb0",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true, TeamID: "1"},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: token.State{Boards: []token.Board{
				{ID: "66c16e54-c14f-4481-ada6-404bca897fb0"},
			}},
			deleteBoardErr: errors.New("delete board failed"),
			wantStatusCode: http.StatusInternalServerError,
			assertFunc:     assert.OnLoggedErr("delete board failed"),
		},
		{
			name:           "Success",
			boardID:        "66c16e54-c14f-4481-ada6-404bca897fb0",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true, TeamID: "1"},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: token.State{Boards: []token.Board{
				{ID: "66c16e54-c14f-4481-ada6-404bca897fb0"},
			}},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusOK,
			assertFunc:     func(*testing.T, *http.Response, string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for sut's dependencies.
			decodeAuth.Err = c.errDecodeAuth
			decodeAuth.Res = c.authDecoded
			decodeState.Err = c.errDecodeState
			decodeState.Res = c.stateDecoded
			deleter.Err = c.deleteBoardErr

			// Prepare request and response recorder.
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
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.Handle(w, r, "")
			res := w.Result()

			// Assert on the status code.
			assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

			// Run case-specific assertions.
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
