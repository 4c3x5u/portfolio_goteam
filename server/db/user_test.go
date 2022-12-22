package db

import (
	"database/sql"
	"testing"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreatorDBUser(t *testing.T) {
	query := `SELECT username FROM users WHERE username = \$1`

	t.Run("CreatedDBUser", func(t *testing.T) {
		db, mock, def := setup(t)
		defer def(db)
		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)
		mock.ExpectClose()
		sut, wantExists := NewExistorUser(db), false

		gotExists, err := sut.Exists("")

		assert.Nil(t, err)
		assert.Equal(t, wantExists, gotExists)
	})

	t.Run("errCreatorUsernameTaken", func(t *testing.T) {
		db, mock, def := setup(t)
		defer def(db)
		mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow(""))
		mock.ExpectClose()
		sut, wantExists := NewExistorUser(db), true

		gotExists, err := sut.Exists("")

		assert.Nil(t, err)
		assert.Equal(t, wantExists, gotExists)
	})

	t.Run("ErrConnDone", func(t *testing.T) {
		db, mock, def := setup(t)
		defer def(db)
		mock.
			ExpectQuery(`SELECT username FROM users WHERE username = \$1`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectClose()
		sut, wantExists, wantErr := NewExistorUser(db), false, sql.ErrConnDone

		gotExists, err := sut.Exists("")

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
