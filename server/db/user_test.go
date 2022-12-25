package db

import (
	"database/sql"
	"errors"
	"testing"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExistorUser(t *testing.T) {
	username := "bob21"
	query := `SELECT username FROM users WHERE username = \$1`

	t.Run("ExistsFalse", func(t *testing.T) {
		wantExists := false
		db, mock, def := setup(t)
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
		db, mock, def := setup(t)
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
		db, mock, def := setup(t)
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
		db, mock, def := setup(t)
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
		db, mock, def := setup(t)
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

// setup is for setting up the DB and Sqlmock for each test case. The first two
// return values are self-explanatory, and the third is the database teardown
// function to be deferred by the test case.
// TODO: Extract this function when it's needed by another test.
func setup(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func(*sql.DB)) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	return db, mock, func(db *sql.DB) {
		mock.ExpectClose()
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
