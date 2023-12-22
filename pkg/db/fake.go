//go:build utest

package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// FakeRetriever is a test fake for Retriever.
type FakeRetriever[T any] struct {
	Res T
	Err error
}

// Retrieve discards params and returns FakeRetriever.Item and
// FakeRetriever.Err.
func (f *FakeRetriever[T]) Retrieve(context.Context, string) (T, error) {
	return f.Res, f.Err
}

// FakeInserter is a test fake for Inserter.
type FakeInserter[T any] struct{ Err error }

// Insert discards params and returns FakeInserter.Err.
func (f *FakeInserter[T]) Insert(context.Context, T) error { return f.Err }

// FakeUpdater is a test fake for Updater.
type FakeUpdater[T any] struct{ Err error }

// Update discards params and returns FakeUpdater.Err.
func (f *FakeUpdater[T]) Update(context.Context, T) error { return f.Err }

// FakeDeleter is a test fake for Deleter.
type FakeDeleter struct{ Err error }

// Delete discards params and returns FakeDeleter.Err.
func (f *FakeDeleter) Delete(context.Context, string) error { return f.Err }

// FakeInserterDualKey is a test fake for InserterDualKey.
type FakeInserterDualKey[T any] struct{ Err error }

// Insert discards params and returns FakeInserterDualKey.Err.
func (f *FakeInserterDualKey[T]) Insert(context.Context, string, T) error {
	return f.Err
}

// FakeUpdaterDualKey is a test fake for UpdaterDualKey.
type FakeUpdaterDualKey[T any] struct{ Err error }

// Update discards params and returns FakeUpdaterDualKey.Err.
func (f *FakeUpdaterDualKey[T]) Update(context.Context, string, T) error {
	return f.Err
}

// FakeDeleterDualKey is a test fake for DeleterDualKey.
type FakeDeleterDualKey struct{ Err error }

// Delete discards params and returns FakeDeleterDualKey.Err.
func (f *FakeDeleterDualKey) Delete(context.Context, string, string) error {
	return f.Err
}

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

// FakeDynamoQueryer is a test fake for DynamoQueryer.
type FakeDynamoQueryer struct {
	Out *dynamodb.QueryOutput
	Err error
}

// Query discards the input parameters and returns Out and Err fields set on
// FakeDynamoQueryer.
func (f *FakeDynamoQueryer) Query(
	context.Context, *dynamodb.QueryInput, ...func(*dynamodb.Options),
) (*dynamodb.QueryOutput, error) {
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

// FakeDynamoTransactWriter is a test fake for DynamoTransactWriter.
type FakeDynamoTransactWriter struct {
	Out *dynamodb.TransactWriteItemsOutput
	Err error
}

// TransactWriteItems discards the input parameters and returns Out and Err
// fields set on FakeDynamoTransactWriter.
func (f *FakeDynamoTransactWriter) TransactWriteItems(
	context.Context,
	*dynamodb.TransactWriteItemsInput,
	...func(*dynamodb.Options),
) (*dynamodb.TransactWriteItemsOutput, error) {
	return f.Out, f.Err
}

// FakeDynamoItemGetPutter is a test fake for DynamoItemGetPutter.
type FakeDynamoItemGetPutter struct {
	OutGet *dynamodb.GetItemOutput
	ErrGet error
	OutPut *dynamodb.PutItemOutput
	ErrPut error
}

// GetItem discards the input parameters and returns OutGet and ErrGet fields
// set on FakeDynamoItemGetPutter.
func (f *FakeDynamoItemGetPutter) GetItem(
	context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options),
) (*dynamodb.GetItemOutput, error) {
	return f.OutGet, f.ErrGet
}

// PutItem discards the input parameters and returns OutPut and ErrPut fields
// set on FakeDynamoItemGetPutter.
func (f *FakeDynamoItemGetPutter) PutItem(
	context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	return f.OutPut, f.ErrPut
}
