package db

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// setUpDBMock sets up the DB and Sqlmock for server/db test cases. The
// first two return values are self-explanatory, and the third is the database
// teardown function to be deferred by the test case.
func setUpDBMock(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	return db, mock, func() {
		mock.ExpectClose()
		if err = db.Close(); err != nil {
			t.Fatal(err)
		}
	}
}
