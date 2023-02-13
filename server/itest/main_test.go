//go:build itest

package itest

import (
	"database/sql"
	"io/ioutil"
	"os"
	"server/log"
	"testing"
)

func TestMain(m *testing.M) {
	logger := log.NewAppLogger()

	// Initialise the database with schema/tables.
	db, err := sql.Open(
		"postgres",
		"postgres://itestuser:itestpwd@localhost:5432/itestdb?sslmode=disable",
	)
	if err != nil {
		logger.Log(log.LevelFatal, err.Error())
		os.Exit(1)
	}
	qInitBytes, err := ioutil.ReadFile("init.sql")
	if err != nil {
		logger.Log(
			log.LevelFatal,
			"could not read db init script: "+err.Error(),
		)
		os.Exit(1)
	}
	if _, err := db.Exec(string(qInitBytes)); err != nil {
		logger.Log(
			log.LevelFatal,
			"could not execute db init script: "+err.Error(),
		)
		os.Exit(1)
	}

	os.Exit(m.Run())
}
