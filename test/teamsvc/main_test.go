//go:build itest

package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/test"
)

// tableName is the name of the team table used in the integration tests.
var tableName = "goteam-test-team"

// TestMain sets up the test table in DynamoDB and runs the tests.
func TestMain(m *testing.M) {
	fmt.Println("setting up team table")
	tearDownTables, err := test.SetUpTestTable(
		"TEAM_TABLE_NAME", tableName, writeReqs, "ID", "",
	)
	if err != nil {
		log.Println("set up team table failed:", err)
		return
	}
	defer tearDownTables()

	m.Run()
}

// writeReqs are the requests sent to the test table to initialise it for tests.
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
						"Members": &types.AttributeValueMemberL{
							Value: []types.AttributeValue{
								&types.AttributeValueMemberS{
									Value: "team1Member",
								},
							},
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
						"Members": &types.AttributeValueMemberL{
							Value: []types.AttributeValue{},
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
						"Members": &types.AttributeValueMemberL{
							Value: []types.AttributeValue{
								&types.AttributeValueMemberS{
									Value: "team1Member",
								},
							},
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
