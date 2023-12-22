// Package db contains code to access and work with DynamoDB tables.
package db

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

var (
	// ErrNoItem means that the item was not found in the table.
	ErrNoItem = errors.New("item not found")

	// ErrDupKey means that the item already exists in the table.
	ErrDupKey = errors.New("duplicate key")

	// ErrTooManyItems means that the limit of items has been reached.
	ErrLimitReached = errors.New("too many items")
)

// Retriever defines a type that can retrieve an item from a DynamoDB table.
type Retriever[T any] interface {
	Retrieve(context.Context, string) (T, error)
}

// Inserter defines a type that can insert an item into a DynamoDB table.
type Inserter[T any] interface {
	Insert(context.Context, T) error
}

// Updater defines a type that can update an item in a DynamoDB table.
type Updater[T any] interface {
	Update(context.Context, T) error
}

// Updater defines a type that can update an item in a DynamoDB table using a
// key separate from the T's ID field.
type UpdaterDualKey[T any] interface {
	Update(context.Context, string, T) error
}

// Deleter defines a type that can delete an item from a DynamoDB table.
type Deleter interface {
	Delete(context.Context, string) error
}

// DeleterDualKey defines a type that can delete an item from a DynamoDB table
// using a partition and a sort key.
type DeleterDualKey interface {
	Delete(context.Context, string, string) error
}

// DynamoItemGetter defines a type that can be used to get an item from a
// DynamoDB table. It is used to dependency-inject the DynamoDB client into
// Retrievers.
type DynamoItemGetter interface {
	GetItem(
		context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.GetItemOutput, error)
}

// DynamoQueryer defines a type that can be used to query a DynamoDB table. It
// is used to dependency-inject the DynamoDB client into Retrievers that are
// used to retrieve a collection of items.
type DynamoQueryer interface {
	Query(
		context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options),
	) (*dynamodb.QueryOutput, error)
}

// DynamoItemPutter defines a type that can be used to put an item into a
// DynamoDB table. It is used to dependency-inject the DynamoDB client into
// Inserters and Updaters.
type DynamoItemPutter interface {
	PutItem(
		context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.PutItemOutput, error)
}

// DynamoItemDeleter defines a type that can be used to delete an item from a
// DynamoDB table. It is used to dependency-inject the DynamoDB client into
// Deleters.
type DynamoItemDeleter interface {
	DeleteItem(
		context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options),
	) (*dynamodb.DeleteItemOutput, error)
}

// DynamoTransactWriter defines a type that can be used to write multiple items
// to a DynamoDB table in a transaction. It is used to dependency-inject the
// DynamoDB client into Updaters that are used to update a collecton of items.
type DynamoTransactWriter interface {
	TransactWriteItems(
		context.Context,
		*dynamodb.TransactWriteItemsInput,
		...func(*dynamodb.Options),
	) (*dynamodb.TransactWriteItemsOutput, error)
}

// DynamoItemGetter defines a type that can be used to get and put an item
// from/to a DynamoDB table. It is used to dependency-inject the DynamoDB client
// into some deleters and putters that operate in an item's internal fields.
type DynamoItemGetPutter interface {
	DynamoItemGetter
	DynamoItemPutter
}
