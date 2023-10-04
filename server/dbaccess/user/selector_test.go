//go:build utest

package user

import (
	"errors"
	"testing"

	"server/assert"
	"server/dbaccess"
)

// TestSelector tests the Select method of Selector to assert that it
// sends the correct query to the database with the correct arguments, and
// returns whatever error occurs.
func TestSelector(t *testing.T) {
	const (
		username = "bob123"
		query    = `SELECT password FROM app.\"user\" WHERE username = \$1`
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

	t.Run("Success", func(t *testing.T) {
		wantPwd := "Myp4ssword!"

		db, mock, teardown := dbaccess.SetUpDBTest(t)
		defer teardown()
		mock.ExpectQuery(query).WithArgs(username).WillReturnRows(
			mock.NewRows([]string{"password"}).AddRow(wantPwd),
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
	})
}
