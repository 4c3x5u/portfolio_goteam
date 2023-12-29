//go:build utest

package taskapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/tasktbl"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/validator"
)

// TestPostHandler tests the Handle method of PostHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPostHandler(t *testing.T) {
	authDecoder := &cookie.FakeDecoder[cookie.Auth]{}
	validate := &validator.FakeFunc[PostReq]{}
	taskInserter := &db.FakeInserter[tasktbl.Task]{}
	log := &log.FakeErrorer{}
	sut := NewPostHandler(
		authDecoder,
		validate.Func,
		taskInserter,
		log,
	)

	for _, c := range []struct {
		name          string
		authToken     string
		authDecoded   cookie.Auth
		errDecodeAuth error
		errValidate   error
		errInsertTask error
		wantStatus    int
		assertFunc    func(*testing.T, *http.Response, []any)
	}{
		{
			name:          "NoAuth",
			authToken:     "",
			errDecodeAuth: cookie.ErrInvalid,
			errValidate:   nil,
			errInsertTask: nil,
			wantStatus:    http.StatusUnauthorized,
			assertFunc:    assert.OnRespErr("Auth token not found."),
		},
		{
			name:          "InvalidAuth",
			authToken:     "nonempty",
			errDecodeAuth: cookie.ErrInvalid,
			errValidate:   nil,
			errInsertTask: nil,
			wantStatus:    http.StatusUnauthorized,
			assertFunc:    assert.OnRespErr("Invalid auth token."),
		},
		{
			name:          "NotAdmin",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{},
			errDecodeAuth: nil,
			errValidate:   nil,
			errInsertTask: nil,
			wantStatus:    http.StatusForbidden,
			assertFunc: assert.OnRespErr(
				"Only team admins can create tasks.",
			),
		},
		{
			name:          "ErrBoardIDEmpty",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			errValidate:   errBoardIDEmpty,
			errInsertTask: nil,
			wantStatus:    http.StatusBadRequest,
			assertFunc:    assert.OnRespErr("Board ID cannot be empty."),
		},
		{
			name:          "ErrParseBoardID",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			errValidate:   errParseBoardID,
			errInsertTask: nil,
			wantStatus:    http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Board ID is must be a valid UUID.",
			),
		},
		{
			name:          "ErrColNoOutOfBounds",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			errValidate:   errColNoOutOfBounds,
			errInsertTask: nil,
			wantStatus:    http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Column number must be between 1 and 4.",
			),
		},
		{
			name:          "ErrTitleEmpty",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			errValidate:   errTitleEmpty,
			errInsertTask: nil,
			wantStatus:    http.StatusBadRequest,
			assertFunc:    assert.OnRespErr("Task title cannot be empty."),
		},
		{
			name:          "ErrTitleTooLong",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			errValidate:   errTitleTooLong,
			errInsertTask: nil,
			wantStatus:    http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:          "ErrDescTooLong",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			errValidate:   errDescTooLong,
			errInsertTask: nil,
			wantStatus:    http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Task description cannot be longer than 500 characters.",
			),
		},
		{
			name:          "ErrSubtaskTitleEmpty",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			errValidate:   errSubtaskTitleEmpty,
			errInsertTask: nil,
			wantStatus:    http.StatusBadRequest,
			assertFunc:    assert.OnRespErr("Subtask title cannot be empty."),
		},
		{
			name:          "ErrSubtaskTitleTooLong",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			errValidate:   errSubtaskTitleTooLong,
			errInsertTask: nil,
			wantStatus:    http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Subtask title cannot be longer than 50 characters.",
			),
		},
		{
			name:          "ErrOrderNegative",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			errValidate:   errOrderNegative,
			errInsertTask: nil,
			wantStatus:    http.StatusBadRequest,
			assertFunc:    assert.OnRespErr("Order cannot be negative."),
		},
		{
			name:          "ErrValidate",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			errValidate:   errors.New("validate failed"),
			errInsertTask: nil,
			wantStatus:    http.StatusInternalServerError,
			assertFunc:    assert.OnLoggedErr("validate failed"),
		},
		{
			name:          "ErrPutTask",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			errValidate:   nil,
			errInsertTask: errors.New("put task failed"),
			wantStatus:    http.StatusInternalServerError,
			assertFunc:    assert.OnLoggedErr("put task failed"),
		},
		{
			name:          "OK",
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			errValidate:   nil,
			errInsertTask: nil,
			wantStatus:    http.StatusOK,
			assertFunc:    func(*testing.T, *http.Response, []any) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			authDecoder.Res = c.authDecoded
			authDecoder.Err = c.errDecodeAuth
			validate.Err = c.errValidate
			taskInserter.Err = c.errInsertTask
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost, "/", strings.NewReader("{}"),
			)
			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name: "auth-token", Value: c.authToken,
				})
			}

			sut.Handle(w, r, "")

			resp := w.Result()
			assert.Equal(t.Error, resp.StatusCode, c.wantStatus)
			c.assertFunc(t, w.Result(), log.Args)
		})
	}
}
