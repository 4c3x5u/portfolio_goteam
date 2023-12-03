package task

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/kxplxn/goteam/pkg/db"
)

// Putter can be used to put a task into the task table.
type Putter struct{ ItemPutter db.ItemPutter }

// NewPutter creates and returns a new Putter.
func NewPutter(ip db.ItemPutter) Putter { return Putter{ItemPutter: ip} }

// Put puts a task into the task table.
func (p Putter) Put(task Task) error {
	item, err := attributevalue.MarshalMap(task)
	if err != nil {
		return err
	}

	_, err = p.ItemPutter.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv(tableName)),
			Item:      item,
		},
	)

	return err
}
