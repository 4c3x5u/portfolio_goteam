//go:build itest

package main

import (
	"database/sql"
	"os"
	"testing"
	"time"

	"server/log"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	_ "github.com/lib/pq"
)

// dbConnStr is the connection string for the database used in integration tests.
var dbConnStr string

// TestMain does setup for and runs the tests in this package.
func TestMain(m *testing.M) {
	logger := log.NewAppLogger()

	pool, err := dockertest.NewPool("")
	if err != nil {
		logger.Log(log.LevelFatal, "could not construct pool: "+err.Error())
		os.Exit(1)
	}

	if err = pool.Client.Ping(); err != nil {
		logger.Log(log.LevelFatal, "could not connect to Docker: "+err.Error())
		os.Exit(1)
	}

	resource, err := pool.RunWithOptions(
		&dockertest.RunOptions{
			Repository: "postgres",
			Tag:        "14",
			Env: []string{
				"POSTGRES_USER=postgres",
				"POSTGRES_PASSWORD=postgres",
				"POSTGRES_DB=goteam",
				"listen_addresses = '*'",
			},
		},
		func(config *docker.HostConfig) {
			config.AutoRemove = true
			config.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		logger.Log(log.LevelFatal, "could not start resource: "+err.Error())
		os.Exit(1)
	}

	hostPort := resource.GetHostPort("5432/tcp")
	dbConnStr = "postgres://postgres:postgres@" +
		hostPort + "/dbname?sslmode=disable"

	resource.Expire(120)

	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, sqlErr := sql.Open("postgres", dbConnStr)
		if sqlErr != nil {
			return sqlErr
		}
		return db.Ping()
	}); err != nil {
		logger.Log(log.LevelFatal, "could not connect to Docker: "+err.Error())
		os.Exit(1)
	}

	code := m.Run()

	if err = pool.Purge(resource); err != nil {
		logger.Log(log.LevelFatal, "could not purge resource: "+err.Error())
	}

	os.Exit(code)
}
