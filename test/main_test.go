//go:build itest

package api

import (
	"context"
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
)

var db *dynamodb.Client

// used as a prefix to a uuid when creating test tables
const (
	userTablePrefix = "goteam-test-user-"
	teamTablePrefix = "goteam-test-team-"
	taskTablePrefix = "goteam-test-task-"
)

// set during DynamoDB setup to be used by tests
var userTableName, teamTableName, taskTableName string

func TestMain(m *testing.M) {
	fmt.Println("setting up test tables")
	tearDownTables, err := setUpTables()
	defer tearDownTables()
	if err != nil {
		log.Printf("test tables setup failed: %s", err)
		return
	}

	m.Run()
}

// setUpTables sets up the test tables on DynamoDB.
func setUpTables() (func() error, error) {
	var tearDown func() error

	// create dynamodb client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return tearDownNothing, err
	}
	db = dynamodb.NewFromConfig(cfg)

	// set up user table
	userTableName = userTablePrefix + uuid.New().String()
	tearDownUserTable, err := createTable(
		db, &userTableName, "Username", "",
	)
	if err != nil {
		return tearDownNothing, err
	}
	tearDown = tearDownUserTable

	// set environvar for user putter & getter to read the table name from
	if err := os.Setenv("DYNAMODB_TABLE_USER", userTableName); err != nil {
		if err != nil {
			return tearDown, err
		}
	}

	// set up team table
	teamTableName = teamTablePrefix + uuid.New().String()
	tearDownTeamTable, err := createTable(db, &teamTableName, "ID", "")
	if err != nil {
		return tearDown, err
	}
	tearDown = joinTeardowns(tearDown, tearDownTeamTable)

	// set environvar for team putter & getter to read the table name from
	if err := os.Setenv("DYNAMODB_TABLE_TEAM", teamTableName); err != nil {
		if err != nil {
			return tearDown, err
		}
	}

	// set up team table
	taskTableName = taskTablePrefix + uuid.New().String()
	tearDownTaskTable, err := createTable(
		db, &taskTableName, "TeamID", "ID", "BoardID",
	)
	if err != nil {
		return tearDown, err
	}
	tearDown = joinTeardowns(tearDown, tearDownTaskTable)

	// set environvar for task putter & getter to read the table name from
	if err := os.Setenv("DYNAMODB_TABLE_TASK", taskTableName); err != nil {
		if err != nil {
			return tearDown, err
		}
	}

	// ensure all test tables are created
	if err := allTablesActive(db); err != nil {
		return tearDown, err
	}

	// populate tables
	_, err = db.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			userTableName: reqsWriteUser,
			teamTableName: reqsWriteTeam,
			taskTableName: reqsWriteTask,
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

// createTable creates a DynamoDB table with the given name and a string
// partition key named ID.
func createTable(
	svc *dynamodb.Client, name *string, partKey string, sortKey string, secINames ...string,
) (func() error, error) {
	attrDefs := []types.AttributeDefinition{
		{AttributeName: &partKey, AttributeType: types.ScalarAttributeTypeS},
	}

	var secIs []types.GlobalSecondaryIndex
	for _, iname := range secINames {
		attrDefs = append(attrDefs, types.AttributeDefinition{
			AttributeName: &iname, AttributeType: types.ScalarAttributeTypeS,
		})

		secIs = append(secIs, types.GlobalSecondaryIndex{
			IndexName: aws.String(iname + "_index"),
			KeySchema: []types.KeySchemaElement{
				{AttributeName: &iname, KeyType: types.KeyTypeHash},
				{AttributeName: &partKey, KeyType: types.KeyTypeRange},
			},
			Projection: &types.Projection{
				ProjectionType: types.ProjectionTypeAll,
			},
			ProvisionedThroughput: &types.ProvisionedThroughput{
				ReadCapacityUnits:  aws.Int64(25),
				WriteCapacityUnits: aws.Int64(25),
			},
		})
	}

	keySchema := []types.KeySchemaElement{
		{AttributeName: &partKey, KeyType: types.KeyTypeHash},
	}
	if sortKey != "" {
		attrDefs = append(attrDefs, types.AttributeDefinition{
			AttributeName: &sortKey, AttributeType: types.ScalarAttributeTypeS,
		})
		keySchema = append(keySchema, types.KeySchemaElement{
			AttributeName: &sortKey, KeyType: types.KeyTypeRange,
		})
	}

	_, err := svc.CreateTable(context.TODO(), &dynamodb.CreateTableInput{
		TableName:            name,
		AttributeDefinitions: attrDefs,
		KeySchema:            keySchema,
		BillingMode:          types.BillingModeProvisioned,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(25),
			WriteCapacityUnits: aws.Int64(25),
		},
		GlobalSecondaryIndexes: secIs,
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
