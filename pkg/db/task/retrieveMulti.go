package task

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"

	"github.com/kxplxn/goteam/pkg/db"
)

// MultiRetriever can be used to retrieve all tasks for a team from the task
// table.
type MultiRetriever struct{ Queryer db.DynamoQueryer }

// NewMultiRetriever creates and returns a new MultiRetriever.
func NewMultiRetriever(dq db.DynamoQueryer) MultiRetriever {
	return MultiRetriever{Queryer: dq}
}

// Retrieve retrieves all tasks for a team from the task table.
func (r MultiRetriever) Retrieve(
	ctx context.Context, teamID string,
) ([]Task, error) {
	keyEx := expression.Key("TeamID").Equal(expression.Value(teamID))
	expr, err := expression.NewBuilder().WithKeyCondition(keyEx).Build()
	if err != nil {
		return nil, err
	}

	out, err := r.Queryer.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(os.Getenv(tableName)),
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		KeyConditionExpression:    expr.KeyCondition(),
	})
	if err != nil {
		return nil, err
	}

	var tasks []Task
	err = attributevalue.UnmarshalListOfMaps(out.Items, &tasks)
	return tasks, err
}
