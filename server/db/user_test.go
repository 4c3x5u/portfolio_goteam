package db

import (
	"errors"
	"testing"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestReaderUser(t *testing.T) {
	username := "bob21"
	query := `SELECT username, password FROM users WHERE username = \$1`

	t.Run("Err", func(t *testing.T) {
		wantErr := errors.New("reader fatal error")

		db, mock, def := setupTest(t)
		defer def(db)
		mock.ExpectQuery(query).WithArgs(username).WillReturnError(wantErr)
		mock.ExpectClose()

		sut := NewReaderUser(db)

		_, err := sut.Read(username)
		if err = assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Res", func(t *testing.T) {
		wantPwd := "Myp4ssword!"

		db, mock, def := setupTest(t)
		defer def(db)
		mock.ExpectQuery(query).WithArgs(username).WillReturnRows(
			mock.NewRows([]string{"username", "password"}).AddRow(username, wantPwd),
		)
		mock.ExpectClose()

		sut := NewReaderUser(db)

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

func TestCreatorUser(t *testing.T) {
	username, password := "bob21", []byte("hashedpwd")
	query := `INSERT INTO users\(username, password\) VALUES \(\$1, \$2\)`

	t.Run("CreateOK", func(t *testing.T) {
		db, mock, def := setupTest(t)
		defer def(db)
		mock.
			ExpectExec(query).
			WithArgs(username, string(password)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sut := NewCreatorUser(db)

		err := sut.Create(NewUser(username, password))

		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
	})

	t.Run("CreateErr", func(t *testing.T) {
		wantErr := errors.New("db: fatal error")
		db, mock, def := setupTest(t)
		defer def(db)
		mock.
			ExpectExec(`INSERT INTO users\(username, password\) VALUES \(\$1, \$2\)`).
			WithArgs(username, string(password)).
			WillReturnError(wantErr)
		sut := NewCreatorUser(db)

		err := sut.Create(NewUser(username, password))

		if err = assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})
}
