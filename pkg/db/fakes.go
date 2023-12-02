package db

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// FakeDynamoDBItemPutter is a test fake for DynamoDBItemPutter.
type FakeDynamoDBItemPutter struct {
	Out *dynamodb.PutItemOutput
	Err error
}

// PutItem discards the input parameters and returns Out and Err fields set on
// FakeItemPutter.
func (f *FakeDynamoDBItemPutter) PutItem(
	context.Context, *dynamodb.PutItemInput, ...func(*dynamodb.Options),
) (*dynamodb.PutItemOutput, error) {
	return f.Out, f.Err
}

// FakeDynamoDBItemGetter is a test fake for DynamoDBItemGetter.
type FakeDynamoDBItemGetter struct {
	Out *dynamodb.GetItemOutput
	Err error
}

// GetItem discards the input parameters and returns Out and Err fields set on
// FakeItemPutter.
func (f *FakeDynamoDBItemGetter) GetItem(
	context.Context, *dynamodb.GetItemInput, ...func(*dynamodb.Options),
) (*dynamodb.GetItemOutput, error) {
	return f.Out, f.Err
}
