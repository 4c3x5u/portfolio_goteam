//go:build itest

package api

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	_ "github.com/lib/pq"
)

const (
	userTablePrefix = "goteam-test-user-"
	teamTablePrefix = "goteam-test-user-"
	taskTablePrefix = "goteam-test-user-"
)

func TestMain(m *testing.M) {
	tearDownDynamoDB, err := setUpDynamoDB()
	if err != nil {
		log.Fatalf("dynamodb setup failed: %s", err)
	}

	// Create and run the docker container for itest database.
	tearDownPostgres, err := setUpPostgres()
	if err != nil {
		log.Fatalf("postgres setup failed: %s", err)
	}

	// Run integration tests.
	code := m.Run()

	if err := tearDownDynamoDB(); err != nil {
		log.Fatalf("dynamodb teardown failed: %s", err)
	}

	// Tear down the database container.
	if err := tearDownPostgres(); err != nil {
		log.Fatalf("postgres teardown failed: %s", err)
	}

	// Done.
	os.Exit(code)
}

func setUpDynamoDB() (func() error, error) {
	return func() error { return nil }, nil
}

// TODO: remove once fully migrated to DynamoDB
func setUpPostgres() (func() error, error) {
	emptyTeardown := func() error { return nil }

	pool, err := dockertest.NewPool("")
	if err != nil {
		return emptyTeardown, fmt.Errorf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		return emptyTeardown, fmt.Errorf("Could not connect to Docker: %s", err)
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
		return emptyTeardown, fmt.Errorf("Could not start resource: %s", err)
	}
	if err := resource.Expire(180); err != nil {
		return emptyTeardown, fmt.Errorf("expire error: %s", err)
	}

	// Get the connection string to the database.
	databaseURL := "postgres://itestdb_usr:itestdb_pwd@" +
		resource.GetHostPort("5432/tcp") + "/itestdb?sslmode=disable"
	log.Println("Connecting to database on url: ", databaseURL)

	// Make sure the container and the database are healthy.
	// IMPORTANT: if it's the first time creating the image, set the maxWait to
	// something higher (e.g. 180 seconds).
	pool.MaxWait = 15 * time.Second
	if err = pool.Retry(func() error {
		db, err = sql.Open("postgres", databaseURL)
		if err != nil {
			return err
		}
		return db.Ping()
	}); err != nil {
		return emptyTeardown, fmt.Errorf("Could not connect to docker: %s", err)
	}

	// Initialise the database with schema and tables.
	qInitBytes, err := os.ReadFile("init.sql")
	if err != nil {
		log.Fatal("+++", err)
	}
	if _, err = db.Exec(string(qInitBytes)); err != nil {
		log.Fatal("+++", err)
	}

	return func() error { return pool.Purge(resource) }, nil
}
