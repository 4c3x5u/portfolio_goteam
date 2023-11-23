//go:build utest

package user

import (
	"errors"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
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
		wantErr := errors.New("user inserter error")

		db, mock, teardown := dbaccess.SetUpDBTest(t)
		defer teardown()
		mock.ExpectQuery(query).WithArgs(username).WillReturnError(wantErr)
		mock.ExpectClose()

		sut := NewSelector(db)

		_, err := sut.Select(username)
		if err = assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		wantPwd := "Myp4ssword!"
		wantTeamID := 21
		wantIsAdmin := true

		db, mock, teardown := dbaccess.SetUpDBTest(t)
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
		if err = assert.Equal(username, user.Username); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(wantPwd, string(user.Password)); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(wantTeamID, user.TeamID); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(wantIsAdmin, user.IsAdmin); err != nil {
			t.Error(err)
		}
	})
}
