package dbaccess

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
)

// SetUpDBTest sets up the *sql.DB and *sqlmock.Sqlmock for server/db test
// cases.
func SetUpDBTest(t *testing.T) (*sql.DB, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	teardown := func() {
		mock.ExpectClose()
		if err = db.Close(); err != nil {
			t.Fatal(err)
		}
	}
	return db, mock, teardown
}
