//go:build utest

package dbaccess

import (
	"errors"
	"testing"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestUserSelector tests the Select method of UserSelector to assert that it
// sends the correct query to the database with the correct arguments, and
// returns whatever error occurs.
func TestUserSelector(t *testing.T) {
	const (
		username = "bob123"
		query    = `SELECT password FROM app.\"user\" WHERE username = \$1`
	)

	t.Run("Error", func(t *testing.T) {
		wantErr := errors.New("user inserter error")

		db, mock, teardown := SetUpDBTest(t)
		defer teardown()
		mock.ExpectQuery(query).WithArgs(username).WillReturnError(wantErr)
		mock.ExpectClose()

		sut := NewUserSelector(db)

		_, err := sut.Select(username)
		if err = assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		wantPwd := "Myp4ssword!"

		db, mock, teardown := SetUpDBTest(t)
		defer teardown()
		mock.ExpectQuery(query).WithArgs(username).WillReturnRows(
			mock.NewRows([]string{"password"}).AddRow(wantPwd),
		)
		mock.ExpectClose()

		sut := NewUserSelector(db)

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

// TestUserInserter tests the Insert method of UserInserter to assert that it
// sends the correct query to the database with the correct arguments, and
// returns whatever error occurs.
func TestUserInserter(t *testing.T) {
	const (
		username = "bob123"
		pwdHash  = "asd..fasdf.asdfa/sdf.asdfa.sdfa"
		query    = `INSERT INTO app.\"user\"\(username, password\) ` +
			`VALUES \(\$1, \$2\)`
	)

	t.Run("Error", func(t *testing.T) {
		wantErr := errors.New("db: fatal error")
		db, mock, teardown := SetUpDBTest(t)
		defer teardown()
		mock.
			ExpectExec(query).
			WithArgs(username, pwdHash).
			WillReturnError(wantErr)
		sut := NewUserInserter(db)

		err := sut.Insert(NewUser(username, []byte(pwdHash)))

		if err = assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		db, mock, teardown := SetUpDBTest(t)
		defer teardown()
		mock.
			ExpectExec(query).
			WithArgs(username, pwdHash).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sut := NewUserInserter(db)

		err := sut.Insert(NewUser(username, []byte(pwdHash)))

		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
	})
}
