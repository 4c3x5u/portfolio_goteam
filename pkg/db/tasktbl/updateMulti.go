package tasktbl

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

// MultiUpdater can be used to update multiple tasks in the task table at once.
type MultiUpdater struct{ TransactWriter db.DynamoTransactWriter }

// NewMultiUpdater creates and returns a new MultiUpdater.
func NewMultiUpdater(tw db.DynamoTransactWriter) MultiUpdater {
	return MultiUpdater{TransactWriter: tw}
}

// Update updates multiple tasks in the task table at once.
func (u MultiUpdater) Update(ctx context.Context, tasks []Task) error {
	tableName := os.Getenv("TASK_TABLE_NAME")

	items := make([]types.TransactWriteItem, len(tasks))
	for i, task := range tasks {
		item, err := attributevalue.MarshalMap(task)
		if err != nil {
			return err
		}
		items[i] = types.TransactWriteItem{
			Put: &types.Put{
				TableName:           &tableName,
				Item:                item,
				ConditionExpression: aws.String("attribute_exists(ID)"),
			},
		}
	}

	_, err := u.TransactWriter.TransactWriteItems(
		ctx,
		&dynamodb.TransactWriteItemsInput{TransactItems: items},
	)

	var ex *types.ConditionalCheckFailedException
	if errors.As(err, &ex) {
		return db.ErrNoItem
	}

	return err
}
