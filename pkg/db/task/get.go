package task

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/db"
)

// Getter can be used to get a task from the task table.
type Getter struct {
	ItemGetter db.DynamoDBGetter
}

// NewGetter creates and returns a new Getter.
func NewGetter(ig db.DynamoDBGetter) Getter { return Getter{ItemGetter: ig} }

// Get gets a task from the task table.
func (g Getter) Get(ctx context.Context, id string) (Task, error) {
	out, err := g.ItemGetter.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv(tableName)),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return Task{}, err
	}
	if out.Item == nil {
		return Task{}, db.ErrNoItem
	}

	var task Task
	err = attributevalue.UnmarshalMap(out.Item, &task)
	return task, err
}
