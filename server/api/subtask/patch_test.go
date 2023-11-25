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
	idValidator := &api.FakeStringValidator{}
	subtaskSelector := &subtaskTable.FakeSelector{}
	taskSelector := &taskTable.FakeSelector{}
	columnSelector := &columnTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	userSelector := &userTable.FakeSelector{}
	subtaskUpdater := &subtaskTable.FakeUpdater{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(
		idValidator,
		subtaskSelector,
		taskSelector,
		columnSelector,
		boardSelector,
		userSelector,
		subtaskUpdater,
		log,
	)

	for _, c := range []struct {
		name               string
		idValidatorErr     error
		subtaskSelectorErr error
		taskSelectorErr    error
		columnSelectorErr  error
		board              boardTable.Record
		boardSelectorErr   error
		user               userTable.Record
		userSelectorErr    error
		subtaskUpdaterErr  error
		wantStatusCode     int
		assertFunc         func(*testing.T, *http.Response, string)
	}{
		{
			name:               "IDEmpty",
			idValidatorErr:     api.ErrStrEmpty,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			user:               userTable.Record{},
			userSelectorErr:    nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusBadRequest,
			assertFunc:         assert.OnResErr("Subtask ID cannot be empty."),
		},
		{
			name:               "IDNotInt",
			idValidatorErr:     api.ErrStrNotInt,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			user:               userTable.Record{},
			userSelectorErr:    nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusBadRequest,
			assertFunc:         assert.OnResErr("Subtask ID must be an integer."),
		},
		{
			name:               "IDUnexpectedErr",
			idValidatorErr:     api.ErrStrTooLong,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			user:               userTable.Record{},
			userSelectorErr:    nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(api.ErrStrTooLong.Error()),
		},
		{
			name:               "SubtaskSelectorErr",
			idValidatorErr:     nil,
			subtaskSelectorErr: sql.ErrConnDone,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			user:               userTable.Record{},
			userSelectorErr:    nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:               "SubtaskNotFound",
			idValidatorErr:     nil,
			subtaskSelectorErr: sql.ErrNoRows,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			user:               userTable.Record{},
			userSelectorErr:    nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusNotFound,
			assertFunc:         assert.OnResErr("Subtask not found."),
		},
		{
			name:               "TaskSelectorErr",
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    sql.ErrNoRows,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			user:               userTable.Record{},
			userSelectorErr:    nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:               "ColumnSelectorErr",
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  sql.ErrNoRows,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			user:               userTable.Record{},
			userSelectorErr:    nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:               "BoardSelectorErr",
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   sql.ErrNoRows,
			user:               userTable.Record{},
			userSelectorErr:    nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:               "UserNotRecognised",
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			user:               userTable.Record{},
			userSelectorErr:    sql.ErrNoRows,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusUnauthorized,
			assertFunc:         assert.OnResErr("Username is not recognised."),
		},
		{
			name:               "UserSelectorErr",
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			user:               userTable.Record{},
			userSelectorErr:    sql.ErrConnDone,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:               "NotAdmin",
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{},
			boardSelectorErr:   nil,
			user:               userTable.Record{IsAdmin: false},
			userSelectorErr:    nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can edit subtasks.",
			),
		},
		{
			name:               "BoardWrongTeam",
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{TeamID: 1},
			boardSelectorErr:   nil,
			user:               userTable.Record{IsAdmin: true, TeamID: 2},
			userSelectorErr:    nil,
			subtaskUpdaterErr:  nil,
			wantStatusCode:     http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:               "SubtaskUpdaterErr",
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{TeamID: 1},
			boardSelectorErr:   nil,
			user:               userTable.Record{IsAdmin: true, TeamID: 1},
			userSelectorErr:    nil,
			subtaskUpdaterErr:  sql.ErrNoRows,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:               "Success",
			idValidatorErr:     nil,
			subtaskSelectorErr: nil,
			taskSelectorErr:    nil,
			columnSelectorErr:  nil,
			board:              boardTable.Record{TeamID: 1},
			boardSelectorErr:   nil,
			user:               userTable.Record{IsAdmin: true, TeamID: 1},
			userSelectorErr:    nil,
			wantStatusCode:     http.StatusOK,
			assertFunc:         func(*testing.T, *http.Response, string) {},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			idValidator.Err = c.idValidatorErr
			subtaskSelector.Err = c.subtaskSelectorErr
			taskSelector.Err = c.taskSelectorErr
			columnSelector.Err = c.columnSelectorErr
			boardSelector.Board = c.board
			boardSelector.Err = c.boardSelectorErr
			userSelector.User = c.user
			userSelector.Err = c.userSelectorErr
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
