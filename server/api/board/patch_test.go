//go:build utest

package board

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
	"github.com/kxplxn/goteam/server/dbaccess"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestPATCHHandler tests the Handle method of PATCHHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	log := &pkgLog.FakeErrorer{}
	idValidator := &api.FakeStringValidator{}
	nameValidator := &api.FakeStringValidator{}
	boardSelector := &boardTable.FakeSelector{}
	userSelector := &userTable.FakeSelector{}
	boardUpdater := &dbaccess.FakeUpdater{}
	sut := NewPATCHHandler(
		idValidator,
		nameValidator,
		boardSelector,
		userSelector,
		boardUpdater,
		log,
	)

	for _, c := range []struct {
		name             string
		idValidatorErr   error
		nameValidatorErr error
		board            boardTable.Record
		boardSelectorErr error
		user             userTable.Record
		selectUserErr    error
		boardUpdaterErr  error
		wantStatusCode   int
		assertFunc       func(*testing.T, *http.Response, string)
	}{
		{
			name:             "IDValidatorErr",
			idValidatorErr:   errors.New("Board ID cannot be empty."),
			nameValidatorErr: nil,
			board:            boardTable.Record{},
			boardSelectorErr: nil,
			user:             userTable.Record{},
			selectUserErr:    nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Board ID cannot be empty.",
			),
		},
		{
			name:             "NameValidatorErr",
			idValidatorErr:   nil,
			nameValidatorErr: errors.New("Board name cannot be empty."),
			board:            boardTable.Record{},
			boardSelectorErr: nil,
			user:             userTable.Record{},
			selectUserErr:    nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Board name cannot be empty.",
			),
		},
		{
			name:             "BoardNotFound",
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{},
			boardSelectorErr: sql.ErrNoRows,
			user:             userTable.Record{},
			selectUserErr:    nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusNotFound,
			assertFunc:       assert.OnResErr("Board not found."),
		},
		{
			name:             "BoardSelectorErr",
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{},
			boardSelectorErr: sql.ErrConnDone,
			user:             userTable.Record{},
			selectUserErr:    nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:             "UserNotFound",
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{},
			boardSelectorErr: nil,
			user:             userTable.Record{},
			selectUserErr:    sql.ErrNoRows,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusUnauthorized,
			assertFunc:       assert.OnResErr("Username is not recognised."),
		},
		{
			name:             "UserSelectorErr",
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{},
			boardSelectorErr: nil,
			user:             userTable.Record{},
			selectUserErr:    sql.ErrConnDone,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:             "WrongTeamID",
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{TeamID: 1},
			boardSelectorErr: nil,
			user:             userTable.Record{IsAdmin: true, TeamID: 2},
			selectUserErr:    nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:             "NotAdmin",
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{TeamID: 1},
			boardSelectorErr: nil,
			user:             userTable.Record{IsAdmin: false, TeamID: 1},
			selectUserErr:    nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only board admins can edit the board.",
			),
		},
		{
			name:             "BoardUpdaterErr",
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{TeamID: 1},
			boardSelectorErr: nil,
			user:             userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:    nil,
			boardUpdaterErr:  sql.ErrNoRows,
			wantStatusCode:   http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrNoRows.Error(),
			),
		},
		{
			name:             "Success",
			idValidatorErr:   nil,
			nameValidatorErr: nil,
			board:            boardTable.Record{TeamID: 1},
			boardSelectorErr: nil,
			user:             userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:    nil,
			boardUpdaterErr:  nil,
			wantStatusCode:   http.StatusOK,
			assertFunc: func(*testing.T, *http.Response, string) {
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			idValidator.Err = c.idValidatorErr
			nameValidator.Err = c.nameValidatorErr
			boardSelector.Board = c.board
			boardSelector.Err = c.boardSelectorErr
			userSelector.User = c.user
			userSelector.Err = c.selectUserErr
			boardUpdater.Err = c.boardUpdaterErr

			reqBody, err := json.Marshal(ReqBody{})
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPatch, "", bytes.NewReader(reqBody),
			)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			sut.Handle(w, req, "")
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
