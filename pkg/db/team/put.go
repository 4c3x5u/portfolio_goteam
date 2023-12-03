package team

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/kxplxn/goteam/pkg/db"
)

// Putter can be used to put a team into the team table.
type Putter struct{ ItemPutter db.ItemPutter }

// NewPutter creates and returns a new Putter.
func NewPutter(ip db.ItemPutter) Putter { return Putter{ItemPutter: ip} }

// Put puts a team into the team table.
func (p Putter) Put(team Team) error {
	item, err := attributevalue.MarshalMap(team)
	if err != nil {
		return err
	}

	_, err = p.ItemPutter.PutItem(
		context.TODO(),
		&dynamodb.PutItemInput{
			TableName: aws.String(os.Getenv("DYNAMODB_TABLE_TEAM")),
			Item:      item,
		},
	)

	return err
}
