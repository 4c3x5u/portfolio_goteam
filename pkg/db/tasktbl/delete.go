package tasktbl

import (
	"context"
	"errors"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/db"
)

// Deleter can be used to retrieve by ID a task from the task table.
type Deleter struct{ idel db.DynamoItemDeleter }

// NewDeleter creates and returns a new Getter.
func NewDeleter(idel db.DynamoItemDeleter) Deleter {
	return Deleter{idel: idel}
}

// Retrieve retrieves by ID a task from the task table.
func (r Deleter) Delete(ctx context.Context, teamID, taskID string) error {
	_, err := r.idel.DeleteItem(ctx, &dynamodb.DeleteItemInput{
		TableName: aws.String(os.Getenv(tableName)),
		Key: map[string]types.AttributeValue{
			"TeamID": &types.AttributeValueMemberS{Value: teamID},
			"ID":     &types.AttributeValueMemberS{Value: taskID},
		},
		ConditionExpression: aws.String("attribute_exists(ID)"),
	})

	var ex *types.ConditionalCheckFailedException
	if errors.As(err, &ex) {
		return db.ErrNoItem
	}

	return err
}
