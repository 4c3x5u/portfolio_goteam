//go:build itest

package api

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	_ "github.com/lib/pq"
)

// used as a prefix to a uuid when creating test tables
const (
	userTablePrefix = "goteam-test-user-"
	teamTablePrefix = "goteam-test-team-"
	taskTablePrefix = "goteam-test-task-"
)

// set during DynamoDB setup to be used by tests
var (
	userTableName string
	teamTableName string
	taskTableName string
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

// tearDownNothing is returned when there is nothing to tear down.
func tearDownNothing() error { return nil }

// setUpDynamoDB sets up the test tables on DynamoDB.
func setUpDynamoDB() (func() error, error) {
	// create dynamodb client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return tearDownNothing, err
	}
	svc := dynamodb.NewFromConfig(cfg)

	// set up user table
	userTableName = userTablePrefix + uuid.New().String()
	tearDownUserTable, err := setUpTable(svc, &userTableName)
	if err != nil {
		return tearDownNothing, err
	}

	// set up team table
	teamTableName = teamTablePrefix + uuid.New().String()
	tearDownTeamTable, err := setUpTable(svc, &teamTableName)
	if err != nil {
		return tearDownUserTable, err
	}

	// set up team table
	taskTableName = taskTablePrefix + uuid.New().String()
	tearDownTaskTable, err := setUpTable(svc, &taskTableName)
	if err != nil {
		return func() error {
			var errs error
			if err = tearDownUserTable(); err != nil {
				errs = err
			}
			if err = tearDownTeamTable(); err != nil {
				return errors.Join(errs, err)
			}
			return nil
		}, err
	}

	// return the teardown function for tables created
	return func() error {
		var errs error
		if err = tearDownUserTable(); err != nil {
			errs = err
		}
		if err = tearDownTeamTable(); err != nil {
			errs = errors.Join(errs, err)
		}
		if err = tearDownTaskTable(); err != nil {
			errs = errors.Join(errs, err)
		}
		return err
	}, nil
}

// setUpTable sets up a DynamoDB table with the given name and a string
// partition key named ID.
func setUpTable(svc *dynamodb.Client, name *string) (func() error, error) {
	_, err := svc.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName: name,
		AttributeDefinitions: []types.AttributeDefinition{
			{AttributeName: aws.String("ID"), AttributeType: "S"},
		},
		KeySchema: []types.KeySchemaElement{
			{AttributeName: aws.String("ID"), KeyType: "HASH"},
		},
		BillingMode: types.BillingModeProvisioned,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(25),
			WriteCapacityUnits: aws.Int64(25),
		},
	})
	if err != nil {
		return tearDownNothing, err
	}

	// create user table teardown function
	return func() error {
		svc.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
			TableName: name,
		})
		return nil
	}, nil
}

// TODO: remove once fully migrated to DynamoDB
func setUpPostgres() (func() error, error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return tearDownNothing, fmt.Errorf("Could not construct pool: %s", err)
	}
	err = pool.Client.Ping()
	if err != nil {
		return tearDownNothing, fmt.Errorf("Could not connect to Docker: %s", err)
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
		return tearDownNothing, fmt.Errorf("Could not start resource: %s", err)
	}
	if err := resource.Expire(180); err != nil {
		return tearDownNothing, fmt.Errorf("expire error: %s", err)
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
		return tearDownNothing, fmt.Errorf("Could not connect to docker: %s", err)
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
