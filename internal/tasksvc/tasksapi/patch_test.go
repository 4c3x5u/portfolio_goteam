//go:build utest

package tasksapi

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
)

func TestPatchHandler(t *testing.T) {
	authDecoder := &cookie.FakeDecoder[cookie.Auth]{}
	colNoVdtor := &api.FakeIntValidator{}
	tasksUpdater := &db.FakeUpdater[[]tasktbl.Task]{}
	log := &log.FakeErrorer{}
	sut := NewPatchHandler(
		authDecoder,
		colNoVdtor,
		tasksUpdater,
		log,
	)

	for _, c := range []struct {
		name             string
		authToken        string
		errDecodeAuth    error
		authDecoded      cookie.Auth
		errValidateColNo error
		errUpdateTasks   error
		errEncodeState   error
		outState         http.Cookie
		wantStatus       int
		assertFunc       func(*testing.T, *http.Response, []any)
	}{
		{
			name:             "NoAuth",
			authToken:        "",
			errDecodeAuth:    nil,
			authDecoded:      cookie.Auth{},
			errValidateColNo: nil,
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         http.Cookie{},
			wantStatus:       http.StatusUnauthorized,
			assertFunc:       assert.OnRespErr("Auth token not found."),
		},
		{
			name:             "ErrDecodeAuth",
			authToken:        "nonempty",
			errDecodeAuth:    errors.New("decode auth failed"),
			authDecoded:      cookie.Auth{},
			errValidateColNo: nil,
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         http.Cookie{},
			wantStatus:       http.StatusUnauthorized,
			assertFunc:       assert.OnRespErr("Invalid auth token."),
		},
		{
			name:             "NotAdmin",
			authToken:        "nonempty",
			errDecodeAuth:    nil,
			authDecoded:      cookie.Auth{IsAdmin: false},
			errValidateColNo: nil,
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         http.Cookie{},
			wantStatus:       http.StatusForbidden,
			assertFunc: assert.OnRespErr(
				"Only team admins can edit tasks.",
			),
		},
		{
			name:             "ColNoInvalid",
			authToken:        "nonempty",
			errDecodeAuth:    nil,
			authDecoded:      cookie.Auth{IsAdmin: true},
			errValidateColNo: errors.New("err validate column number"),
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         http.Cookie{},
			wantStatus:       http.StatusBadRequest,
			assertFunc:       assert.OnRespErr("Invalid column number."),
		},
		{
			name:             "TaskNotFound",
			authToken:        "nonempty",
			errDecodeAuth:    nil,
			authDecoded:      cookie.Auth{IsAdmin: true, TeamID: "1"},
			errValidateColNo: nil,
			errUpdateTasks:   db.ErrNoItem,
			errEncodeState:   nil,
			outState:         http.Cookie{},
			wantStatus:       http.StatusNotFound,
			assertFunc:       assert.OnRespErr("Task not found."),
		},
		{
			name:             "ErrUpdateTasks",
			authToken:        "nonempty",
			errDecodeAuth:    nil,
			authDecoded:      cookie.Auth{IsAdmin: true, TeamID: "1"},
			errValidateColNo: nil,
			errUpdateTasks:   errors.New("update tasks failed"),
			errEncodeState:   nil,
			outState:         http.Cookie{},
			wantStatus:       http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr("update tasks failed"),
		},
		{
			name:             "OK",
			authToken:        "nonempty",
			errDecodeAuth:    nil,
			authDecoded:      cookie.Auth{IsAdmin: true, TeamID: "1"},
			errValidateColNo: nil,
			errUpdateTasks:   nil,
			errEncodeState:   nil,
			outState:         http.Cookie{Name: "foo", Value: "bar"},
			wantStatus:       http.StatusOK,
			assertFunc:       func(*testing.T, *http.Response, []any) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			authDecoder.Res = c.authDecoded
			authDecoder.Err = c.errDecodeAuth
			colNoVdtor.Err = c.errValidateColNo
			tasksUpdater.Err = c.errUpdateTasks
			w := httptest.NewRecorder()
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

			sut.Handle(w, r, "")

			resp := w.Result()
			assert.Equal(t.Error, resp.StatusCode, c.wantStatus)
			c.assertFunc(t, resp, log.Args)
		})
	}
}
