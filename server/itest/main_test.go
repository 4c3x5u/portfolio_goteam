//go:build itest

package itest

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"
	"time"

	"server/log"

	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"

	_ "github.com/lib/pq"
)

// TestMain does setup for and runs the tests in this package.
func TestMain(m *testing.M) {
	logger := log.NewAppLogger()

	// Use a single pool for both containers?
	dbConnStr, tearDownDB := runDBContainer(logger)
	tearDownServer := runServerContainer(dbConnStr, logger)

	code := m.Run()

	tearDownDB()
	tearDownServer()

	os.Exit(code)
}

func runDBContainer(logger log.Logger) (string, func()) {
	logger.Log(log.LevelInfo, "setting up database container...")

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
		func(cfg *docker.HostConfig) {
			cfg.AutoRemove = true
			cfg.RestartPolicy = docker.RestartPolicy{Name: "no"}
		},
	)
	if err != nil {
		logger.Log(log.LevelFatal, "could not start resource: "+err.Error())
		os.Exit(1)
	}

	dbConnStr := "postgres://postgres:postgres@" +
		resource.GetHostPort("5432/tcp") +
		"/goteam?sslmode=disable"

	resource.Expire(300)
	pool.MaxWait = 180 * time.Second
	try := 0
	if err = pool.Retry(func() error {
		try++
		db, sqlErr := sql.Open("postgres", dbConnStr)
		if sqlErr != nil {
			logger.Log(log.LevelInfo, fmt.Sprintf(
				"err for try #%d is %+v", try, sqlErr,
			))
			return sqlErr
		}
		return db.Ping()
	}); err != nil {
		logger.Log(log.LevelFatal, "could not connect to Docker: "+err.Error())
		os.Exit(1)
	}

	return dbConnStr, func() {
		if err := pool.Purge(resource); err != nil {
			logger.Log(log.LevelFatal, "could not purge resource: "+err.Error())
		}
	}
}

func runServerContainer(dbConnStr string, logger log.Logger) func() {
	logger.Log(log.LevelInfo, "setting up server container...")

	pool, err := dockertest.NewPool("")
	if err != nil {
		logger.Log(log.LevelFatal, "could not construct pool: "+err.Error())
		os.Exit(1)
	}

	if err = pool.Client.Ping(); err != nil {
		logger.Log(log.LevelFatal, "could not connect to Docker: "+err.Error())
		os.Exit(1)
	}

	resource, err := pool.BuildAndRunWithOptions(
		"../Dockerfile",
		&dockertest.RunOptions{
			Name: "goteam-server-itest",
			Env: []string{
				"PORT=" + serverPort,
				"DBCONNSTR=" + dbConnStr,
				"JWTKEY=QWERTYQWERTYQWERTYQWERTYQWERTY",
			},
			ExposedPorts: []string{serverPort},
			PortBindings: map[docker.Port][]docker.PortBinding{
				docker.Port(serverPort): {
					{HostIP: serverHost, HostPort: serverPort},
				},
			},
		},
		func(cfg *docker.HostConfig) {
			cfg.AutoRemove = true
			cfg.RestartPolicy = docker.RestartPolicy{Name: "no"}
		})
	if err != nil {
		logger.Log(log.LevelFatal, "could not start resource: "+err.Error())
		os.Exit(1)
	}

	resource.Expire(180)
	pool.MaxWait = 120 * time.Second
	try := 0
	if err = pool.Retry(func() error {
		try++
		if res, errGet := http.Get(serverURL); errGet != nil {
			logger.Log(log.LevelInfo, fmt.Sprintf(
				"err for try #%d is %+v", try, errGet,
			))
			return errGet
		} else if res.StatusCode != 200 {
			logger.Log(log.LevelInfo, fmt.Sprintf(
				"status for try #%d is %s", try, res.Status,
			))
			return errors.New("status: " + res.Status)
		}
		return nil
	}); err != nil {
		logger.Log(log.LevelFatal, "could not connect to Docker: "+err.Error())
		os.Exit(1)
	}

	return func() {
		if err := pool.Purge(resource); err != nil {
			logger.Log(log.LevelFatal, "could not purge resource: "+err.Error())
		}
	}
}
