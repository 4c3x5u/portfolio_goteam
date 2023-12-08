package user

import (
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/db"
)

// Getter can be used to get a user from the user table by ID.
type Getter struct{ ItemGetter db.DynamoDBGetter }

// NewGetter creates and returns a new Getter.
func NewGetter(ig db.DynamoDBGetter) Getter { return Getter{ItemGetter: ig} }

// Get gets user from the user table by ID.
func (g Getter) Get(ctx context.Context, username string) (User, error) {
	out, err := g.ItemGetter.GetItem(ctx, &dynamodb.GetItemInput{
		TableName: aws.String(os.Getenv(tableName)),
		Key: map[string]types.AttributeValue{
			"Username": &types.AttributeValueMemberS{Value: username},
		},
	})
	if err != nil {
		return User{}, err
	}
	if out.Item == nil {
		return User{}, db.ErrNoItem
	}

	var user User
	if attributevalue.UnmarshalMap(out.Item, &user); err != nil {
		return User{}, err
	}
	return user, nil
}
