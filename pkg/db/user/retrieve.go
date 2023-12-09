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

// Retriever can be used to retrieve by username a user from the user table.
type Retriever struct{ ItemGetter db.DynamoDBGetter }

// NewRetriever creates and returns a new Retriever.
func NewRetriever(ig db.DynamoDBGetter) Retriever {
	return Retriever{ItemGetter: ig}
}

// Retrieve retrieves by username a user from the user table.
func (g Retriever) Retrieve(ctx context.Context, username string) (User, error) {
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
