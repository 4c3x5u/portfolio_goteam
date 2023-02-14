//go:build itest

package itest

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	_ "github.com/lib/pq"
)

// dbConnPool is the database connection pool used during integration testing.
// It is set in main_test.go/TestMain.
var dbConnPool *sql.DB

func TestMain(m *testing.M) {
	// create a new pool to run the resource (i.e. test db) in
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pulls the db image, creates a container based on it, and run it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14",
		Env: []string{
			"POSTGRES_USER=itestdb_usr",
			"POSTGRES_PASSWORD=itestdb_pwd",
			"POSTGRES_DB=itestdb",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	databaseURL := "postgres://itestdb_usr:itestdb_pwd@" +
		resource.GetHostPort("5432/tcp") + "/itestdb?sslmode=disable"

	log.Println("Connecting to database on url: ", databaseURL)

	// Tell docker to hard kill the container in 180 seconds
	resource.Expire(180)

	// exponential backoff-retry, because the application in the container might
	// not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		dbConnPool, err = sql.Open("postgres", databaseURL)
		if err != nil {
			return err
		}
		return dbConnPool.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// initialise the database with schema and tables
	qInitBytes, err := ioutil.ReadFile("init.sql")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = dbConnPool.Exec(string(qInitBytes)); err != nil {
		log.Fatal(err)
	}

	//Run tests
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}
