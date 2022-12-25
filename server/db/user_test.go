package db

import (
	"database/sql"
	"testing"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestExistorUser(t *testing.T) {
	username := "bob21"
	query := `SELECT username FROM users WHERE username = \$1`

	t.Run("ExistsFalse", func(t *testing.T) {
		db, mock, def := setup(t)
		defer def(db)
		mock.ExpectQuery(query).WithArgs(username).WillReturnError(sql.ErrNoRows)
		mock.ExpectClose()
		sut, wantExists := NewExistorUser(db), false

		gotExists, err := sut.Exists(username)

		assert.Nil(t, err)
		assert.Equal(t, wantExists, gotExists)
	})

	t.Run("ExistsTrue", func(t *testing.T) {
		db, mock, def := setup(t)
		defer def(db)
		mock.ExpectQuery(query).WithArgs(username).WillReturnRows(
			sqlmock.NewRows([]string{"username"}).AddRow(""),
		)
		mock.ExpectClose()
		sut, wantExists := NewExistorUser(db), true

		gotExists, err := sut.Exists(username)

		assert.Nil(t, err)
		assert.Equal(t, wantExists, gotExists)
	})

	t.Run("ExistsErr", func(t *testing.T) {
		db, mock, def := setup(t)
		defer def(db)
		mock.ExpectQuery(query).WithArgs(username).WillReturnError(sql.ErrConnDone)
		mock.ExpectClose()
		sut, wantExists, wantErr := NewExistorUser(db), false, sql.ErrConnDone

		gotExists, err := sut.Exists(username)

		assert.Equal(t, wantExists, gotExists)
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
		if err := db.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
