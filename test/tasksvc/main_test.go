//go:build itest

package tasksvc

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	"github.com/kxplxn/goteam/test"
)

// used as a prefix to a uuid when creating test table
const testTablePrefix = "goteam-test-task-"

var (
	db        *dynamodb.Client
	tableName string
)

// TestMain sets up the test tables in DynamoDB and runs the tests.
func TestMain(m *testing.M) {
	fmt.Println("setting up task table")
	tearDown, err := setUpTestTable()
	defer tearDown()
	if err != nil {
		log.Println("set up task failed:", err)
		return
	}

	m.Run()
}

// setUpTestTable sets up the test table in DynamoDB.
func setUpTestTable() (func() error, error) {
	// create dynamodb client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return nil, err
	}
	db = dynamodb.NewFromConfig(cfg)

	// set up team table
	tableName = testTablePrefix + uuid.New().String()
	tearDown, err := test.CreateTable(
		db, &tableName, "TeamID", "ID", "BoardID",
	)
	if err != nil {
		return nil, err
	}

	// set environvar for task putter & getter to read the table name from
	if err := os.Setenv("TASK_TABLE_NAME", tableName); err != nil {
		if err != nil {
			return tearDown, err
		}
	}

	// ensure all test tables are created
	if err := test.EnsureTableActive(db, tableName); err != nil {
		return tearDown, err
	}

	// populate tables
	_, err = db.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: writeReqs,
		},
	})
	if err != nil {
		return tearDown, err
	}

	// return the teardown function for tables created
	return tearDown, nil
}

