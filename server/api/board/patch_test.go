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

	"server/api"
	"server/assert"
	"server/dbaccess"
	boardTable "server/dbaccess/board"
	userboardTable "server/dbaccess/userboard"
	pkgLog "server/log"
)

// TestPATCHHandler tests the Handle method of PATCHHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	log := &pkgLog.FakeErrorer{}
	idValidator := &api.FakeStringValidator{}
	nameValidator := &api.FakeStringValidator{}
	boardSelector := &boardTable.FakeSelector{}
	userBoardSelector := &userboardTable.FakeSelector{}
	boardUpdater := &dbaccess.FakeUpdater{}
	sut := NewPATCHHandler(
		idValidator,
		nameValidator,
		boardSelector,
		userBoardSelector,
		boardUpdater,
		log,
	)

	for _, c := range []struct {
		name                 string
		idValidatorErr       error
		nameValidatorErr     error
		boardSelectorErr     error
		userIsAdmin          bool
		userBoardSelectorErr error
		boardUpdaterErr      error
		wantStatusCode       int
		assertFunc           func(*testing.T, *http.Response, string)
	}{
		{
			name:                 "IDValidatorErr",
			idValidatorErr:       errors.New("Board ID cannot be empty."),
			nameValidatorErr:     nil,
			boardSelectorErr:     nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			boardUpdaterErr:      nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Board ID cannot be empty.",
			),
		},
		{
			name:                 "NameValidatorErr",
			idValidatorErr:       nil,
			nameValidatorErr:     errors.New("Board name cannot be empty."),
			boardSelectorErr:     nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			boardUpdaterErr:      nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc: assert.OnResErr(
				"Board name cannot be empty.",
			),
		},
		{
			name:                 "BoardNotFound",
			idValidatorErr:       nil,
			nameValidatorErr:     nil,
			boardSelectorErr:     sql.ErrNoRows,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			boardUpdaterErr:      nil,
			wantStatusCode:       http.StatusNotFound,
			assertFunc:           assert.OnResErr("Board not found."),
		},
		{
			name:                 "BoardSelectorErr",
			idValidatorErr:       nil,
			nameValidatorErr:     nil,
			boardSelectorErr:     sql.ErrConnDone,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			boardUpdaterErr:      nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                 "UserDoesNotHaveAccess",
			idValidatorErr:       nil,
			nameValidatorErr:     nil,
			boardSelectorErr:     nil,
			userIsAdmin:          false,
			userBoardSelectorErr: sql.ErrNoRows,
			boardUpdaterErr:      nil,
			wantStatusCode:       http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:                 "UserBoardSelectorErr",
			idValidatorErr:       nil,
			nameValidatorErr:     nil,
			boardSelectorErr:     nil,
			userIsAdmin:          false,
			userBoardSelectorErr: sql.ErrConnDone,
			boardUpdaterErr:      nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                 "NotAdmin",
			idValidatorErr:       nil,
			nameValidatorErr:     nil,
			boardSelectorErr:     nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			boardUpdaterErr:      nil,
			wantStatusCode:       http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only board admins can edit the board.",
			),
		},
		{
			name:                 "BoardUpdaterErr",
			idValidatorErr:       nil,
			nameValidatorErr:     nil,
			boardSelectorErr:     nil,
			userIsAdmin:          true,
			userBoardSelectorErr: nil,
			boardUpdaterErr:      sql.ErrNoRows,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrNoRows.Error(),
			),
		},
		{
			name:                 "Success",
			idValidatorErr:       nil,
			nameValidatorErr:     nil,
			boardSelectorErr:     nil,
			userIsAdmin:          true,
			userBoardSelectorErr: nil,
			boardUpdaterErr:      nil,
			wantStatusCode:       http.StatusOK,
			assertFunc: func(*testing.T, *http.Response, string) {
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			idValidator.Err = c.idValidatorErr
			nameValidator.Err = c.nameValidatorErr
			boardSelector.Err = c.boardSelectorErr
			userBoardSelector.IsAdmin = c.userIsAdmin
			userBoardSelector.Err = c.userBoardSelectorErr
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
