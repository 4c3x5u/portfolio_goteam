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

// FakeDeleter is a test fake for Deleter.
type FakeDeleter struct{ Err error }

// Delete discards params and returns FakeDeleter.Err.
func (f *FakeDeleter) Delete(_ context.Context) error { return f.Err }

// FakeDynamoItemGetter is a test fake for DynamoItemGetter.
type FakeDynamoItemGetter struct {
	Out *dynamodb.GetItemOutput
	Err error
}

// GetItem discards the input parameters and returns Out and Err fields set on
// FakeDynamoItemGetter.
func (f *FakeDynamoItemGetter) GetItem(
	context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options),
) (*dynamodb.GetItemOutput, error) {
	return f.Out, f.Err
}

// FakeDynamoItemPutter is a test fake for DynamoItemPutter.
type FakeDynamoItemPutter struct {
	Out *dynamodb.PutItemOutput
	Err error
}

// PutItem discards the input parameters and returns Out and Err fields set on
// FakeDynamoItemPutter.
func (f *FakeDynamoItemPutter) PutItem(
	context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	return f.Out, f.Err
}

// FakeDynamoItemDeleter is a test fake for DynamoItemDeleter.
type FakeDynamoItemDeleter struct {
	Out *dynamodb.DeleteItemOutput
	Err error
}

// PutItem discards the input parameters and returns Out and Err fields set on
// FakeDynamoItemDeleter.
func (f *FakeDynamoItemDeleter) DeleteItem(
	context.Context, *dynamodb.DeleteItemInput, ...func(*dynamodb.Options),
) (*dynamodb.DeleteItemOutput, error) {
	return f.Out, f.Err
}
