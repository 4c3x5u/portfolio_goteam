//go:build itest

package test

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

var (
	db        *dynamodb.Client
	tableName string
)

// used as a prefix to a uuid when creating test tables
const (
	testTablePrefix = "goteam-test-team-"
)

func TestMain(m *testing.M) {
	fmt.Println("setting up team table")
	tearDownTables, err := setUpTestTable()
	if err != nil {
		log.Println("set up team table failed:", err)
		return
	}
	defer tearDownTables()

	m.Run()
}

// setUpTestTable sets up the test tables on DynamoDB.
func setUpTestTable() (func() error, error) {
	// create dynamodb client
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return test.TearDownNone, err
	}
	db = dynamodb.NewFromConfig(cfg)

	// set up team table
	tableName = testTablePrefix + uuid.New().String()
	tearDownTable, err := test.CreateTable(db, &tableName, "ID", "")
	if err != nil {
		return tearDownTable, err
	}

	// set environvar for team putter & getter to read the table name from
	if err := os.Setenv("TEAM_TABLE_NAME", tableName); err != nil {
		if err != nil {
			return tearDownTable, err
		}
	}

	// ensure all test tables are created
	if err := test.EnsureTableActive(db, tableName); err != nil {
		return tearDownTable, err
	}

	// populate tables
	_, err = db.BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: writeReqs,
		},
	})
	if err != nil {
		return tearDownTable, err
	}

	// return the teardown function for tables created
	return tearDownTable, nil
}

// writeReqs are the requests sent to the team table to initialise it for test
// use.
var writeReqs = []types.WriteRequest{
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
				&types.AttributeValueMemberS{Value: "team4Member"},
			},
		},
		"Boards": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ID": &types.AttributeValueMemberS{
							Value: "ca47fbec-269e-4ef4-a74a-bcfbcd599fd5",
						},
						"Name": &types.AttributeValueMemberS{
							Value: "Team 4 Board 1",
						},
					},
				},
			},
		},
	}}},
}
