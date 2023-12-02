package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// Putter defines a type that wraps a DynamoDBItemPutter. It can be used
// instead of DynamoDBItemPutter to put an item into a DynamoDB table for
// simplifying DynamoDB usage.
type Putter[T any] interface{ Put(T) error }

// DynamoDBItemPutter defines a type that can be used to put an item into
// a DynamoDB table.
type DynamoDBItemPutter interface {
	PutItem(
		context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)
}

// DynamoDBItemPutter defines a type that can be used to get an item from
// a DynamoDB table.
type DynamoDBItemGetter interface {
	GetItem(
		context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.GetItemOutput, error)
}
