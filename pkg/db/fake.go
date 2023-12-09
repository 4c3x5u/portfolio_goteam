//go:build utest

package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// FakeRetriever is a test fake for Retriever.
type FakeRetriever[T any] struct {
	Item T
	Err  error
}

// Retrieve discards params and returns FakeRetriever.Item and
// FakeRetriever.Err.
func (f *FakeRetriever[T]) Retrieve(
	ctx context.Context, username string,
) (T, error) {
	return f.Item, f.Err
}

// FakeInserter is a test fake for Inserter.
type FakeInserter[T any] struct{ Err error }

// Insert discards params and returns FakeInserter.Err.
func (f *FakeInserter[T]) Insert(_ context.Context, _ T) error { return f.Err }

// FakeUpdater is a test fake for Updater.
type FakeUpdater[T any] struct{ Err error }

// Update discards params and returns FakeUpdater.Err.
func (f *FakeUpdater[T]) Update(_ context.Context, _ T) error { return f.Err }

// FakeDynamoDBGetter is a test fake for ItemGetter.
type FakeDynamoDBGetter struct {
	Out *dynamodb.GetItemOutput
	Err error
}

// GetItem discards the input parameters and returns Out and Err fields set on
// FakeItemGetter.
func (f *FakeDynamoDBGetter) GetItem(
	context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options),
) (*dynamodb.GetItemOutput, error) {
	return f.Out, f.Err
}

// FakeDynamoDBPutter is a test fake for ItemPutter.
type FakeDynamoDBPutter struct {
	Out *dynamodb.PutItemOutput
	Err error
}

// PutItem discards the input parameters and returns Out and Err fields set on
// FakeItemPutter.
func (f *FakeDynamoDBPutter) PutItem(
	context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	return f.Out, f.Err
}
