package register

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"server/assert"
)

func TestCreatorDBUser(t *testing.T) {
	query := `SELECT username FROM users WHERE username = \$1`

	t.Run("CreatedDBUser", func(t *testing.T) {
		db, mock, def := setup(t)
		defer def(db)
		mock.ExpectQuery(query).WillReturnError(sql.ErrNoRows)
		mock.ExpectClose()
		sut := NewCreatorDBUser(db)

		err := sut.CreateUser("", "")

		assert.Nil(t, err)
	})

	t.Run("errCreatorUsernameTaken", func(t *testing.T) {
		db, mock, def := setup(t)
		defer def(db)
		mock.ExpectQuery(query).WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow(""))
		mock.ExpectClose()
		sut, wantErr := NewCreatorDBUser(db), errCreatorUsernameTaken

		err := sut.CreateUser("", "")

		assert.Equal(t, wantErr.Error(), err.Error())
	})

	t.Run("ErrConnDone", func(t *testing.T) {
		db, mock, def := setup(t)
		defer def(db)
		mock.
			ExpectQuery(`SELECT username FROM users WHERE username = \$1`).
			WillReturnError(sql.ErrConnDone)
		mock.ExpectClose()
		sut, wantErr := NewCreatorDBUser(db), sql.ErrConnDone

		err := sut.CreateUser("", "")

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
