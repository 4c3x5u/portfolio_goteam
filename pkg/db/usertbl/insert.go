package usertbl

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
type Inserter struct{ iput db.DynamoItemPutter }

// NewInserter creates and returns a new Inserter.
func NewInserter(iput db.DynamoItemPutter) Inserter {
	return Inserter{iput: iput}
}

// Insert inserts a new user into the user table.
func (i Inserter) Insert(ctx context.Context, user User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = i.iput.PutItem(ctx, &dynamodb.PutItemInput{
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