// reqsWriteTask are the requests sent to the task table to initialise it for
// test use.
var writeReqs = []types.WriteRequest{
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "74c80ae5-64f3-4298-a8ff-48f8f920c7d4",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "c146486d-7260-4d3d-9da5-2545a5109ca1",
		},
		"Title": &types.AttributeValueMemberS{Value: "task 1"},
		"Order": &types.AttributeValueMemberN{Value: "1"},
		"Subtasks": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Title": &types.AttributeValueMemberS{
							Value: "subtask 1",
						},
						"IsDone": &types.AttributeValueMemberBOOL{Value: false},
					},
				},
			},
		},
		"BoardID": &types.AttributeValueMemberS{
			Value: "f0c5d521-ccb5-47cc-ba40-313ddb901165",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "0"},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "74c80ae5-64f3-4298-a8ff-48f8f920c7d4",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "379a94ac-3af4-4ca0-8469-5b41567e1bf1",
		},
		"Title": &types.AttributeValueMemberS{Value: "task 2"},
		"Order": &types.AttributeValueMemberN{Value: "1"},
		"Subtasks": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Title": &types.AttributeValueMemberS{
							Value: "subtask 2",
						},
						"IsDone": &types.AttributeValueMemberBOOL{Value: false},
					},
				},
			},
		},
		"BoardID": &types.AttributeValueMemberS{
			Value: "f0c5d521-ccb5-47cc-ba40-313ddb901165",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "1"},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "74c80ae5-64f3-4298-a8ff-48f8f920c7d4",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "b59bcff3-9829-4630-a21f-83977dfc4665",
		},
		"Title": &types.AttributeValueMemberS{Value: "task 3"},
		"Order": &types.AttributeValueMemberN{Value: "1"},
		"Subtasks": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Title": &types.AttributeValueMemberS{
							Value: "subtask 3",
						},
						"IsDone": &types.AttributeValueMemberBOOL{Value: false},
					},
				},
			},
		},
		"BoardID": &types.AttributeValueMemberS{
			Value: "f0c5d521-ccb5-47cc-ba40-313ddb901165",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "2"},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "74c80ae5-64f3-4298-a8ff-48f8f920c7d4",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "8fd4d2a3-6247-4dcc-bc6a-5077d8e57be1",
		},
		"Title": &types.AttributeValueMemberS{Value: "task 4"},
		"Order": &types.AttributeValueMemberN{Value: "1"},
		"Subtasks": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Title": &types.AttributeValueMemberS{
							Value: "subtask 4",
						},
						"IsDone": &types.AttributeValueMemberBOOL{Value: false},
					},
				},
			},
		},
		"BoardID": &types.AttributeValueMemberS{
			Value: "f0c5d521-ccb5-47cc-ba40-313ddb901165",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "3"},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "c684a6a0-404d-46fa-9fa5-1497f9874567",
		},
		"Title": &types.AttributeValueMemberS{Value: "task 5"},
		"Order": &types.AttributeValueMemberN{Value: "1"},
		"BoardID": &types.AttributeValueMemberS{
			Value: "91536664-9749-4dbb-a470-6e52aa353ae4",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "0"},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "8fb040a2-910c-47af-a4ab-9dee49f16d1d",
		},
		"Title": &types.AttributeValueMemberS{Value: "task 6"},
		"Order": &types.AttributeValueMemberN{Value: "1"},
		"BoardID": &types.AttributeValueMemberS{
			Value: "1559a33c-54c5-42c8-8e5f-fe096f7760fa",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "2"},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "a2e5b55f-01cc-4eac-8882-d76acb94a5b9",
		},
		"Title": &types.AttributeValueMemberS{Value: "task 7"},
		"Order": &types.AttributeValueMemberN{Value: "2"},
		"BoardID": &types.AttributeValueMemberS{
			Value: "1559a33c-54c5-42c8-8e5f-fe096f7760fa",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "2"},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "e0021a56-6a1e-4007-b773-395d3991fb7e",
		},
		"Title": &types.AttributeValueMemberS{Value: "task 8"},
		"Order": &types.AttributeValueMemberN{Value: "3"},
		"Subtasks": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Title": &types.AttributeValueMemberS{
							Value: "subtask 5",
						},
						"IsDone": &types.AttributeValueMemberBOOL{Value: false},
					},
				},
			},
		},
		"BoardID": &types.AttributeValueMemberS{
			Value: "1559a33c-54c5-42c8-8e5f-fe096f7760fa",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "2"},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "9362dcd5-408b-4e26-9dda-68056ba7b833",
		},
		"Title": &types.AttributeValueMemberS{Value: "task 9"},
		"Order": &types.AttributeValueMemberN{Value: "1"},
		"Subtasks": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Title": &types.AttributeValueMemberS{
							Value: "subtask 6",
						},
						"IsDone": &types.AttributeValueMemberBOOL{Value: false},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Title": &types.AttributeValueMemberS{
							Value: "subtask 7",
						},
						"IsDone": &types.AttributeValueMemberBOOL{Value: false},
					},
				},
			},
		},
		"BoardID": &types.AttributeValueMemberS{
			Value: "1559a33c-54c5-42c8-8e5f-fe096f7760fa",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "2"},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "01a3168d-6d2a-46fb-aed9-70c26a4d71e9",
		},
		"Title":       &types.AttributeValueMemberS{Value: "task 10"},
		"Description": &types.AttributeValueMemberS{Value: "some description"},
		"Order":       &types.AttributeValueMemberN{Value: "1"},
		"Subtasks": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Title": &types.AttributeValueMemberS{
							Value: "subtask 8",
						},
						"IsDone": &types.AttributeValueMemberBOOL{Value: false},
					},
				},
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Title": &types.AttributeValueMemberS{
							Value: "subtask 9",
						},
						"IsDone": &types.AttributeValueMemberBOOL{Value: true},
					},
				},
			},
		},
		"BoardID": &types.AttributeValueMemberS{
			Value: "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "0"},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "9dd9c982-8d1c-49ac-a412-3b01ba74b634",
		},
		"Title": &types.AttributeValueMemberS{Value: "task 11"},
		"Order": &types.AttributeValueMemberN{Value: "1"},
		"BoardID": &types.AttributeValueMemberS{
			Value: "fdb82637-f6a5-4d55-9dc3-9f60061e632f",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "2"},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "3c3ec4ea-a850-4fc5-aab0-24e9e7223bbc",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "55e275e4-de80-4241-b73b-88e784d5522b",
		},
		"Title": &types.AttributeValueMemberS{Value: "team 4 task 1"},
		"Description": &types.AttributeValueMemberS{
			Value: "team 4 task 1 description",
		},
		"Order": &types.AttributeValueMemberN{Value: "1"},
		"Subtasks": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Title": &types.AttributeValueMemberS{
							Value: "team 4 subtask 1",
						},
						"IsDone": &types.AttributeValueMemberBOOL{Value: false},
					},
				},
			},
		},
		"BoardID": &types.AttributeValueMemberS{
			Value: "ca47fbec-269e-4ef4-a74a-bcfbcd599fd5",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "0"},
	}}},
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"TeamID": &types.AttributeValueMemberS{
			Value: "3c3ec4ea-a850-4fc5-aab0-24e9e7223bbc",
		},
		"ID": &types.AttributeValueMemberS{
			Value: "5ccd750d-3783-4832-891d-025f24a4944f",
		},
		"Title": &types.AttributeValueMemberS{Value: "team 4 task 2"},
		"Description": &types.AttributeValueMemberS{
			Value: "team 4 task 2 description",
		},
		"Order": &types.AttributeValueMemberN{Value: "0"},
		"Subtasks": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"Title": &types.AttributeValueMemberS{
							Value: "team 4 subtask 2",
						},
						"IsDone": &types.AttributeValueMemberBOOL{Value: true},
					},
				},
			},
		},
		"BoardID": &types.AttributeValueMemberS{
			Value: "ca47fbec-269e-4ef4-a74a-bcfbcd599fd5",
		},
		"ColumnNumber": &types.AttributeValueMemberN{Value: "0"},
	}}},
}
