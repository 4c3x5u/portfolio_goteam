//go:build utest

package column

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestHandler(t *testing.T) {
	idValidator := &api.FakeStringValidator{}
	columnSelector := &columnTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	userSelector := &userTable.FakeSelector{}
	columnUpdater := &columnTable.FakeUpdater{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(
		idValidator,
		columnSelector,
		boardSelector,
		userSelector,
		columnUpdater,
		log,
	)

	for _, c := range []struct {
		name            string
		idValidatorErr  error
		column          columnTable.Record
		selectColumnErr error
		board           boardTable.Record
		selectBoardErr  error
		user            userTable.Record
		selectUserErr   error
		updateColumnErr error
		wantStatusCode  int
		assertFunc      func(*testing.T, *http.Response, string)
	}{
		{
			name:            "IDValidatorErr",
			idValidatorErr:  errors.New("invalid id"),
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			user:            userTable.Record{},
			selectUserErr:   nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusBadRequest,
			assertFunc:      assert.OnResErr("invalid id"),
		},
		{
			name:            "ColumnNotFound",
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: sql.ErrNoRows,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			user:            userTable.Record{},
			selectUserErr:   nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusNotFound,
			assertFunc:      assert.OnResErr("Column not found."),
		},
		{
			name:            "ColumnSelectorErr",
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: sql.ErrConnDone,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			user:            userTable.Record{},
			selectUserErr:   nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:            "BoardNotFound",
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  sql.ErrNoRows,
			user:            userTable.Record{},
			selectUserErr:   sql.ErrNoRows,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusNotFound,
			assertFunc:      assert.OnResErr("Board not found."),
		},
		{
			name:            "BoardSelectorErr",
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  sql.ErrConnDone,
			user:            userTable.Record{},
			selectUserErr:   sql.ErrNoRows,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:            "UserNotFound",
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			user:            userTable.Record{},
			selectUserErr:   sql.ErrNoRows,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusUnauthorized,
			assertFunc:      assert.OnResErr("Username is not recognised."),
		},
		{
			name:            "UserSelectorErr",
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			user:            userTable.Record{},
			selectUserErr:   sql.ErrConnDone,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:            "NotAdmin",
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			user:            userTable.Record{IsAdmin: false},
			selectUserErr:   nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can move tasks.",
			),
		},
		{
			name:            "BoardWrongTeam",
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{TeamID: 1},
			selectBoardErr:  nil,
			user:            userTable.Record{IsAdmin: true, TeamID: 2},
			selectUserErr:   nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:            "TaskNotFound",
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{TeamID: 1},
			selectBoardErr:  nil,
			user:            userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:   nil,
			updateColumnErr: sql.ErrNoRows,
			wantStatusCode:  http.StatusNotFound,
			assertFunc:      assert.OnResErr("Task not found."),
		},
		{
			name:            "ColumnUpdaterErr",
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{TeamID: 1},
			selectBoardErr:  nil,
			user:            userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:   nil,
			updateColumnErr: sql.ErrConnDone,
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:            "OK",
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{TeamID: 1},
			selectBoardErr:  nil,
			user:            userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:   nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusOK,
			assertFunc:      func(_ *testing.T, _ *http.Response, _ string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for sut's dependencies.
			idValidator.Err = c.idValidatorErr
			columnSelector.Column = c.column
			columnSelector.Err = c.selectColumnErr
			boardSelector.Board = c.board
			boardSelector.Err = c.selectBoardErr
			userSelector.User = c.user
			userSelector.Err = c.selectUserErr
			columnUpdater.Err = c.updateColumnErr

			// Prepare request and response recorder.
			tasks, err := json.Marshal([]map[string]int{{"id": 0, "order": 0}})
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPatch, "", bytes.NewReader(tasks),
			)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.Handle(w, req, "")
			res := w.Result()

			// Assert on the status code.
			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			// Run case-specific assertions.
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
