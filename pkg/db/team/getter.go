package team

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kxplxn/goteam/pkg/db"
)

// Getter can be used to get an item from the team table.
type Getter struct {
	ItemGetter db.DynamoDBItemGetter
}

// NewGetter creates and returns a new Getter.
func NewGetter(ig db.DynamoDBItemGetter) Getter {
	return Getter{ItemGetter: ig}
}

func (g Getter) Get(id string) (Team, error) {
	out, err := g.ItemGetter.GetItem(context.TODO(), &dynamodb.GetItemInput{})
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
