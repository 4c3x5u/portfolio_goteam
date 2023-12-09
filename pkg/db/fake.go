//go:build utest

package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// FakePutter is a test fake for Putter.
type FakePutter[T any] struct{ Err error }

// Put discards params and returns FakePutter.Err.
func (f *FakePutter[T]) Put(_ context.Context, _ T) error { return f.Err }

// FakeGetter is a test fake for Getter.
type FakeGetter[T any] struct {
	Item T
	Err  error
}

// Get discards params and returns FakeGetter.User and FakeGetter.Err.
func (f *FakeGetter[T]) Get(ctx context.Context, username string) (T, error) {
	return f.Item, f.Err
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
