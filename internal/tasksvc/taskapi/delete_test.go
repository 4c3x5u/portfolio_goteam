//go:build utest

package taskapi

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
	taskDeleter := &db.FakeDeleterDualKey{}
	log := &log.FakeErrorer{}
	sut := NewDeleteHandler(authDecoder, taskDeleter, log)

	for _, c := range []struct {
		name          string
		authToken     string
		errDecodeAuth error
		auth          cookie.Auth
		errDeleteTask error
		wantStatus    int
		assertFunc    func(*testing.T, *http.Response, []any)
	}{
		{
			name:          "NoAuth",
			authToken:     "",
			errDecodeAuth: nil,
			auth:          cookie.Auth{},
			errDeleteTask: nil,
			wantStatus:    http.StatusUnauthorized,
			assertFunc:    assert.OnRespErr("Auth token not found."),
		},
		{
			name:          "ErrDecodeAuth",
			authToken:     "nonempty",
			errDecodeAuth: errors.New("decode auth failed"),
			auth:          cookie.Auth{},
			errDeleteTask: nil,
			wantStatus:    http.StatusUnauthorized,
			assertFunc:    assert.OnRespErr("Invalid auth token."),
		},
		{
			name:          "NotAdmin",
			authToken:     "nonempty",
			errDecodeAuth: nil,
			auth:          cookie.Auth{IsAdmin: false},
			errDeleteTask: nil,
			wantStatus:    http.StatusForbidden,
			assertFunc: assert.OnRespErr(
				"Only team admins can delete tasks.",
			),
		},
		{
			name:          "NotFound",
			authToken:     "nonempty",
			errDecodeAuth: nil,
			auth:          cookie.Auth{IsAdmin: true},
			errDeleteTask: db.ErrNoItem,
			wantStatus:    http.StatusNotFound,
			assertFunc:    assert.OnRespErr("Task not found."),
		},
		{
			name:          "ErrDeleteTask",
			authToken:     "nonempty",
			errDecodeAuth: nil,
			auth:          cookie.Auth{IsAdmin: true},
			errDeleteTask: errors.New("delete task failed"),
			wantStatus:    http.StatusInternalServerError,
			assertFunc:    assert.OnLoggedErr("delete task failed"),
		},
		{
			name:          "Success",
			authToken:     "nonempty",
			errDecodeAuth: nil,
			auth:          cookie.Auth{IsAdmin: true},
			errDeleteTask: nil,
			wantStatus:    http.StatusOK,
			assertFunc:    func(*testing.T, *http.Response, []any) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			authDecoder.Res = c.auth
			authDecoder.Err = c.errDecodeAuth
			taskDeleter.Err = c.errDeleteTask

			r := httptest.NewRequest("", "/?id=foo", nil)
			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  "auth-token",
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
