//go:build itest

package itest

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestMain(m *testing.M) {
	// Initialise the database.
	db, err := sql.Open("postgres", dbConnStr)
	if err != nil {
		log.Fatal(err)
	}
	qInitBytes, err := ioutil.ReadFile("init.sql")
	if err != nil {
		log.Fatal(err)
	}
	if _, err := db.Exec(string(qInitBytes)); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}
