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

// TestPATCHHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	userSelector := &userTable.FakeSelector{}
	idValidator := &api.FakeStringValidator{}
	columnSelector := &columnTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	columnUpdater := &columnTable.FakeUpdater{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(
		userSelector,
		idValidator,
		columnSelector,
		boardSelector,
		columnUpdater,
		log,
	)

	for _, c := range []struct {
		name            string
		user            userTable.Record
		selectUserErr   error
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
			name:            "UserNotRecognised",
			user:            userTable.Record{},
			selectUserErr:   sql.ErrNoRows,
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusUnauthorized,
			assertFunc:      assert.OnResErr("Username is not recognised."),
		},
		{
			name:            "UserSelectorErr",
			user:            userTable.Record{},
			selectUserErr:   sql.ErrConnDone,
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:            "NotAdmin",
			user:            userTable.Record{IsAdmin: false},
			selectUserErr:   nil,
			idValidatorErr:  nil,
			column:          columnTable.Record{},
			selectColumnErr: nil,
			board:           boardTable.Record{},
			selectBoardErr:  nil,
			updateColumnErr: nil,
			wantStatusCode:  http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can move tasks.",
			),
		},
		{
			name:            "IDValidatorErr",
			user:            userTable.Record{IsAdmin: true},
			selectUserErr:   nil,
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
			user:            userTable.Record{IsAdmin: true},
			selectUserErr:   nil,
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
			user:            userTable.Record{IsAdmin: true},
			selectUserErr:   nil,
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
			user:            userTable.Record{IsAdmin: true},
			selectUserErr:   nil,
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
			user:            userTable.Record{IsAdmin: true},
			selectUserErr:   nil,
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
			user:            userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:   nil,
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
			user:            userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:   nil,
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
			user:            userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:   nil,
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
			user:            userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:   nil,
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
			// Set pre-determinate return values for sut's dependencies.
			userSelector.Rec = c.user
			userSelector.Err = c.selectUserErr
			idValidator.Err = c.idValidatorErr
			columnSelector.Column = c.column
			columnSelector.Err = c.selectColumnErr
			boardSelector.Board = c.board
			boardSelector.Err = c.selectBoardErr
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
