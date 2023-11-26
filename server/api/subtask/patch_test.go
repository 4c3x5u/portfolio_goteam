//go:build utest

package subtask

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	subtaskTable "github.com/kxplxn/goteam/server/dbaccess/subtask"
	taskTable "github.com/kxplxn/goteam/server/dbaccess/task"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestPATCHHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	userSelector := &userTable.FakeSelector{}
	idValidator := &api.FakeStringValidator{}
	subtaskSelector := &subtaskTable.FakeSelector{}
	taskSelector := &taskTable.FakeSelector{}
	columnSelector := &columnTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	subtaskUpdater := &subtaskTable.FakeUpdater{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(
		userSelector,
		idValidator,
		subtaskSelector,
		taskSelector,
		columnSelector,
		boardSelector,
		subtaskUpdater,
		log,
	)

	for _, c := range []struct {
		name               string
		user               userTable.Record
		userSelectorErr    error
		idValidatorErr     error
		subtaskSelectorErr error
		taskSelectorErr    error
		columnSelectorErr  error
		board              boardTable.Record
		boardSelectorErr   error
		subtaskUpdaterErr  error
		wantStatusCode     int
		assertFunc         func(*testing.T, *http.Response, string)
	}{
		{
			name:               "UserNotRecognised",
			user:               userTable.Record{},
			userSelectorErr:    sql.ErrNoRows,
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusUnauthorized,
			assertFunc:         assert.OnResErr("Username is not recognised."),
		},
		{
			name:               "UserSelectorErr",
			user:               userTable.Record{},
			userSelectorErr:    sql.ErrConnDone,
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:               "NotAdmin",
			user:               userTable.Record{IsAdmin: false},
			userSelectorErr:    nil,
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can edit subtasks.",
			),
		},
		{
			name:               "IDEmpty",
			user:               userTable.Record{IsAdmin: true},
			userSelectorErr:    nil,
			idValidatorErr:     api.ErrStrEmpty,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusBadRequest,
			assertFunc:         assert.OnResErr("Subtask ID cannot be empty."),
		},
		{
			name:               "IDNotInt",
			user:               userTable.Record{IsAdmin: true},
			userSelectorErr:    nil,
			idValidatorErr:     api.ErrStrNotInt,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusBadRequest,
			assertFunc:         assert.OnResErr("Subtask ID must be an integer."),
		},
		{
			name:               "IDUnexpectedErr",
			user:               userTable.Record{IsAdmin: true},
			userSelectorErr:    nil,
			idValidatorErr:     api.ErrStrTooLong,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(api.ErrStrTooLong.Error()),
		},
		{
			name:               "SubtaskSelectorErr",
			user:               userTable.Record{IsAdmin: true},
			userSelectorErr:    nil,
			idValidatorErr:     nil,
			subtaskSelectorErr: sql.ErrConnDone,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:               "SubtaskNotFound",
			user:               userTable.Record{IsAdmin: true},
			userSelectorErr:    nil,
			idValidatorErr:     nil,
			subtaskSelectorErr: sql.ErrNoRows,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusNotFound,
			assertFunc:         assert.OnResErr("Subtask not found."),
		},
		{
			name:               "TaskSelectorErr",
			user:               userTable.Record{IsAdmin: true},
			userSelectorErr:    nil,
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    sql.ErrNoRows,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:               "ColumnSelectorErr",
			user:               userTable.Record{IsAdmin: true},
			userSelectorErr:    nil,
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  sql.ErrNoRows,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:               "BoardSelectorErr",
			user:               userTable.Record{IsAdmin: true},
			userSelectorErr:    nil,
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   sql.ErrNoRows,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:               "BoardWrongTeam",
			user:               userTable.Record{IsAdmin: true, TeamID: 2},
			userSelectorErr:    nil,
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{TeamID: 1},
			boardSelectorErr:   nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:               "SubtaskUpdaterErr",
			user:               userTable.Record{IsAdmin: true, TeamID: 1},
			userSelectorErr:    nil,
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{TeamID: 1},
			boardSelectorErr:   nil,
			subtaskUpdaterErr:  sql.ErrNoRows,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:               "Success",
			user:               userTable.Record{IsAdmin: true, TeamID: 1},
			userSelectorErr:    nil,
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{TeamID: 1},
			boardSelectorErr:   nil,
			wantStatusCode:     http.StatusOK,
			assertFunc:         func(*testing.T, *http.Response, string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userSelector.Rec = c.user
			userSelector.Err = c.userSelectorErr
			idValidator.Err = c.idValidatorErr
			subtaskSelector.Err = c.subtaskSelectorErr
			taskSelector.Err = c.taskSelectorErr
			columnSelector.Err = c.columnSelectorErr
			boardSelector.Board = c.board
			boardSelector.Err = c.boardSelectorErr
			subtaskUpdater.Err = c.subtaskUpdaterErr

			reqBody, err := json.Marshal(map[string]any{"done": false})
			if err != nil {
				t.Fatal(err)
			}
			r, err := http.NewRequest("", "?id=", bytes.NewReader(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			sut.Handle(w, r, "")
			res := w.Result()

			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			c.assertFunc(t, res, log.InMessage)
		})
	}
}
