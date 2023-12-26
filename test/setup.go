//go:build itest

package test

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

// AddAuthCookie is used in various test cases to authenticate the request being
// sent to a handler. It returns a higher-order function that adds the given
// token as the auth cookie value to the request so that it can be initialised
// with a token in table-driven tests before the test is actually run.
func AddAuthCookie(token string) func(*http.Request) {
	return func(r *http.Request) {
		r.AddCookie(&http.Cookie{Name: "auth-token", Value: token})
	}
}

// AddStateCookie adds the given token as the state cookie value to the request.
// It returns a higher-order function that adds the given token as the state
// cookie value to the request so that it can be initialised with a token in
// table-driven tests before the test is actually run.
func AddStateCookie(token string) func(*http.Request) {
	return func(r *http.Request) {
		r.AddCookie(&http.Cookie{Name: "state-token", Value: token})
	}
}

// DB returns the DynamoDB client used in integration tests. If the client has
// not yet been created, it is created and returned.
func DB() *dynamodb.Client {
	if db == nil {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			fmt.Println("error loading default config")
			return nil
		}
		db = dynamodb.NewFromConfig(cfg)
	}
	return db
}

// SetUpTestTable sets up a test table in DynamoDB.
func SetUpTestTable(
	envVar string,
	tableName string,
	writeReqs []types.WriteRequest,
	partKey string,
	sortKey string,
	secINames ...string,
) (func() error, error) {
	// set up test table
	tearDown, err := createTable(
		DB(), &tableName, partKey, sortKey, secINames...,
	)
	if err != nil {
		return tearDownNone, err
	}

	// set environvar for putters/getter to read the table name from
	if err := os.Setenv(envVar, tableName); err != nil {
		return tearDown, err
	}

	// ensure test table is created and active
	if err := ensureTableActive(db, tableName); err != nil {
		return tearDown, err
	}

	// populate test table with given write requests
	_, err = DB().BatchWriteItem(context.TODO(), &dynamodb.BatchWriteItemInput{
		RequestItems: map[string][]types.WriteRequest{
			tableName: writeReqs,
		},
	})
	if err != nil {
		return tearDown, err
	}

	// return the teardown function
	return tearDown, nil
}

// createTable creates a DynamoDB table with the given name, and given sort and
// partition keys, and secondary index names.
func createTable(
	svc *dynamodb.Client,
	name *string,
	partKey string,
	sortKey string,
	secINames ...string,
) (func() error, error) {
	fmt.Println("creating", name, "table")

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
		return tearDownNone, err
	}

	// create user table teardown function
	return func() error {
		svc.DeleteTable(context.TODO(), &dynamodb.DeleteTableInput{
			TableName: name,
		})
		return nil
	}, nil
}

// ensureTableActive checks whether the test table is created and its status is
// "ACTIVE" every 500 milliseconds until it is true.
func ensureTableActive(svc *dynamodb.Client, tableName string) error {
	fmt.Println("ensuring all test tables are active")
	var teamTableActive bool
	for {
		if !teamTableActive {
			resp, err := svc.DescribeTable(
				context.TODO(), &dynamodb.DescribeTableInput{
					TableName: &tableName,
				},
			)
			if err != nil {
				return err
			}
			if resp.Table.TableStatus == types.TableStatusActive {
				teamTableActive = true
			}
		} else {
			break
		}

		time.Sleep(500 * time.Millisecond)
	}
	return nil
}

// tearDownNone is returned when there is nothing to tear down. This is done so
// the tear-down function can be called before checking for error because there
// are cases when a non-empty tear-down function is returned with an error.
func tearDownNone() error { return nil }
