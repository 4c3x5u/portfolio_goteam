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

// used as a prefix to a uuid when creating test tables
const testTablePrefix = "goteam-test-user-"

var (
	db        *dynamodb.Client
	tableName string
)

func TestMain(m *testing.M) {
	fmt.Println("setting up user table")
	tearDownTables, err := setUpTables()
	defer tearDownTables()
	if err != nil {
		log.Println("set up user table failed:", err)
		return
	}

	m.Run()
}

// setUpTables sets up the test tables on DynamoDB.
func setUpTables() (func() error, error) {
	// create dynamodb client
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		return test.TearDownNone, err
	}
	db = dynamodb.NewFromConfig(cfg)

	// set up user table
	tableName = testTablePrefix + uuid.New().String()
	tearDown, err := test.CreateTable(
		db, &tableName, "Username", "",
	)
	if err != nil {
		return test.TearDownNone, err
	}

	// set environvar for user putter & getter to read the table name from
	if err := os.Setenv("USER_TABLE_NAME", tableName); err != nil {
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

// writeReqs are the requests sent to the user table to initialise it for test
// use.
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
