//go:build utest

package task

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/assert"
	boardTable "github.com/kxplxn/goteam/pkg/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/pkg/dbaccess/column"
	taskTable "github.com/kxplxn/goteam/pkg/dbaccess/task"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// TestPatchHandler tests the PATCH handler.
func TestPatchHandler(t *testing.T) {
	decodeAuth := &token.FakeDecode[token.Auth]{}
	decodeState := &token.FakeDecode[token.State]{}
	titleValidator := &api.FakeStringValidator{}
	subtTitleValidator := &api.FakeStringValidator{}
	taskSelector := &taskTable.FakeSelector{}
	columnSelector := &columnTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	taskUpdater := &taskTable.FakeUpdater{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPatchHandler(
		decodeAuth.Func,
		decodeState.Func,
		titleValidator,
		subtTitleValidator,
		taskSelector,
		columnSelector,
		boardSelector,
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
		errSelectTask        error
		errSelectColumn      error
		board                boardTable.Record
		boardSelectorErr     error
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{},
			boardSelectorErr:     nil,
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{},
			boardSelectorErr:     nil,
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{},
			boardSelectorErr:     nil,
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{},
			boardSelectorErr:     nil,
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{},
			boardSelectorErr:     nil,
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{},
			boardSelectorErr:     nil,
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     nil,
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Task title cannot be longer than 50 characters.",
			),
		},
		{
			name:          "TaskTitleUnexpectedErr",
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     nil,
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     nil,
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Subtask title cannot be longer than 50 characters.",
			),
		},
		{
			name:          "SubtaskTitleUnexpectedErr",
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     nil,
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
			errSelectTask:        sql.ErrNoRows,
			errSelectColumn:      nil,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusNotFound,
			assertFunc:           assert.OnResErr("Task not found."),
		},
		{
			name:          "TaskSelectorErr",
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
			errSelectTask:        sql.ErrConnDone,
			errSelectColumn:      nil,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:          "ColumnSelectorErr",
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
			errSelectTask:        nil,
			errSelectColumn:      sql.ErrNoRows,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:          "BoardSelectorErr",
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     sql.ErrNoRows,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:          "NoAccess",
			authToken:     "nonempty",
			authDecoded:   token.Auth{IsAdmin: true, TeamID: "32"},
			errDecodeAuth: nil,
			stateToken:    "nonempty",
			stateDecoded: token.State{Boards: []token.Board{{
				Columns: []token.Column{{Tasks: []token.Task{{ID: "qwerty"}}}}},
			}},
			errDecodeState:       nil,
			errValidateTitle:     nil,
			errValidateSubtTitle: nil,
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
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
			errDecodeState: nil,
			// user: userTable.Record{
			// 	IsAdmin: true, TeamID: 1,
			// },
			errValidateTitle:     nil,
			errValidateSubtTitle: nil,
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     nil,
			taskUpdaterErr:       sql.ErrConnDone,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
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
			errSelectTask:        nil,
			errSelectColumn:      nil,
			board:                boardTable.Record{TeamID: 21},
			boardSelectorErr:     nil,
			taskUpdaterErr:       nil,
			wantStatusCode:       http.StatusOK,
			assertFunc:           func(*testing.T, *http.Response, string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Decoded = c.authDecoded
			decodeAuth.Err = c.errDecodeAuth
			decodeState.Decoded = c.stateDecoded
			decodeState.Err = c.errDecodeState
			titleValidator.Err = c.errValidateTitle
			subtTitleValidator.Err = c.errValidateSubtTitle
			taskSelector.Err = c.errSelectTask
			columnSelector.Err = c.errSelectColumn
			boardSelector.Board = c.board
			boardSelector.Err = c.boardSelectorErr
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
