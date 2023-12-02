package user

import (
	"context"

	"github.com/kxplxn/goteam/pkg/db"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// GetterByTeamID can be used to get an item from the user table by its team ID.
type GetterByTeamID struct {
	ItemGetter db.DynamoDBItemGetter
}

// NewGetterByTeamID creates and returns a new GetterByTeamID.
func NewGetterByTeamID(ig db.DynamoDBItemGetter) GetterByTeamID {
	return GetterByTeamID{ItemGetter: ig}
}

// Get gets an item from the user table by its team ID.
func (g GetterByTeamID) Get(teamID string) (User, error) {
	out, err := g.ItemGetter.GetItem(context.TODO(), &dynamodb.GetItemInput{})
	if err != nil {
		return User{}, err
	}
	if out == nil {
		return User{}, db.ErrNoItem
	}

	var user User
	if err := attributevalue.UnmarshalMap(out.Item, &user); err != nil {
		return User{}, err
	}

	return user, nil
}
