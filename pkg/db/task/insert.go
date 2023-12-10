package task

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

// Inserter can be used to insert a new task into the task table.
type Inserter struct{ ItemPutter db.DynamoItemPutter }

// NewInserter creates and returns a new Inserter.
func NewInserter(ip db.DynamoItemPutter) Inserter {
	return Inserter{ItemPutter: ip}
}

// Insert inserts a new task into the task table.
func (u Inserter) Insert(ctx context.Context, task Task) error {
	item, err := attributevalue.MarshalMap(task)
	if err != nil {
		return err
	}

	_, err = u.ItemPutter.PutItem(ctx, &dynamodb.PutItemInput{
		TableName:           aws.String(os.Getenv(tableName)),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(ID)"),
	})

	var ex *types.ConditionalCheckFailedException
	if errors.As(err, &ex) {
		return db.ErrDupKey
	}

	return err
}
