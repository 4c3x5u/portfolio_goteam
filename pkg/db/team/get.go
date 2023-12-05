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

// Getter can be used to get a team from the team table by ID.
type Getter struct{ ItemGetter db.ItemGetter }

// NewGetter creates and returns a new Getter.
func NewGetter(ig db.ItemGetter) Getter { return Getter{ItemGetter: ig} }

// Get gets a team from the team table by ID.
func (g Getter) Get(ctx context.Context, id string) (Team, error) {
	out, err := g.ItemGetter.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv(tableName)),
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
	if err != nil {
		return Team{}, err
	}
	if out == nil {
		return Team{}, db.ErrNoItem
	}

	var t Team
	if err := attributevalue.UnmarshalMap(out.Item, &t); err != nil {
		return Team{}, err
	}

	return t, nil
}
