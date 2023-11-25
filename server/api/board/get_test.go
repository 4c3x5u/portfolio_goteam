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

	t.Run("UserNotRecognised", func(t *testing.T) {
		wantErr := sql.ErrNoRows
		userSelector.Err = wantErr
		wantStatusCode := http.StatusUnauthorized

		w := httptest.NewRecorder()
		sut.Handle(w, nil, "")
		res := w.Result()

		if err := assert.Equal(wantStatusCode, res.StatusCode); err != nil {
			t.Error(err)
		}

		assert.OnLoggedErr(wantErr.Error())(t, res, log.InMessage)
	})
}
