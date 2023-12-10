package team

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/db"
)

// Retriever can be used to retrieve by ID a team from the team table.
type Retriever struct{ ItemGetter db.DynamoItemGetter }

// NewRetriever creates and returns a new Retriever.
func NewRetriever(ig db.DynamoItemGetter) Retriever { return Retriever{ItemGetter: ig} }

// Retrieve retrieves by ID a team from the team table.
func (r Retriever) Retrieve(ctx context.Context, id string) (Team, error) {
	out, err := r.ItemGetter.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv(tableName)),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return Team{}, err
	}
	if out.Item == nil {
		return Team{}, db.ErrNoItem
	}

	var t Team
	if err := attributevalue.UnmarshalMap(out.Item, &t); err != nil {
		return Team{}, err
	}

	return t, nil
}
