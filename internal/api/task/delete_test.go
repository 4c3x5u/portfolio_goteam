//go:build utest

package task

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/dbaccess"
	boardTable "github.com/kxplxn/goteam/pkg/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/pkg/dbaccess/column"
	taskTable "github.com/kxplxn/goteam/pkg/dbaccess/task"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// TestDeleteHandler tests the Handle method of DeleteHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDeleteHandler(t *testing.T) {
	decodeAuth := &token.FakeDecode[token.Auth]{}
	idValidator := &api.FakeStringValidator{}
	taskSelector := &taskTable.FakeSelector{}
	columnSelector := &columnTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	taskDeleter := &dbaccess.FakeDeleter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDeleteHandler(
		decodeAuth.Func,
		idValidator,
		taskSelector,
		columnSelector,
		boardSelector,
		taskDeleter,
		log,
	)

	for _, c := range []struct {
		authToken       string
		authDecoded     token.Auth
		errDecodeAuth   error
		name            string
		errValidateID   error
		errSelectTask   error
		errSelectColumn error
		board           boardTable.Record
		errSelectBoard  error
		errDeleteTask   error
		wantStatus      int
		assertFunc      func(*testing.T, *http.Response, string)
	}{
		{
			name:            "NoAuth",
			authToken:       "",
			authDecoded:     token.Auth{},
			errDecodeAuth:   nil,
			errValidateID:   nil,
			errSelectTask:   nil,
			errSelectColumn: nil,
			board:           boardTable.Record{},
			errSelectBoard:  nil,
			errDeleteTask:   nil,
			wantStatus:      http.StatusUnauthorized,
			assertFunc:      assert.OnResErr("Auth token not found."),
		},
		{
			name:            "ErrEncodeAuth",
			authToken:       "nonempty",
			authDecoded:     token.Auth{},
			errDecodeAuth:   errors.New("encode auth failed"),
			errValidateID:   nil,
			errSelectTask:   nil,
			errSelectColumn: nil,
			board:           boardTable.Record{},
			errSelectBoard:  nil,
			errDeleteTask:   nil,
			wantStatus:      http.StatusUnauthorized,
			assertFunc:      assert.OnResErr("Invalid auth token."),
		},
		{
			name:            "NotAdmin",
			authToken:       "nonempty",
			authDecoded:     token.Auth{IsAdmin: false},
			errDecodeAuth:   nil,
			errValidateID:   nil,
			errSelectTask:   nil,
			errSelectColumn: nil,
			board:           boardTable.Record{},
			errSelectBoard:  nil,
			errDeleteTask:   nil,
			wantStatus:      http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only board admins can delete tasks.",
			),
		},
		{
			name:            "IDEmpty",
			authToken:       "nonempty",
			authDecoded:     token.Auth{IsAdmin: true},
			errDecodeAuth:   nil,
			errValidateID:   api.ErrEmpty,
			errSelectTask:   nil,
			errSelectColumn: nil,
			board:           boardTable.Record{},
			errSelectBoard:  nil,
			errDeleteTask:   nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc:      assert.OnResErr("Task ID cannot be empty."),
		},
		{
			name:            "IDNotInt",
			authToken:       "nonempty",
			authDecoded:     token.Auth{IsAdmin: true},
			errDecodeAuth:   nil,
			errValidateID:   api.ErrNotInt,
			errSelectTask:   nil,
			errSelectColumn: nil,
			board:           boardTable.Record{},
			errSelectBoard:  nil,
			errDeleteTask:   nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc:      assert.OnResErr("Task ID must be an integer."),
		},
		{
			name:            "IDUnexpectedErr",
			authToken:       "nonempty",
			authDecoded:     token.Auth{IsAdmin: true},
			errDecodeAuth:   nil,
			errValidateID:   api.ErrTooLong,
			errSelectTask:   nil,
			errSelectColumn: nil,
			board:           boardTable.Record{},
			errSelectBoard:  nil,
			errDeleteTask:   nil,
			wantStatus:      http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr(api.ErrTooLong.Error()),
		},
		{
			name:            "TaskSelectorErr",
			authToken:       "nonempty",
			authDecoded:     token.Auth{IsAdmin: true},
			errDecodeAuth:   nil,
			errValidateID:   nil,
			errSelectTask:   sql.ErrConnDone,
			errSelectColumn: nil,
			board:           boardTable.Record{},
			errSelectBoard:  nil,
			errDeleteTask:   nil,
			wantStatus:      http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:            "TaskNotFound",
			authToken:       "nonempty",
			authDecoded:     token.Auth{IsAdmin: true},
			errDecodeAuth:   nil,
			errValidateID:   nil,
			errSelectTask:   sql.ErrNoRows,
			errSelectColumn: nil,
			board:           boardTable.Record{},
			errSelectBoard:  nil,
			errDeleteTask:   nil,
			wantStatus:      http.StatusNotFound,
			assertFunc:      assert.OnResErr("Task not found."),
		},
		{
			name:            "ColumnSelectorErr",
			authToken:       "nonempty",
			authDecoded:     token.Auth{IsAdmin: true},
			errDecodeAuth:   nil,
			errValidateID:   nil,
			errSelectTask:   nil,
			errSelectColumn: sql.ErrNoRows,
			board:           boardTable.Record{},
			errSelectBoard:  nil,
			errDeleteTask:   nil,
			wantStatus:      http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:            "BoardSelectorErr",
			authToken:       "nonempty",
			authDecoded:     token.Auth{IsAdmin: true},
			errDecodeAuth:   nil,
			errValidateID:   nil,
			errSelectTask:   nil,
			errSelectColumn: nil,
			board:           boardTable.Record{},
			errSelectBoard:  sql.ErrNoRows,
			errDeleteTask:   nil,
			wantStatus:      http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:            "BoardWrongTeam",
			authToken:       "nonempty",
			authDecoded:     token.Auth{IsAdmin: true, TeamID: "1"},
			errDecodeAuth:   nil,
			errValidateID:   nil,
			errSelectTask:   nil,
			errSelectColumn: nil,
			board:           boardTable.Record{TeamID: 2},
			errSelectBoard:  nil,
			errDeleteTask:   nil,
			wantStatus:      http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:            "ErrDeleteTask",
			authToken:       "nonempty",
			authDecoded:     token.Auth{IsAdmin: true, TeamID: "2"},
			errDecodeAuth:   nil,
			errValidateID:   nil,
			errSelectTask:   nil,
			errSelectColumn: nil,
			board:           boardTable.Record{TeamID: 2},
			errSelectBoard:  nil,
			errDeleteTask:   sql.ErrNoRows,
			wantStatus:      http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:            "Success",
			authToken:       "nonempty",
			authDecoded:     token.Auth{IsAdmin: true, TeamID: "2"},
			errDecodeAuth:   nil,
			errValidateID:   nil,
			errSelectTask:   nil,
			errSelectColumn: nil,
			board:           boardTable.Record{TeamID: 2},
			errSelectBoard:  nil,
			errDeleteTask:   nil,
			wantStatus:      http.StatusOK,
			assertFunc:      func(*testing.T, *http.Response, string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Decoded = c.authDecoded
			decodeAuth.Err = c.errDecodeAuth
			idValidator.Err = c.errValidateID
			taskSelector.Err = c.errSelectTask
			columnSelector.Err = c.errSelectColumn
			boardSelector.Board = c.board
			boardSelector.Err = c.errSelectBoard
			taskDeleter.Err = c.errDeleteTask

			r := httptest.NewRequest("", "/", nil)
			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  "auth-token",
					Value: c.authToken,
				})
			}

			w := httptest.NewRecorder()

			sut.Handle(w, r, "")
			res := w.Result()

			assert.Equal(t.Error, res.StatusCode, c.wantStatus)

			c.assertFunc(t, res, log.InMessage)
		})
	}
}
