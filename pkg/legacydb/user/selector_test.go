//go:build utest

package user

import (
	"errors"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/legacydb"
)

// TestSelector tests the Select method of Selector to assert that it sends the
// correct query to the database with the correct arguments, and returns all
// data or returns whatever error occurs.
func TestSelector(t *testing.T) {
	const (
		username = "bob123"
		query    = `SELECT password, teamID, isAdmin FROM app.\"user\" ` +
			`WHERE username = \$1`
	)

	t.Run("Error", func(t *testing.T) {
		wantErr := errors.New("error selecting user")

		db, mock, teardown := legacydb.SetUpDBTest(t)
		defer teardown()
		mock.ExpectQuery(query).WithArgs(username).WillReturnError(wantErr)
		mock.ExpectClose()

		sut := NewSelector(db)

		_, err := sut.Select(username)
		assert.Equal(t.Error, wantErr, err)
	})

	t.Run("OK", func(t *testing.T) {
		wantPwd := "Myp4ssword!"
		wantTeamID := 21
		wantIsAdmin := true

		db, mock, teardown := legacydb.SetUpDBTest(t)
		defer teardown()
		mock.ExpectQuery(query).WithArgs(username).WillReturnRows(
			mock.NewRows([]string{"password", "teamID", "isAdmin"}).
				AddRow(wantPwd, wantTeamID, wantIsAdmin),
		)
		mock.ExpectClose()

		sut := NewSelector(db)

		user, err := sut.Select(username)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t.Error, user.Username, username)
		assert.Equal(t.Error, string(user.Password), wantPwd)
		assert.Equal(t.Error, user.TeamID, wantTeamID)
		assert.Equal(t.Error, user.IsAdmin, wantIsAdmin)
	})
}
