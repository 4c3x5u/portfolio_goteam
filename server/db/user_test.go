package db

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"testing"

	"server/assert"
)

func TestUserReader(t *testing.T) {
	username := "bob21"
	query := `SELECT username, password FROM users WHERE username = \$1`

	t.Run("Error", func(t *testing.T) {
		wantErr := errors.New("user reader error")

		db, mock, teardown := setUpDBMock(t)
		defer teardown()
		mock.ExpectQuery(query).WithArgs(username).WillReturnError(wantErr)
		mock.ExpectClose()

		sut := NewUserReader(db)

		_, err := sut.Read(username)
		if err = assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		wantPwd := "Myp4ssword!"

		db, mock, teardown := setUpDBMock(t)
		defer teardown()
		mock.ExpectQuery(query).WithArgs(username).WillReturnRows(
			mock.NewRows([]string{"username", "password"}).AddRow(username, wantPwd),
		)
		mock.ExpectClose()

		sut := NewUserReader(db)

		user, err := sut.Read(username)
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

func TestUserCreator(t *testing.T) {
	username, password := "bob21", []byte("hashedpwd")
	query := `INSERT INTO users\(username, password\) VALUES \(\$1, \$2\)`

	t.Run("Error", func(t *testing.T) {
		wantErr := errors.New("db: fatal error")
		db, mock, teardown := setUpDBMock(t)
		defer teardown()
		mock.
			ExpectExec(`INSERT INTO users\(username, password\) VALUES \(\$1, \$2\)`).
			WithArgs(username, string(password)).
			WillReturnError(wantErr)
		sut := NewUserCreator(db)

		err := sut.Create(NewUser(username, password))

		if err = assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		db, mock, teardown := setUpDBMock(t)
		defer teardown()
		mock.
			ExpectExec(query).
			WithArgs(username, string(password)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sut := NewUserCreator(db)

		err := sut.Create(NewUser(username, password))

		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
	})
}
