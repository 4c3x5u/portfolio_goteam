package user

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/db"
)

// Getter can be used to get a user from the user table by ID.
type Getter struct{ ItemGetter db.ItemGetteer }

// NewGetter creates and returns a new Getter.
func NewGetter(ig db.ItemGetteer) Getter { return Getter{ItemGetter: ig} }

// Get gets user from the user table by ID.
func (g Getter) Get(id string) (User, error) {
	out, err := g.ItemGetter.GetItem(context.TODO(), &dynamodb.GetItemInput{
		Key: map[string]types.AttributeValue{
			"ID": &types.AttributeValueMemberS{Value: id},
		},
	})
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
