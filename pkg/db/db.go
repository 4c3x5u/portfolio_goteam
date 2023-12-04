// Package db contains code to access and work with DynamoDB tables.
package db

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	// ErrNoItem means that no items were found.
	ErrNoItem = errors.New("item not found")

	// ErrDupKey means that operation does not allow an update on an existing
	// item and the key passed in was not unique.
	ErrDupKey = errors.New("duplicate key")
)

// ItemPutter defines a type that can be used to put an item into a DynamoDB
// table.
type ItemPutter interface {
	PutItem(
		context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)
}

// ItemGetter defines a type that can be used to get an item from a DynamoDB
// table.
type ItemGetter interface {
	GetItem(
		context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.GetItemOutput, error)
}
