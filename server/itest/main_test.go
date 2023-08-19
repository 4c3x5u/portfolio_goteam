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

func TestMain(m *testing.M) {
	// Create and run the docker container for itest database.
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}
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
	resource.Expire(180)

	// Get the connection string to the database.
	databaseURL := "postgres://itestdb_usr:itestdb_pwd@" +
		resource.GetHostPort("5432/tcp") + "/itestdb?sslmode=disable"
	log.Println("Connecting to database on url: ", databaseURL)

	// Make sure the container and the database are healthy.
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Initialise the database with schema and tables.
	qInitBytes, err := ioutil.ReadFile("init.sql")
	if err != nil {
		log.Fatal("+++", err)
	}
	if _, err = db.Exec(string(qInitBytes)); err != nil {
		log.Fatal("+++", err)
	}

	// Run integration tests.
	code := m.Run()

	// Tear down the database container.
	if err = pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	// Done.
	os.Exit(code)
}
