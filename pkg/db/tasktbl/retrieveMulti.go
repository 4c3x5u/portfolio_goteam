package tasktbl

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
type MultiRetriever struct{ queryer db.DynamoQueryer }

// NewMultiRetriever creates and returns a new MultiRetriever.
func NewMultiRetriever(queryer db.DynamoQueryer) MultiRetriever {
	return MultiRetriever{queryer: queryer}
}

// Retrieve retrieves all tasks for a team from the task table.
func (r MultiRetriever) Retrieve(
	ctx context.Context, boardID string,
) ([]Task, error) {
	keyCond := expression.Key("BoardID").Equal(expression.Value(boardID))
	expr, err := expression.NewBuilder().WithKeyCondition(keyCond).Build()
	if err != nil {
		return nil, err
	}

	out, err := r.queryer.Query(ctx, &dynamodb.QueryInput{
		TableName:                 aws.String(os.Getenv(tableName)),
		IndexName:                 aws.String("BoardID_index"),
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
