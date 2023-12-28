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

// TestDeleteHandler tests the Handle method of DELETEHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDeleteHandler(t *testing.T) {
	authDecoder := &cookie.FakeDecoder[cookie.Auth]{}
	deleter := &db.FakeDeleterDualKey{}
	log := &log.FakeErrorer{}
	sut := NewDeleteHandler(authDecoder, deleter, log)

	for _, c := range []struct {
		name           string
		boardID        string
		authToken      string
		errDecodeAuth  error
		authDecoded    cookie.Auth
		deleteBoardErr error
		wantStatusCode int
		assertFunc     func(*testing.T, *http.Response, []any)
	}{
		{
			name:           "NoAuth",
			boardID:        "",
			authToken:      "",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusUnauthorized,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "InvalidAuth",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  cookie.ErrInvalid,
			authDecoded:    cookie.Auth{},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusUnauthorized,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "NotAdmin",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: false},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "EmptyID",
			boardID:        "",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "InvalidID",
			boardID:        "adksfjahsd",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "ErrNoItem",
			boardID:        "66c16e54-c14f-4481-ada6-404bca897fb0",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true, TeamID: "1"},
			deleteBoardErr: db.ErrNoItem,
			wantStatusCode: http.StatusNotFound,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
		{
			name:           "DeleteErr",
			boardID:        "66c16e54-c14f-4481-ada6-404bca897fb0",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true, TeamID: "1"},
			deleteBoardErr: errors.New("delete board failed"),
			wantStatusCode: http.StatusInternalServerError,
			assertFunc:     assert.OnLoggedErr("delete board failed"),
		},
		{
			name:           "Success",
			boardID:        "66c16e54-c14f-4481-ada6-404bca897fb0",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    cookie.Auth{IsAdmin: true, TeamID: "1"},
			deleteBoardErr: nil,
			wantStatusCode: http.StatusOK,
			assertFunc:     func(*testing.T, *http.Response, []any) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			authDecoder.Err = c.errDecodeAuth
			authDecoder.Res = c.authDecoded
			deleter.Err = c.deleteBoardErr
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodPost, "/?id="+c.boardID, nil)
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
