package user

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/kxplxn/goteam/pkg/db"
)

// Getter can be used to get an item from the user table.
type Getter struct {
	ItemGetter db.DynamoDBItemGetter
}

// NewGetter creates and returns a new Getter.
func NewGetter(ig db.DynamoDBItemGetter) Getter {
	return Getter{ItemGetter: ig}
}

// Get gets an item from the user table.
func (g Getter) Get(id string) (User, error) {
	out, err := g.ItemGetter.GetItem(context.TODO(), &dynamodb.GetItemInput{})
	if err != nil {
		return User{}, err
	}
	if out == nil {
		return User{}, db.ErrNoItem
	}

	var user User
	if attributevalue.UnmarshalMap(out.Item, &user); err != nil {
		return User{}, err
	}

	return user, nil
}
