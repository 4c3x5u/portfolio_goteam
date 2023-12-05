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
var userTableName, teamTableName, taskTableName string

func TestMain(m *testing.M) {
	fmt.Println("setting up dynamodb test tables")
	tearDownDynamoDB, err := setUpDynamoDB()
	if err != nil {
		// must tear down here too as some tables might be created while others
		// have failed
		tearDownDynamoDB()
		log.Fatalf("dynamodb setup failed: %s", err)
	}
	defer tearDownDynamoDB()

	fmt.Println("setting up postgres test tables")
	tearDownPostgres, err := setUpPostgres()
	defer tearDownPostgres()
	if err != nil {
		log.Fatalf("postgres setup failed: %s", err)
	}

	m.Run()
}

// setUpDynamoDB sets up the test tables on DynamoDB.
func setUpDynamoDB() (func() error, error) {
	var tearDown func() error

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
	tearDown = tearDownUserTable

	// set up team table
	teamTableName = teamTablePrefix + uuid.New().String()
	tearDownTeamTable, err := setUpTable(svc, &teamTableName)
	if err != nil {
		return tearDown, err
	}
	tearDown = joinTeardowns(tearDown, tearDownTeamTable)

	// set up team table
	taskTableName = taskTablePrefix + uuid.New().String()
	tearDownTaskTable, err := setUpTable(svc, &taskTableName)
	if err != nil {
		return tearDown, err
	}
	tearDown = joinTeardowns(tearDown, tearDownTaskTable)

	// ensure all test tables are created
	if err := allTablesActive(svc); err != nil {
		return tearDown, err
	}

	// populate tables
	_, err = svc.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			userTableName: reqsWriteUser,
			teamTableName: reqsWriteTeam,
		},
	})
	if err != nil {
		return tearDown, err
	}

	// return the teardown function for tables created
	return tearDown, nil
}

// allTablesActive checks whether all tables are created and their status are
// "ACTIVE" every 500 milliseconds until all pass.
func allTablesActive(svc *dynamodb.Client) error {
	fmt.Println("ensuring all test tables are active")
	var userTableActive, teamTableActive, taskTableActive bool
	for {
		time.Sleep(500 * time.Millisecond)

		if !userTableActive {
			resp, err := svc.DescribeTable(
				context.TODO(), &dynamodb.DescribeTableInput{
					TableName: &userTableName,
				},
			)
			if err != nil {
				return err
			}
			if resp.Table.TableStatus == types.TableStatusActive {
				userTableActive = true
			}
		}

		if !teamTableActive {
			resp, err := svc.DescribeTable(
				context.TODO(), &dynamodb.DescribeTableInput{
					TableName: &teamTableName,
				},
			)
			if err != nil {
				return err
			}
			if resp.Table.TableStatus == types.TableStatusActive {
				teamTableActive = true
			}
		}

		if !taskTableActive {
			resp, err := svc.DescribeTable(
				context.TODO(), &dynamodb.DescribeTableInput{
					TableName: &taskTableName,
				},
			)
			if err != nil {
				return err
			}
			if resp.Table.TableStatus == types.TableStatusActive {
				taskTableActive = true
			}
		}

		if userTableActive && teamTableActive && taskTableActive {
			break
		}
	}
	return nil
}

// reqsWriteUser are the requests sent to the user table to initialise it for
// test use.
var reqsWriteUser = []types.WriteRequest{
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: "team1Admin"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: true,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: "team1Member"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: false,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: "team2Admin"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: true,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "66ca0ddf-5f62-4713-bcc9-36cb0954eb7b",
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: "team2Member"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: false,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "66ca0ddf-5f62-4713-bcc9-36cb0954eb7b",
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: "team3Admin"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: true,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "74c80ae5-64f3-4298-a8ff-48f8f920c7d4",
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{Value: "team4Admin"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: true,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "3c3ec4ea-a850-4fc5-aab0-24e9e7223bbc",
		},
	}}},
}

// reqsWriteTeam are the requests sent to the team table to initialise it for
// test use.
var reqsWriteTeam = []types.WriteRequest{
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
		"Members": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "team1Admin"},
				&types.AttributeValueMemberS{Value: "team1Member"},
			},
		},
		"Boards": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ID": &types.AttributeValueMemberS{
							Value: "91536664-9749-4dbb-a470-6e52aa353ae4",
						},
						"Name": &types.AttributeValueMemberS{
							Value: "Team 1 Board 1",
						},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ID": &types.AttributeValueMemberS{
							Value: "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
						},
						"Name": &types.AttributeValueMemberS{
							Value: "Team 1 Board 2",
						},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ID": &types.AttributeValueMemberS{
							Value: "1559a33c-54c5-42c8-8e5f-fe096f7760fa",
						},
						"Name": &types.AttributeValueMemberS{
							Value: "Team 1 Board 3",
						},
					},
				},
			},
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{
			Value: "66ca0ddf-5f62-4713-bcc9-36cb0954eb7b",
		},
		"Members": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "team2Admin"},
				&types.AttributeValueMemberS{Value: "team2Member"},
			},
		},
		"Boards": &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{
			Value: "74c80ae5-64f3-4298-a8ff-48f8f920c7d4",
		},
		"Members": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "team3Admin"},
			},
		},
		"Boards": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ID": &types.AttributeValueMemberS{
							Value: "f0c5d521-ccb5-47cc-ba40-313ddb901165",
						},
						"Name": &types.AttributeValueMemberS{
							Value: "Team 3 Board 1",
						},
					},
				},
			},
		},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"ID": &types.AttributeValueMemberS{
			Value: "3c3ec4ea-a850-4fc5-aab0-24e9e7223bbc",
		},
		"Members": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberS{Value: "team4Admin"},
			},
		},
		"Boards": &types.AttributeValueMemberL{Value: []types.AttributeValue{}},
	}}},
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

// tearDownNothing is returned when there is nothing to tear down.
func tearDownNothing() error { return nil }

// joinTeardowns joins multiple teardowns together into one teardown that
// invokes each of the child teardowns and joins their errors.
func joinTeardowns(tearDowns ...func() error) func() error {
	return func() error {
		var jointErr error
		for _, td := range tearDowns {
			if err := td(); err != nil {
				jointErr = errors.Join(jointErr, err)
			}
		}
		return jointErr
	}
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
