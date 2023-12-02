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

// ErrIDExists means that the ID passed in the user data already exists in user
// table.
var ErrIDExists = errors.New("user ID already exists")

// ItemPutter can be used to "put" an item into user table.
type Putter struct {
	ItemPutter db.DynamoDBItemPutter
}

// NewPutter creates and returns a new ItemPutter.
func NewPutter(ip db.DynamoDBItemPutter) Putter {
	return Putter{ItemPutter: ip}
}

// Put puts an item into the user table with the given user data. If an item
// with the given ID already exists, it returns ErrIDExists instead.
func (p Putter) Put(user User) error {
	item, err := attributevalue.MarshalMap(user)
	if err != nil {
		return err
	}

	_, err = p.ItemPutter.PutItem(context.TODO(), &dynamodb.PutItemInput{
		TableName:           aws.String(os.Getenv("DYNAMODB_TABLE_USER")),
		Item:                item,
		ConditionExpression: aws.String("attribute_not_exists(ID)"),
	})

	var ex *types.ConditionalCheckFailedException
	if errors.As(err, &ex) {
		return ErrIDExists
	}

	return err
}
