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
	ItemGetter db.ItemGetter
}

// NewGetter creates and returns a new Getter.
func NewGetter(ig db.ItemGetter) Getter { return Getter{ItemGetter: ig} }

// Get gets a task from the task table.
func (g Getter) Get(id string) (Task, error) {
	out, err := g.ItemGetter.GetItem(
		context.TODO(),
		&dynamodb.GetItemInput{
			TableName: aws.String(os.Getenv("DYNAMODB_TABLE_TASK")),
			Key: map[string]types.AttributeValue{
				"ID": &types.AttributeValueMemberS{Value: id},
			},
		},
	)
	if err != nil {
		return Task{}, err
	}
	if out == nil {
		return Task{}, db.ErrNoItem
	}

	var task Task
	err = attributevalue.UnmarshalMap(out.Item, &task)
	return task, err
}
