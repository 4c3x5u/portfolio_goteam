//go:build itest

package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/google/uuid"

	"github.com/kxplxn/goteam/test"
)

// tableName is the name of the user table used in the integration tests.
var tableName = "goteam-test-user-" + uuid.New().String()

// TestMain sets up the test table in DynamoDB and runs the tests.
func TestMain(m *testing.M) {
	fmt.Println("setting up user table")
	tearDownTables, err := test.SetUpTestTable(
		"USER_TABLE_NAME", tableName, writeReqs, "Username", "",
	)
	defer tearDownTables()
	if err != nil {
		log.Println("set up user table failed:", err)
		return
	}

	m.Run()
}

// writeReqs are the requests sent to the test table to initialise it for tests.
var writeReqs = []types.WriteRequest{
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"Username": &types.AttributeValueMemberS{Value: "team1Admin"},
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
		"Username": &types.AttributeValueMemberS{Value: "team1Member"},
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
		"Username": &types.AttributeValueMemberS{Value: "team2Admin"},
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
		"Username": &types.AttributeValueMemberS{Value: "team2Member"},
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
		"Username": &types.AttributeValueMemberS{Value: "team3Admin"},
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
		"Username": &types.AttributeValueMemberS{Value: "team4Admin"},
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
	{PutRequest: &types.PutRequest{Item: map[string]types.AttributeValue{
		"Username": &types.AttributeValueMemberS{Value: "team4Member"},
		"Password": &types.AttributeValueMemberB{
			Value: []byte(
				"$2a$11$kZfdRfTOjhfmel7J4WRG3eltzH9lavxp5qyrpFnzc9MIYLhZNCqTO",
			),
		},
		"IsAdmin": &types.AttributeValueMemberBOOL{
			Value: false,
		},
		"TeamID": &types.AttributeValueMemberS{
			Value: "3c3ec4ea-a850-4fc5-aab0-24e9e7223bbc",
		},
	}}},
}
