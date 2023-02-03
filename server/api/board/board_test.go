//go:build itest

package board

import (
	"database/sql"
	"os"
	"testing"
)

// TestBoardAPIRoute tests the board route served by the server to assert that
// it
func TestBoardAPIRoute(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("DBCONNSTR"))
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Ping(); err != nil {
		t.Fatal(err)
	}
}
