package db

import (
	"database/sql"
	"errors"
	"testing"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestReaderUser(t *testing.T) {
	username := "bob21"
	query := `SELECT password FROM users WHERE username = \$1`

	t.Run("Err", func(t *testing.T) {
		wantErr := sql.ErrNoRows

		db, mock, def := setupTest(t)
		defer def(db)
		mock.ExpectQuery(query).WithArgs(username).WillReturnError(wantErr)
		mock.ExpectClose()

		sut := NewReaderUserPwd(db)

		_, err := sut.Read(username)
		assert.Equal(t, wantErr.Error(), err.Error())
	})

	t.Run("Res", func(t *testing.T) {
		wantPwd := "Myp4ssword!"

		db, mock, def := setupTest(t)
		defer def(db)
		mock.ExpectQuery(query).WithArgs(username).WillReturnRows(
			mock.NewRows([]string{"password"}).AddRow(wantPwd),
		)
		mock.ExpectClose()

		sut := NewReaderUserPwd(db)

		pwd, err := sut.Read(username)
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, wantPwd, string(pwd))
	})
}

func TestExistorUser(t *testing.T) {
	username := "bob21"
	query := `SELECT username FROM users WHERE username = \$1`

	t.Run("ExistsFalse", func(t *testing.T) {
		wantExists := false
		db, mock, def := setupTest(t)
		defer def(db)
		mock.ExpectQuery(query).WithArgs(username).WillReturnError(sql.ErrNoRows)
		mock.ExpectClose()
		sut := NewExistorUser(db)

		gotExists, err := sut.Exists(username)

		assert.Nil(t, err)
		assert.Equal(t, wantExists, gotExists)
	})

	t.Run("ExistsTrue", func(t *testing.T) {
		wantExists := true
		db, mock, def := setupTest(t)
		defer def(db)
		mock.ExpectQuery(query).WithArgs(username).WillReturnRows(
			sqlmock.NewRows([]string{"username"}).AddRow(""),
		)
		mock.ExpectClose()
		sut := NewExistorUser(db)

		gotExists, err := sut.Exists(username)

		assert.Nil(t, err)
		assert.Equal(t, wantExists, gotExists)
	})

	t.Run("ExistsErr", func(t *testing.T) {
		wantExists, wantErr := false, sql.ErrConnDone
		db, mock, def := setupTest(t)
		defer def(db)
		mock.ExpectQuery(query).WithArgs(username).WillReturnError(sql.ErrConnDone)
		mock.ExpectClose()
		sut := NewExistorUser(db)

		gotExists, err := sut.Exists(username)

		assert.Equal(t, wantExists, gotExists)
		assert.Equal(t, wantErr.Error(), err.Error())
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

		err := sut.Create(username, password)

		assert.Nil(t, err)
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

		err := sut.Create(username, password)

		assert.Equal(t, wantErr.Error(), err.Error())
	})
}
