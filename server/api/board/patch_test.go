//go:build utest

package board

import (
	"bytes"
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestPATCHHandler tests the Handle method of PATCHHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	userSelector := &userTable.FakeSelector{}
	idValidator := &api.FakeStringValidator{}
	nameValidator := &api.FakeStringValidator{}
	boardSelector := &boardTable.FakeSelector{}
	boardUpdater := &dbaccess.FakeUpdater{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(
		userSelector,
		idValidator,
		nameValidator,
		boardSelector,
		boardUpdater,
		log,
	)

	for _, c := range []struct {
		name             string
		user             userTable.Record
		selectUserErr    error
		idValidatorErr   error
		nameValidatorErr error
		board            boardTable.Record
		boardSelectorErr error
		boardUpdaterErr  error
		wantStatusCode   int
		assertFunc       func(*testing.T, *http.Response, string)
	}{
		{
			name:             "UserNotRecognised",
			user:             userTable.Record{},
			selectUserErr:    sql.ErrNoRows,
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{},
			boardSelectorErr: nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusUnauthorized,
			assertFunc:       assert.OnResErr("Username is not recognised."),
		},
		{
			name:             "UserSelectorErr",
			user:             userTable.Record{},
			selectUserErr:    sql.ErrConnDone,
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{},
			boardSelectorErr: nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:             "NotAdmin",
			user:             userTable.Record{IsAdmin: false, TeamID: 1},
			selectUserErr:    nil,
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{TeamID: 1},
			boardSelectorErr: nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only team admins can edit the board.",
			),
		},
		{
			name:             "IDValidatorErr",
			user:             userTable.Record{IsAdmin: true},
			selectUserErr:    nil,
			idValidatorErr:   errors.New("Board ID cannot be empty."),
			nameValidatorErr: nil,
			board:            boardTable.Record{},
			boardSelectorErr: nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Board ID cannot be empty.",
			),
		},
		{
			name:             "NameValidatorErr",
			user:             userTable.Record{IsAdmin: true},
			selectUserErr:    nil,
			idValidatorErr:   nil,
			nameValidatorErr: errors.New("Board name cannot be empty."),
			board:            boardTable.Record{},
			boardSelectorErr: nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Board name cannot be empty.",
			),
		},
		{
			name:             "BoardNotFound",
			user:             userTable.Record{IsAdmin: true},
			selectUserErr:    nil,
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{},
			boardSelectorErr: sql.ErrNoRows,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusNotFound,
			assertFunc:       assert.OnResErr("Board not found."),
		},
		{
			name:             "BoardSelectorErr",
			user:             userTable.Record{IsAdmin: true},
			selectUserErr:    nil,
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{},
			boardSelectorErr: sql.ErrConnDone,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:             "WrongTeamID",
			user:             userTable.Record{IsAdmin: true, TeamID: 2},
			selectUserErr:    nil,
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{TeamID: 1},
			boardSelectorErr: nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:             "BoardUpdaterErr",
			user:             userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:    nil,
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{TeamID: 1},
			boardSelectorErr: nil,
			boardUpdaterErr:  sql.ErrNoRows,
			wantStatusCode:   http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrNoRows.Error(),
			),
		},
		{
			name:             "Success",
			user:             userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:    nil,
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{TeamID: 1},
			boardSelectorErr: nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusOK,
			assertFunc: func(*testing.T, *http.Response, string) {
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userSelector.Rec = c.user
			userSelector.Err = c.selectUserErr
			idValidator.Err = c.idValidatorErr
			nameValidator.Err = c.nameValidatorErr
			boardSelector.Board = c.board
			boardSelector.Err = c.boardSelectorErr
			boardUpdater.Err = c.boardUpdaterErr

			req := httptest.NewRequest("", "/", bytes.NewReader([]byte("{}")))
			w := httptest.NewRecorder()

			sut.Handle(w, req, "")
			res := w.Result()

			if err := assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			c.assertFunc(t, res, log.InMessage)
		})
	}
}
