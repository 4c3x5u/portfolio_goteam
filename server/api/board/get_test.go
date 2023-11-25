//go:build utest

package board

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

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
	teamSelector := &teamTable.FakeSelector{}
	boardSelector := &boardTable.FakeRecursiveSelector{}
	log := &pkgLog.FakeErrorer{}

	sut := NewGETHandler(userSelector, teamSelector, boardSelector, log)

	for _, c := range []struct {
		name            string
		userSelectorErr error
		wantStatusCode  int
		assertFunc      func(*testing.T, *http.Response, string)
	}{
		{
			name:            "UserIsNotRecognised",
			userSelectorErr: sql.ErrNoRows,
			wantStatusCode:  http.StatusUnauthorized,
			assertFunc:      assert.OnLoggedErr(sql.ErrNoRows.Error()),
		},
		{
			name:            "UserSelectorErr",
			userSelectorErr: sql.ErrConnDone,
			wantStatusCode:  http.StatusInternalServerError,
			assertFunc:      assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userSelector.Err = c.userSelectorErr

			w := httptest.NewRecorder()
			sut.Handle(w, nil, "")
			res := w.Result()

			if err := assert.Equal(c.wantStatusCode, res.StatusCode); err != nil {
				t.Error(err)
			}

			c.assertFunc(t, res, log.InMessage)
		})
	}
}
