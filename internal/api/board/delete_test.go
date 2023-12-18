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
	deleter := &db.FakeDeleter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDELETEHandler(decodeAuth.Func, decodeState.Func, deleter, log)

	// Used on cases where no case-specific assertions are required.
	emptyAssertFunc := func(*testing.T, *http.Response, string) {}

	for _, c := range []struct {
		name           string
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
			authToken:      "",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{},
			stateToken:     "",
			errDecodeState: nil,
			stateDecoded:   token.State{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusUnauthorized,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "InvalidAuth",
			authToken:      "nonempty",
			errDecodeAuth:  token.ErrInvalid,
			authDecoded:    token.Auth{},
			stateToken:     "",
			errDecodeState: nil,
			stateDecoded:   token.State{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusUnauthorized,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "NotAdmin",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: false},
			stateToken:     "",
			errDecodeState: nil,
			stateDecoded:   token.State{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "NoState",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "",
			errDecodeState: nil,
			stateDecoded:   token.State{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "InvalidState",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: token.ErrInvalid,
			stateDecoded:   token.State{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "NoAccess",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded:   token.State{Boards: []token.Board{{ID: "3"}}},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "DeleteErr",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true, TeamID: "1"},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded:   token.State{Boards: []token.Board{{ID: "2"}}},
			deleteBoardErr: errors.New("delete board failed"),
			wantStatusCode: http.StatusInternalServerError,
			assertFunc:     assert.OnLoggedErr("delete board failed"),
		},
		{
			name:           "Success",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true, TeamID: "1"},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded:   token.State{Boards: []token.Board{{ID: "2"}}},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusOK,
			assertFunc:     emptyAssertFunc,
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
			r := httptest.NewRequest(http.MethodPost, "/?id=2", nil)
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
