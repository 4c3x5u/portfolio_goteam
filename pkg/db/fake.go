package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// FakeItemPutter is a test fake for ItemPutter.
type FakeItemPutter struct {
	Out *dynamodb.PutItemOutput
	Err error
}

// PutItem discards the input parameters and returns Out and Err fields set on
// FakeItemPutter.
func (f *FakeItemPutter) PutItem(
	context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	return f.Out, f.Err
}

// FakeItemGetter is a test fake for ItemGetter.
type FakeItemGetter struct {
	Out *dynamodb.GetItemOutput
	Err error
}

// GetItem discards the input parameters and returns Out and Err fields set on
// FakeItemGetter.
func (f *FakeItemGetter) GetItem(
	context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options),
) (*dynamodb.GetItemOutput, error) {
	return f.Out, f.Err
}
