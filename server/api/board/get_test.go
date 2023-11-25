//go:build utest

package board

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	teamTable "github.com/kxplxn/goteam/server/dbaccess/team"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestGETHandler tests the Handle method of GETHandler to assert that it
// behaves correctly in all possible scenarios.
func TestGETHandler(t *testing.T) {
	userSelector := &userTable.FakeSelector{}
	idValidator := &api.FakeStringValidator{}
	boardSelector := &boardTable.FakeRecursiveSelector{}
	teamSelector := &teamTable.FakeSelector{}
	log := &pkgLog.FakeErrorer{}

	sut := NewGETHandler(userSelector, idValidator, boardSelector, teamSelector, log)

	for _, c := range []struct {
		name             string
		userSelectorErr  error
		idValidatorErr   error
		boardSelectorErr error
		teamSelectorErr  error
		wantStatusCode   int
		assertFunc       func(*testing.T, *http.Response, string)
	}{
		{
			name:             "UserIsNotRecognised",
			userSelectorErr:  sql.ErrNoRows,
			idValidatorErr:   nil,
			boardSelectorErr: nil,
			teamSelectorErr:  nil,
			wantStatusCode:   http.StatusUnauthorized,
			assertFunc:       assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:             "UserSelectorErr",
			userSelectorErr:  sql.ErrConnDone,
			idValidatorErr:   nil,
			boardSelectorErr: nil,
			teamSelectorErr:  nil,
			wantStatusCode:   http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:             "InvalidID",
			userSelectorErr:  nil,
			idValidatorErr:   errors.New("error invalid id"),
			boardSelectorErr: nil,
			teamSelectorErr:  nil,
			wantStatusCode:   http.StatusBadRequest,
			assertFunc:       func(_ *testing.T, _ *http.Response, _ string) {},
		},
		{
			name:             "BoardNotFound",
			userSelectorErr:  nil,
			idValidatorErr:   nil,
			boardSelectorErr: sql.ErrNoRows,
			teamSelectorErr:  nil,
			wantStatusCode:   http.StatusNotFound,
			assertFunc:       func(_ *testing.T, _ *http.Response, _ string) {},
		},
		{
			name:             "BoardSelectorErr",
			userSelectorErr:  nil,
			idValidatorErr:   nil,
			boardSelectorErr: sql.ErrConnDone,
			teamSelectorErr:  nil,
			wantStatusCode:   http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:             "TeamSelectorErr",
			userSelectorErr:  nil,
			idValidatorErr:   nil,
			boardSelectorErr: nil,
			teamSelectorErr:  sql.ErrNoRows,
			wantStatusCode:   http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userSelector.Err = c.userSelectorErr
			idValidator.Err = c.idValidatorErr
			boardSelector.Err = c.boardSelectorErr
			teamSelector.Err = c.teamSelectorErr

			r, err := http.NewRequest(http.MethodGet, "?boardID=1", nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			sut.Handle(w, r, "")

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
