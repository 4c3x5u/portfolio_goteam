//go:build utest

package task

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
	taskTable "github.com/kxplxn/goteam/pkg/db/task"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// TestPatchHandler tests the PATCH handler.
func TestPatchHandler(t *testing.T) {
	decodeAuth := &token.FakeDecode[token.Auth]{}
	decodeState := &token.FakeDecode[token.State]{}
	titleValidator := &api.FakeStringValidator{}
	subtTitleValidator := &api.FakeStringValidator{}
	taskUpdater := &db.FakeUpdater[taskTable.Task]{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPatchHandler(
		decodeAuth.Func,
		decodeState.Func,
		titleValidator,
		subtTitleValidator,
		taskUpdater,
		log,
	)

	for _, c := range []struct {
		name                 string
		authToken            string
		authDecoded          token.Auth
		errDecodeAuth        error
		stateToken           string
		stateDecoded         token.State
		errDecodeState       error
		errValidateTitle     error
		errValidateSubtTitle error
		taskUpdaterErr       error
		wantStatusCode       int
		assertFunc           func(*testing.T, *http.Response, string)
	}{
		{
			name:                 "NoAuth",
			authToken:            "",
			authDecoded:          token.Auth{},
			errDecodeAuth:        nil,
			stateToken:           "",
			stateDecoded:         token.State{},
			errDecodeState:       nil,
			errValidateTitle:     nil,
			errValidateSubtTitle: nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusUnauthorized,
			assertFunc:           assert.OnResErr("Auth token not found."),
		},
		{
			name:                 "ErrDecodeAuth",
			authToken:            "nonempty",
			authDecoded:          token.Auth{},
			errDecodeAuth:        token.ErrInvalid,
			stateToken:           "",
			stateDecoded:         token.State{},
			errDecodeState:       nil,
			errValidateTitle:     nil,
			errValidateSubtTitle: nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusUnauthorized,
			assertFunc:           assert.OnResErr("Invalid auth token."),
		},
		{
			name:                 "NotAdmin",
			authToken:            "nonempty",
			authDecoded:          token.Auth{IsAdmin: false},
			errDecodeAuth:        nil,
			stateToken:           "",
			stateDecoded:         token.State{},
			errDecodeState:       nil,
			errValidateTitle:     nil,
			errValidateSubtTitle: nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can edit tasks.",
			),
		},
		{
			name:                 "NoState",
			authToken:            "nonempty",
			authDecoded:          token.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			stateToken:           "",
			stateDecoded:         token.State{},
			errDecodeState:       nil,
			errValidateTitle:     nil,
			errValidateSubtTitle: nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc:           assert.OnResErr("State token not found."),
		},
		{
			name:                 "ErrDecodeState",
			authToken:            "nonempty",
			authDecoded:          token.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			stateToken:           "nonempty",
			stateDecoded:         token.State{},
			errDecodeState:       token.ErrInvalid,
			errValidateTitle:     nil,
			errValidateSubtTitle: nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc:           assert.OnResErr("Invalid state token."),
		},
		// task id is invalid when it is not found in state
		{
			name:                 "TaskIDInvalid",
			authToken:            "nonempty",
			authDecoded:          token.Auth{IsAdmin: true},
			errDecodeAuth:        nil,
			stateToken:           "nonempty",
			stateDecoded:         token.State{},
			errDecodeState:       nil,
			errValidateTitle:     nil,
			errValidateSubtTitle: nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc:           assert.OnResErr("Invalid task ID."),
		},
		{
			name:          "TaskTitleEmpty",
			authToken:     "nonempty",
			authDecoded:   token.Auth{IsAdmin: true, TeamID: "21"},
			errDecodeAuth: nil,
			stateToken:    "nonempty",
			stateDecoded: token.State{Boards: []token.Board{{
				Columns: []token.Column{{Tasks: []token.Task{{ID: "qwerty"}}}}},
			}},
			errDecodeState:       nil,
			errValidateTitle:     api.ErrEmpty,
			errValidateSubtTitle: nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be empty.",
			),
		},
		{
			name:          "TaskTitleTooLong",
			authToken:     "nonempty",
			authDecoded:   token.Auth{IsAdmin: true, TeamID: "21"},
			errDecodeAuth: nil,
			stateToken:    "nonempty",
			stateDecoded: token.State{Boards: []token.Board{{
				Columns: []token.Column{{Tasks: []token.Task{{ID: "qwerty"}}}}},
			}},
			errDecodeState:       nil,
			errValidateTitle:     api.ErrTooLong,
			errValidateSubtTitle: nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:          "TaskTitleErr",
			authToken:     "nonempty",
			authDecoded:   token.Auth{IsAdmin: true, TeamID: "21"},
			errDecodeAuth: nil,
			stateToken:    "nonempty",
			stateDecoded: token.State{Boards: []token.Board{{
				Columns: []token.Column{{Tasks: []token.Task{{ID: "qwerty"}}}}},
			}},
			errDecodeState:       nil,
			errValidateTitle:     api.ErrNotInt,
			errValidateSubtTitle: nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				api.ErrNotInt.Error(),
			),
		},
		{
			name:          "SubtaskTitleEmpty",
			authToken:     "nonempty",
			authDecoded:   token.Auth{IsAdmin: true, TeamID: "21"},
			errDecodeAuth: nil,
			stateToken:    "nonempty",
			stateDecoded: token.State{Boards: []token.Board{{
				Columns: []token.Column{{Tasks: []token.Task{{ID: "qwerty"}}}}},
			}},
			errDecodeState:       nil,
			errValidateTitle:     nil,
			errValidateSubtTitle: api.ErrEmpty,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be empty.",
			),
		},
		{
			name:          "SubtaskTitleTooLong",
			authToken:     "nonempty",
			authDecoded:   token.Auth{IsAdmin: true, TeamID: "21"},
			errDecodeAuth: nil,
			stateToken:    "nonempty",
			stateDecoded: token.State{Boards: []token.Board{{
				Columns: []token.Column{{Tasks: []token.Task{{ID: "qwerty"}}}}},
			}},
			errDecodeState:       nil,
			errValidateTitle:     nil,
			errValidateSubtTitle: api.ErrTooLong,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be longer than 50 characters.",
			),
		},
		{
			name:          "SubtaskTitleErr",
			authToken:     "nonempty",
			authDecoded:   token.Auth{IsAdmin: true, TeamID: "21"},
			errDecodeAuth: nil,
			stateToken:    "nonempty",
			stateDecoded: token.State{Boards: []token.Board{{
				Columns: []token.Column{{Tasks: []token.Task{{ID: "qwerty"}}}}},
			}},
			errDecodeState:       nil,
			errValidateTitle:     nil,
			errValidateSubtTitle: api.ErrNotInt,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				api.ErrNotInt.Error(),
			),
		},
		{
			name:          "TaskNotFound",
			authToken:     "nonempty",
			authDecoded:   token.Auth{IsAdmin: true, TeamID: "21"},
			errDecodeAuth: nil,
			stateToken:    "nonempty",
			stateDecoded: token.State{Boards: []token.Board{{
				Columns: []token.Column{{Tasks: []token.Task{{ID: "qwerty"}}}}},
			}},
			errDecodeState:       nil,
			errValidateTitle:     nil,
			errValidateSubtTitle: nil,
			taskUpdaterErr:       db.ErrNoItem,
			wantStatusCode:       http.StatusNotFound,
			assertFunc:           assert.OnResErr("Task not found."),
		},
		{
			name:          "TaskUpdaterErr",
			authToken:     "nonempty",
			authDecoded:   token.Auth{IsAdmin: true, TeamID: "21"},
			errDecodeAuth: nil,
			stateToken:    "nonempty",
			stateDecoded: token.State{Boards: []token.Board{{
				Columns: []token.Column{{Tasks: []token.Task{{ID: "qwerty"}}}}},
			}},
			errDecodeState:       nil,
			errValidateTitle:     nil,
			errValidateSubtTitle: nil,
			taskUpdaterErr:       errors.New("update task failed"),
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr("update task failed"),
		},
		{
			name:          "Success",
			authToken:     "nonempty",
			authDecoded:   token.Auth{IsAdmin: true, TeamID: "21"},
			errDecodeAuth: nil,
			stateToken:    "nonempty",
			stateDecoded: token.State{Boards: []token.Board{{
				Columns: []token.Column{{Tasks: []token.Task{{ID: "qwerty"}}}}},
			}},
			errDecodeState:       nil,
			errValidateTitle:     nil,
			errValidateSubtTitle: nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusOK,
			assertFunc:           func(*testing.T, *http.Response, string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Res = c.authDecoded
			decodeAuth.Err = c.errDecodeAuth
			decodeState.Res = c.stateDecoded
			decodeState.Err = c.errDecodeState
			titleValidator.Err = c.errValidateTitle
			subtTitleValidator.Err = c.errValidateSubtTitle
			taskUpdater.Err = c.taskUpdaterErr
			r := httptest.NewRequest("", "/?id=qwerty", strings.NewReader(`{
				"column":      0,
				"title":       "",
				"description": "",
				"subtasks":    [{"title": ""}]
			}`))
			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  "auth-token",
					Value: c.authToken,
				})
			}
			if c.stateToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  "state-token",
					Value: c.stateToken,
				})
			}

			w := httptest.NewRecorder()

			sut.Handle(w, r, "")

			res := w.Result()
			assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
