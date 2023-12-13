//go:build utest

package tasks

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

func TestPatchHandler(t *testing.T) {
	decodeAuth := token.FakeDecode[token.Auth]{}
	decodeState := token.FakeDecode[token.State]{}
	colNoVdtor := &api.FakeIntValidator{}
	tasksUpdater := &db.FakeUpdater[[]taskTable.Task]{}
	encodeState := token.FakeEncode[token.State]{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPatchHandler(
		decodeAuth.Func,
		decodeState.Func,
		colNoVdtor,
		tasksUpdater,
		encodeState.Func,
		log,
	)

	for _, c := range []struct {
		name             string
		authToken        string
		errDecodeAuth    error
		authDecoded      token.Auth
		stateToken       string
		errDecodeState   error
		stateDecoded     token.State
		errValidateColNo error
		errUpdateTasks   error
		errEncodeState   error
		outState         string
		wantStatus       int
		assertFunc       func(*testing.T, *http.Response, string)
	}{
		{
			name:             "NoAuth",
			authToken:        "",
			errDecodeAuth:    nil,
			authDecoded:      token.Auth{},
			stateToken:       "",
			errDecodeState:   nil,
			stateDecoded:     token.State{},
			errValidateColNo: nil,
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         "",
			wantStatus:       http.StatusUnauthorized,
			assertFunc:       assert.OnResErr("Auth token not found."),
		},
		{
			name:             "ErrDecodeAuth",
			authToken:        "nonempty",
			errDecodeAuth:    errors.New("decode auth failed"),
			authDecoded:      token.Auth{},
			stateToken:       "",
			errDecodeState:   nil,
			stateDecoded:     token.State{},
			errValidateColNo: nil,
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         "",
			wantStatus:       http.StatusUnauthorized,
			assertFunc:       assert.OnResErr("Invalid auth token."),
		},
		{
			name:             "NotAdmin",
			authToken:        "nonempty",
			errDecodeAuth:    nil,
			authDecoded:      token.Auth{IsAdmin: false},
			stateToken:       "",
			errDecodeState:   nil,
			stateDecoded:     token.State{},
			errValidateColNo: nil,
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         "",
			wantStatus:       http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can edit tasks.",
			),
		},
		{
			name:             "NoState",
			authToken:        "nonempty",
			errDecodeAuth:    nil,
			authDecoded:      token.Auth{IsAdmin: true},
			stateToken:       "",
			errDecodeState:   nil,
			stateDecoded:     token.State{},
			errValidateColNo: nil,
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         "",
			wantStatus:       http.StatusBadRequest,
			assertFunc:       assert.OnResErr("State token not found."),
		},
		{
			name:             "ErrDecodeState",
			authToken:        "nonempty",
			errDecodeAuth:    nil,
			authDecoded:      token.Auth{IsAdmin: true},
			stateToken:       "nonempty",
			errDecodeState:   errors.New("decode state failed"),
			stateDecoded:     token.State{},
			errValidateColNo: nil,
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         "",
			wantStatus:       http.StatusBadRequest,
			assertFunc:       assert.OnResErr("Invalid state token."),
		},
		{
			name:             "NoAccess",
			authToken:        "nonempty",
			errDecodeAuth:    nil,
			authDecoded:      token.Auth{IsAdmin: true},
			stateToken:       "nonempty",
			errDecodeState:   nil,
			stateDecoded:     token.State{},
			errValidateColNo: nil,
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         "",
			wantStatus:       http.StatusBadRequest,
			assertFunc:       assert.OnResErr("Invalid task ID."),
		},
		{
			name:           "ColNoInvalid",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: token.State{Boards: []token.Board{{
				ID:      "1",
				Columns: []token.Column{{Tasks: []token.Task{{ID: "taskid"}}}},
			}}},
			errValidateColNo: errors.New("err validate column number"),
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         "",
			wantStatus:       http.StatusBadRequest,
			assertFunc:       assert.OnResErr("Invalid column number."),
		},
		{
			name:           "TaskNotFound",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true, TeamID: "1"},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: token.State{Boards: []token.Board{{
				ID:      "1",
				Columns: []token.Column{{Tasks: []token.Task{{ID: "taskid"}}}},
			}}},
			errValidateColNo: nil,
			errUpdateTasks:   db.ErrNoItem,
			errEncodeState:   nil,
			outState:         "",
			wantStatus:       http.StatusNotFound,
			assertFunc:       assert.OnResErr("Task not found."),
		},
		{
			name:           "ErrUpdateTasks",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true, TeamID: "1"},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: token.State{Boards: []token.Board{{
				ID:      "1",
				Columns: []token.Column{{Tasks: []token.Task{{ID: "taskid"}}}},
			}}},
			errValidateColNo: nil,
			errUpdateTasks:   errors.New("update tasks failed"),
			errEncodeState:   nil,
			outState:         "",
			wantStatus:       http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr("update tasks failed"),
		},
		{
			name:           "ErrEncodeState",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true, TeamID: "1"},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: token.State{Boards: []token.Board{{
				ID:      "1",
				Columns: []token.Column{{Tasks: []token.Task{{ID: "taskid"}}}},
			}}},
			errValidateColNo: nil,
			errUpdateTasks:   nil,
			errEncodeState:   errors.New("encode state failed"),
			outState:         "",
			wantStatus:       http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr("encode state failed"),
		},
		{
			name:           "OK",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true, TeamID: "1"},
			stateToken:     "nonempty",
			errDecodeState: nil,
			stateDecoded: token.State{Boards: []token.Board{{
				ID:      "1",
				Columns: []token.Column{{Tasks: []token.Task{{ID: "taskid"}}}},
			}}},
			errValidateColNo: nil,
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         "aklsdjhfalks",
			wantStatus:       http.StatusOK,
			assertFunc: func(t *testing.T, r *http.Response, _ string) {
				ck := r.Cookies()[0]
				assert.Equal(t.Error, ck.Name, "state-token")
				assert.Equal(t.Error, ck.Value, "aklsdjhfalks")
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Res = c.authDecoded
			decodeAuth.Err = c.errDecodeAuth
			decodeState.Res = c.stateDecoded
			decodeState.Err = c.errDecodeState
			colNoVdtor.Err = c.errValidateColNo
			tasksUpdater.Err = c.errUpdateTasks
			encodeState.Err = c.errEncodeState
			encodeState.Res = c.outState

			// Prepare request and response recorder.
			r := httptest.NewRequest("", "/", strings.NewReader(`[{
                "id": "taskid",
                "order": 3,
                "column": 0
            }]`))
			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name: "auth-token", Value: c.authToken,
				})
			}
			if c.stateToken != "" {
				r.AddCookie(&http.Cookie{
					Name: "state-token", Value: c.stateToken,
				})
			}
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.Handle(w, r, "")
			res := w.Result()

			// Assert on the status code.
			assert.Equal(t.Error, res.StatusCode, c.wantStatus)

			// Run case-specific assertions.
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
