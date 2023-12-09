package user

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/db"
)

// Inserter can be used to insert a new user into the user table.
type Inserter struct{ ItemPutter db.DynamoDBPutter }

// NewInserter creates and returns a new Inserter.
func NewInserter(ip db.DynamoDBPutter) Inserter {
	return Inserter{ItemPutter: ip}
}

// Insert inserts a new user into the user table.
func (i Inserter) Insert(ctx context.Context, user User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = i.ItemPutter.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(os.Getenv(tableName)),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(Username)"),
	})

	var ex *types.ConditionalCheckFailedException
	if errors.As(err, &ex) {
		return db.ErrDupKey
	}

	return err
}
