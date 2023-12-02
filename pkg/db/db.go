package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// DynamoDBItemPutter defines a type that can be used to put an item into
// a DynamoDB table.
type DynamoDBItemPutter interface {
	PutItem(
		context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)
}

// Putter defines a type that wraps a DynamoDBItemPutter. It can be used
// instead of DynamoDBItemPutter to simplify DynamoDB usage within the app.
type Putter[T any] interface{ Put(T) error }
