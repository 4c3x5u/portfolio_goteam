package tasktbl

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/db"
)

// Retriever can be used to retrieve by ID a task from the task table.
type Retriever struct{ ItemGetter db.DynamoItemGetter }

// NewRetriever creates and returns a new Getter.
func NewRetriever(ig db.DynamoItemGetter) Retriever {
	return Retriever{ItemGetter: ig}
}

// Retrieve retrieves by ID a task from the task table.
func (r Retriever) Retrieve(ctx context.Context, id string) (Task, error) {
	out, err := r.ItemGetter.GetItem(ctx, &dynamodb.GetItemInput{
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
