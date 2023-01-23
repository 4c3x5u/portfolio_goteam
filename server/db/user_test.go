package db

import (
	"errors"
	"testing"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserSelector(t *testing.T) {
	id := "bob123"
	query := `SELECT id, password FROM app.\"user\" WHERE id = \$1`

	t.Run("Error", func(t *testing.T) {
		wantErr := errors.New("user inserter error")

		db, mock, teardown := setUpDBTest(t)
		defer teardown()
		mock.ExpectQuery(query).WithArgs(id).WillReturnError(wantErr)
		mock.ExpectClose()

		sut := NewUserSelector(db)

		_, err := sut.Select(id)
		if err = assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		wantPwd := "Myp4ssword!"

		db, mock, teardown := setUpDBTest(t)
		defer teardown()
		mock.ExpectQuery(query).WithArgs(id).WillReturnRows(
			mock.NewRows([]string{"id", "password"}).AddRow(id, wantPwd),
		)
		mock.ExpectClose()

		sut := NewUserSelector(db)

		user, err := sut.Select(id)
		if err != nil {
			t.Fatal(err)
		}
		if err = assert.Equal(id, user.ID); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(wantPwd, string(user.Password)); err != nil {
			t.Error(err)
		}
	})
}

func TestUserInserter(t *testing.T) {
	id, password := "bob123", []byte("hashedpwd")
	query := `INSERT INTO app.\"user\"\(id, password\) VALUES \(\$1, \$2\)`

	t.Run("Error", func(t *testing.T) {
		wantErr := errors.New("db: fatal error")
		db, mock, teardown := setUpDBTest(t)
		defer teardown()
		mock.
			ExpectExec(query).
			WithArgs(id, string(password)).
			WillReturnError(wantErr)
		sut := NewUserInserter(db)

		err := sut.Insert(NewUser(id, password))

		if err = assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		db, mock, teardown := setUpDBTest(t)
		defer teardown()
		mock.
			ExpectExec(query).
			WithArgs(id, string(password)).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sut := NewUserInserter(db)

		err := sut.Insert(NewUser(id, password))

		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
	})
}
