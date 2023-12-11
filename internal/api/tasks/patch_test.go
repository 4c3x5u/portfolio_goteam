//go:build utest

package tasks

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/assert"
	boardTable "github.com/kxplxn/goteam/pkg/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/pkg/dbaccess/column"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// TestPATCHHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	decodeAuth := token.FakeDecode[token.Auth]{}
	idValidator := &api.FakeStringValidator{}
	columnSelector := &columnTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	columnUpdater := &columnTable.FakeUpdater{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(
		decodeAuth.Func,
		idValidator,
		columnSelector,
		boardSelector,
		columnUpdater,
		log,
	)

	for _, c := range []struct {
		name            string
		authToken       string
		errDecodeAuth   error
		auth            token.Auth
		idValidatorErr  error
		column          columnTable.Record
		selectColumnErr error
		board           boardTable.Record
		selectBoardErr  error
		updateColumnErr error
		wantStatusCode  int
		assertFunc      func(*testing.T, *http.Response, string)
	}{
		{
			name:            "NoAuth",
			authToken:       "",
			errDecodeAuth:   nil,
			auth:            token.Auth{},
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusUnauthorized,
			assertFunc:      assert.OnResErr("Auth token not found."),
		},
		{
			name:            "ErrDecodeAuth",
			authToken:       "nonempty",
			errDecodeAuth:   errors.New("decode auth failed"),
			auth:            token.Auth{},
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusUnauthorized,
			assertFunc:      assert.OnResErr("Invalid auth token."),
		},
		{
			name:            "NotAdmin",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			auth:            token.Auth{IsAdmin: false},
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can edit tasks.",
			),
		},
		{
			name:            "IDValidatorErr",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			auth:            token.Auth{IsAdmin: true},
			idValidatorErr:  errors.New("invalid id"),
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc:      assert.OnResErr("invalid id"),
		},
		{
			name:            "ColumnNotFound",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			auth:            token.Auth{IsAdmin: true},
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: sql.ErrNoRows,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusNotFound,
			assertFunc:      assert.OnResErr("Column not found."),
		},
		{
			name:            "ColumnSelectorErr",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			auth:            token.Auth{IsAdmin: true},
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: sql.ErrConnDone,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:            "BoardNotFound",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			auth:            token.Auth{IsAdmin: true},
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  sql.ErrNoRows,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusNotFound,
			assertFunc:      assert.OnResErr("Board not found."),
		},
		{
			name:            "BoardSelectorErr",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			auth:            token.Auth{IsAdmin: true},
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  sql.ErrConnDone,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:            "BoardWrongTeam",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			auth:            token.Auth{IsAdmin: true, TeamID: "1"},
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{TeamID: 2},
			selectBoardErr:  nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:            "TaskNotFound",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			auth:            token.Auth{IsAdmin: true, TeamID: "1"},
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{TeamID: 1},
			selectBoardErr:  nil,
			updateColumnErr: sql.ErrNoRows,
			wantStatusCode:  http.StatusNotFound,
			assertFunc:      assert.OnResErr("Task not found."),
		},
		{
			name:            "ColumnUpdaterErr",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			auth:            token.Auth{IsAdmin: true, TeamID: "1"},
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{TeamID: 1},
			selectBoardErr:  nil,
			updateColumnErr: sql.ErrConnDone,
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:            "OK",
			authToken:       "nonempty",
			errDecodeAuth:   nil,
			auth:            token.Auth{IsAdmin: true, TeamID: "1"},
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{TeamID: 1},
			selectBoardErr:  nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusOK,
			assertFunc:      func(_ *testing.T, _ *http.Response, _ string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			decodeAuth.Res = c.auth
			decodeAuth.Err = c.errDecodeAuth
			idValidator.Err = c.idValidatorErr
			columnSelector.Column = c.column
			columnSelector.Err = c.selectColumnErr
			boardSelector.Board = c.board
			boardSelector.Err = c.selectBoardErr
			columnUpdater.Err = c.updateColumnErr

			// Prepare request and response recorder.
			r := httptest.NewRequest("", "/", strings.NewReader("[]"))
			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  "auth-token",
					Value: c.authToken,
				})
			}
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.Handle(w, r, "")
			res := w.Result()

			// Assert on the status code.
			assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

			// Run case-specific assertions.
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
