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
	stateDecoder := &cookie.FakeDecoder[cookie.State]{}
	titleValidator := &api.FakeStringValidator{}
	subtTitleValidator := &api.FakeStringValidator{}
	colNoValidator := &api.FakeIntValidator{}
	taskInserter := &db.FakeInserter[tasktbl.Task]{}
	stateEncoder := &cookie.FakeEncoder[cookie.State]{}
	log := &log.FakeErrorer{}
	sut := NewPostHandler(
		authDecoder,
		stateDecoder,
		titleValidator,
		subtTitleValidator,
		colNoValidator,
		taskInserter,
		stateEncoder,
		log,
	)

	for _, c := range []struct {
		name                 string
		rBody                string
		authToken            string
		authDecoded          cookie.Auth
		errDecodeAuth        error
		inStateToken         string
		inStateDecoded       cookie.State
		errDecodeInState     error
		errValidateColNo     error
		errValidateTaskTitle error
		errValidateSubtTitle error
		errInsertTask        error
		outState             http.Cookie
		errEncodeState       error
		wantStatus           int
		assertFunc           func(*testing.T, *http.Response, []any)
	}{
		{
			name:       "NoAuth",
			wantStatus: http.StatusUnauthorized,
			assertFunc: assert.OnRespErr("Auth token not found."),
		},
		{
			name:                 "InvalidAuth",
			authToken:            "nonempty",
			errDecodeAuth:        cookie.ErrInvalid,
			inStateToken:         "",
			inStateDecoded:       cookie.State{},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       nil,
			wantStatus:           http.StatusUnauthorized,
			assertFunc:           assert.OnRespErr("Invalid auth token."),
		},
		{
			name:                 "NotAdmin",
			rBody:                "",
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{},
			errDecodeAuth:        nil,
			inStateToken:         "",
			inStateDecoded:       cookie.State{},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       nil,
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
			inStateToken:         "",
			inStateDecoded:       cookie.State{},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       nil,
			wantStatus:           http.StatusBadRequest,
			assertFunc:           assert.OnRespErr("State token not found."),
		},
		{
			name:                 "InvalidState",
			rBody:                "",
			authToken:            "nonempty",
			authDecoded:          cookie.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			inStateToken:         "nonempty",
			inStateDecoded:       cookie.State{},
			errDecodeInState:     cookie.ErrInvalid,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       nil,
			wantStatus:           http.StatusBadRequest,
			assertFunc:           assert.OnRespErr("Invalid state token."),
		},
		{
			name:          "ColNoOutOfBounds",
			rBody:         `{"board": "boardid"}`,
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			inStateToken:  "nonempty",
			inStateDecoded: cookie.State{
				Boards: []cookie.Board{{ID: "boardid"}},
			},
			errDecodeInState:     nil,
			errValidateColNo:     validator.ErrOutOfBounds,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       nil,
			wantStatus:           http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Column number out of bounds.",
			),
		},
		{
			name:          "NoBoardAccess",
			rBody:         `{"board": "boardid"}`,
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			inStateToken:  "nonempty",
			inStateDecoded: cookie.State{
				Boards: []cookie.Board{{ID: "foo"}},
			},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       nil,
			wantStatus:           http.StatusForbidden,
			assertFunc: assert.OnRespErr(
				"You do not have access to this board.",
			),
		},
		{
			name:          "TaskTitleEmpty",
			rBody:         `{"board": "boardid", "column": 0}`,
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			inStateToken:  "nonempty",
			inStateDecoded: cookie.State{Boards: []cookie.Board{{
				ID: "boardid",
				Columns: []cookie.Column{{Tasks: []cookie.Task{{
					ID: "taskid", Order: 0,
				}}}},
			}}},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: validator.ErrEmpty,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       nil,
			wantStatus:           http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Task title cannot be empty.",
			),
		},
		{
			name:          "TaskTitleTooLong",
			rBody:         `{"board": "boardid", "column": 0}`,
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			inStateToken:  "nonempty",
			inStateDecoded: cookie.State{Boards: []cookie.Board{{
				ID: "boardid",
				Columns: []cookie.Column{{Tasks: []cookie.Task{{
					ID: "taskid", Order: 0,
				}}}},
			}}},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: validator.ErrTooLong,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       nil,
			wantStatus:           http.StatusBadRequest,
			assertFunc: assert.OnRespErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:          "TaskTitleUnexpectedErr",
			rBody:         `{"board": "boardid", "column": 0}`,
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			inStateToken:  "nonempty",
			inStateDecoded: cookie.State{Boards: []cookie.Board{{
				ID: "boardid",
				Columns: []cookie.Column{{Tasks: []cookie.Task{{
					ID: "taskid", Order: 0,
				}}}},
			}}},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: validator.ErrWrongFormat,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       nil,
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
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			inStateToken:  "nonempty",
			inStateDecoded: cookie.State{Boards: []cookie.Board{{
				ID: "boardid",
				Columns: []cookie.Column{{Tasks: []cookie.Task{{
					ID: "taskid", Order: 0,
				}}}},
			}}},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: validator.ErrEmpty,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       nil,
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
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			inStateToken:  "nonempty",
			inStateDecoded: cookie.State{Boards: []cookie.Board{{
				ID: "boardid",
				Columns: []cookie.Column{{Tasks: []cookie.Task{{
					ID: "taskid", Order: 0,
				}}}},
			}}},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: validator.ErrTooLong,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       nil,
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
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			inStateToken:  "nonempty",
			inStateDecoded: cookie.State{Boards: []cookie.Board{{
				ID: "boardid",
				Columns: []cookie.Column{{Tasks: []cookie.Task{{
					ID: "taskid", Order: 0,
				}}}},
			}}},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: validator.ErrWrongFormat,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       nil,
			wantStatus:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				validator.ErrWrongFormat.Error(),
			),
		},
		{
			name:          "ErrPutTask",
			rBody:         `{"board": "boardid", "column": 0}`,
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			inStateToken:  "nonempty",
			inStateDecoded: cookie.State{Boards: []cookie.Board{{
				ID: "boardid",
				Columns: []cookie.Column{{Tasks: []cookie.Task{{
					ID: "taskid", Order: 0,
				}}}},
			}}},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        errors.New("failed to put task"),
			outState:             http.Cookie{},
			errEncodeState:       nil,
			wantStatus:           http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				"failed to put task",
			),
		},
		{
			name:          "ErrEncodeState",
			rBody:         `{"board": "boardid", "column": 0}`,
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			inStateToken:  "nonempty",
			inStateDecoded: cookie.State{Boards: []cookie.Board{{
				ID: "boardid",
				Columns: []cookie.Column{{Tasks: []cookie.Task{{
					ID: "taskid", Order: 0,
				}}}},
			}}},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			outState:             http.Cookie{},
			errEncodeState:       errors.New("encode state failed"),
			wantStatus:           http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr("encode state failed"),
		},
		{
			name:          "OK",
			rBody:         `{"board": "boardid"}`,
			authToken:     "nonempty",
			authDecoded:   cookie.Auth{IsAdmin: true},
			errDecodeAuth: nil,
			inStateToken:  "nonempty",
			inStateDecoded: cookie.State{Boards: []cookie.Board{{
				ID: "boardid",
				Columns: []cookie.Column{{Tasks: []cookie.Task{{
					ID: "taskid", Order: 0,
				}}}},
			}}},
			errDecodeInState:     nil,
			errValidateColNo:     nil,
			errValidateTaskTitle: nil,
			errValidateSubtTitle: nil,
			errInsertTask:        nil,
			outState:             http.Cookie{Name: "foo", Value: "bar"},
			errEncodeState:       nil,
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
			stateDecoder.Res = c.inStateDecoded
			stateDecoder.Err = c.errDecodeInState
			colNoValidator.Err = c.errValidateColNo
			titleValidator.Err = c.errValidateTaskTitle
			subtTitleValidator.Err = c.errValidateSubtTitle
			taskInserter.Err = c.errInsertTask
			stateEncoder.Res = c.outState
			stateEncoder.Err = c.errEncodeState
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost, "/", strings.NewReader(c.rBody),
			)
			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name: "auth-token", Value: c.authToken,
				})
			}
			if c.inStateToken != "" {
				r.AddCookie(&http.Cookie{
					Name: "state-token", Value: c.inStateToken,
				})
			}

			sut.Handle(w, r, "")

			resp := w.Result()
			assert.Equal(t.Error, resp.StatusCode, c.wantStatus)
			c.assertFunc(t, w.Result(), log.Args)
		})
	}
}
