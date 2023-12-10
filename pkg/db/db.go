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

// Retriever defines a type that can retrieve an item from a DynamoDB table by
// a string value.
type Retriever[T any] interface {
	Retrieve(context.Context, string) (T, error)
}

// Creator defines a type that can insert a new item into a DynamoDB table.
type Inserter[T any] interface {
	Insert(context.Context, T) error
}

// Updater defines a type that can update an item in a DynamoDB table.
type Updater[T any] interface {
	Update(context.Context, T) error
}

// Deleter defines a type that can delete an item from a DynamoDB table.
type Deleter interface {
	Delete(context.Context, string) error
}

// DynamoItemGetter defines a type that can be used to get an item from a
// DynamoDB table. It is used to dependency-inject the DynamoDB client into
// Retriever.
type DynamoItemGetter interface {
	GetItem(
		context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.GetItemOutput, error)
}

// DynamoItemPutter defines a type that can be used to put an item into a
// DynamoDB table. It is used to dependency-inject the DynamoDB client into
// Inserter and Updater.
type DynamoItemPutter interface {
	PutItem(
		context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)
}

// DynamoItemDeleter defines a type that can be used to delete an item from a
// DynamoDB table. It is used to dependency-inject the DynamoDB client into
// Deleter.
type DynamoItemDeleter interface {
	DeleteItem(
		context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.DeleteItemOutput, error)
}
