//go:build utest

package task

import (
	"bytes"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/pkg/api"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/tasktable"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
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
	taskInserter := &db.FakeInserter[tasktable.Task]{}
	stateEncoder := &cookie.FakeEncoder[cookie.State]{}
	log := &pkgLog.FakeErrorer{}
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
		reqBody              string
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
		assertFunc           func(*testing.T, *http.Response, string)
	}{
		{
			name:       "NoAuth",
			wantStatus: http.StatusUnauthorized,
			assertFunc: assert.OnResErr("Auth token not found."),
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
			assertFunc:           assert.OnResErr("Invalid auth token."),
		},
		{
			name:                 "NotAdmin",
			reqBody:              "",
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
			assertFunc: assert.OnResErr(
				"Only team admins can create tasks.",
			),
		},
		{
			name:                 "NoState",
			reqBody:              "",
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
			assertFunc:           assert.OnResErr("State token not found."),
		},
		{
			name:                 "InvalidState",
			reqBody:              "",
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
			assertFunc:           assert.OnResErr("Invalid state token."),
		},
		{
			name:          "ColNoOutOfBounds",
			reqBody:       `{"board": "boardid"}`,
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
			assertFunc: assert.OnResErr(
				"Column number out of bounds.",
			),
		},
		{
			name:          "NoBoardAccess",
			reqBody:       `{"board": "boardid"}`,
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
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:          "TaskTitleEmpty",
			reqBody:       `{"board": "boardid", "column": 0}`,
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
			assertFunc: assert.OnResErr(
				"Task title cannot be empty.",
			),
		},
		{
			name:          "TaskTitleTooLong",
			reqBody:       `{"board": "boardid", "column": 0}`,
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
			assertFunc: assert.OnResErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:          "TaskTitleUnexpectedErr",
			reqBody:       `{"board": "boardid", "column": 0}`,
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
			reqBody: `{
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
			assertFunc: assert.OnResErr(
				"Subtask title cannot be empty.",
			),
		},
		{
			name: "SubtTitleTooLong",
			reqBody: `{
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
			assertFunc: assert.OnResErr(
				"Subtask title cannot be longer than 50 characters.",
			),
		},
		{
			name: "ValidateSubtTitleErr",
			reqBody: `{
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
			reqBody:       `{"board": "boardid", "column": 0}`,
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
			reqBody:       `{"board": "boardid", "column": 0}`,
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
			reqBody:       `{"board": "boardid"}`,
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
			assertFunc: func(t *testing.T, resp *http.Response, _ string) {
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

			req := httptest.NewRequest(
				http.MethodPost, "/", bytes.NewReader([]byte(c.reqBody)),
			)
			if c.authToken != "" {
				req.AddCookie(&http.Cookie{
					Name: "auth-token", Value: c.authToken,
				})
			}
			if c.inStateToken != "" {
				req.AddCookie(&http.Cookie{
					Name: "state-token", Value: c.inStateToken,
				})
			}

			w := httptest.NewRecorder()

			sut.Handle(w, req, "")
			res := w.Result()

			assert.Equal(t.Error, res.StatusCode, c.wantStatus)

			c.assertFunc(t, w.Result(), log.InMessage)
		})
	}
}
