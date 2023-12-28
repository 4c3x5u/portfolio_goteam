//go:build utest

package taskapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/pkg/api"
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
	titleValidator := &api.FakeStringValidator{}
	subtTitleValidator := &api.FakeStringValidator{}
	colNoValidator := &api.FakeIntValidator{}
	taskInserter := &db.FakeInserter[tasktbl.Task]{}
	log := &log.FakeErrorer{}
	sut := NewPostHandler(
		authDecoder,
		titleValidator,
		subtTitleValidator,
		colNoValidator,
		taskInserter,
		log,
	)

	for _, c := range []struct {
		name                 string
		rBody                string
		authToken            string
		authDecoded          cookie.Auth
		errDecodeAuth        error
		errValidateColNo     error
		errValidateTaskTitle error
		errValidateSubtTitle error
		errInsertTask        error
		wantStatus           int
		assertFunc           func(*testing.T, *http.Response, []any)
	}{
		{
			name:                 "NoAuth",
			authToken:            "",
			errDecodeAuth:        cookie.ErrInvalid,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			wantStatus:           http.StatusUnauthorized,
			assertFunc:           assert.OnRespErr("Auth token not found."),
		},
		{
			name:                 "InvalidAuth",
			authToken:            "nonempty",
			errDecodeAuth:        cookie.ErrInvalid,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			wantStatus:           http.StatusUnauthorized,
			assertFunc:           assert.OnRespErr("Invalid auth token."),
		},
		{
			name:                 "NotAdmin",
			rBody:                "",
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			wantStatus:           http.StatusForbidden,
			assertFunc: assert.OnRespErr(
				"Only team admins can create tasks.",
			),
		},
		{
			name:                 "NoState",
			rBody:                "",
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			wantStatus:           http.StatusBadRequest,
			assertFunc:           assert.OnRespErr("State token not found."),
		},
		{
			name:                 "InvalidState",
			rBody:                "",
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			wantStatus:           http.StatusBadRequest,
			assertFunc:           assert.OnRespErr("Invalid state token."),
		},
		{
			name:                 "ColNoOutOfBounds",
			rBody:                `{"board": "boardid"}`,
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     validator.ErrOutOfBounds,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			wantStatus:           http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Column number out of bounds.",
			),
		},
		{
			name:                 "NoBoardAccess",
			rBody:                `{"board": "boardid"}`,
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			wantStatus:           http.StatusForbidden,
			assertFunc: assert.OnRespErr(
				"You do not have access to this board.",
			),
		},
		{
			name:                 "TaskTitleEmpty",
			rBody:                `{"board": "boardid", "column": 0}`,
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: validator.ErrEmpty,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			wantStatus:           http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Task title cannot be empty.",
			),
		},
		{
			name:                 "TaskTitleTooLong",
			rBody:                `{"board": "boardid", "column": 0}`,
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: validator.ErrTooLong,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			wantStatus:           http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:                 "TaskTitleUnexpectedErr",
			rBody:                `{"board": "boardid", "column": 0}`,
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: validator.ErrWrongFormat,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			wantStatus:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				validator.ErrWrongFormat.Error(),
			),
		},
		{
			name: "SubtTitleEmpty",
			rBody: `{
                "board": "boardid", "column": 0, "subtasks": ["foo"]
            }`,
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: validator.ErrEmpty,
			errInsertTask:        nil,
			wantStatus:           http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Subtask title cannot be empty.",
			),
		},
		{
			name: "SubtTitleTooLong",
			rBody: `{
                "board": "boardid", "column": 0, "subtasks": ["foo"]
            }`,
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: validator.ErrTooLong,
			errInsertTask:        nil,
			wantStatus:           http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Subtask title cannot be longer than 50 characters.",
			),
		},
		{
			name: "ValidateSubtTitleErr",
			rBody: `{
                "board": "boardid", "column": 0, "subtasks": ["foo"]
            }`,
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: validator.ErrWrongFormat,
			errInsertTask:        nil,
			wantStatus:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				validator.ErrWrongFormat.Error(),
			),
		},
		{
			name:                 "ErrPutTask",
			rBody:                `{"board": "boardid", "column": 0}`,
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        errors.New("failed to put task"),
			wantStatus:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				"failed to put task",
			),
		},
		{
			name:                 "ErrEncodeState",
			rBody:                `{"board": "boardid", "column": 0}`,
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			wantStatus:           http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr("encode state failed"),
		},
		{
			name:                 "OK",
			rBody:                `{"board": "boardid"}`,
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			wantStatus:           http.StatusOK,
			assertFunc: func(t *testing.T, resp *http.Response, _ []any) {
				ckState := resp.Cookies()[0]
				assert.Equal(t.Error, ckState.Name, "foo")
				assert.Equal(t.Error, ckState.Value, "bar")
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			authDecoder.Res = c.authDecoded
			authDecoder.Err = c.errDecodeAuth
			colNoValidator.Err = c.errValidateColNo
			titleValidator.Err = c.errValidateTaskTitle
			subtTitleValidator.Err = c.errValidateSubtTitle
			taskInserter.Err = c.errInsertTask
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost, "/", strings.NewReader(c.rBody),
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
