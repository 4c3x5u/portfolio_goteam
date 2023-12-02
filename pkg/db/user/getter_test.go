//go:build utest

package user

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestGetter(t *testing.T) {
	ig := &db.FakeDynamoDBItemGetter{}
	sut := NewGetter(ig)

	t.Run("RDirectErrReturnWhenOtherErreturnErrWhenOccurs", func(t *testing.T) {
		wantErr := errors.New("failed to get item")
		ig.Err = wantErr

		_, err := sut.Get("")

		assert.ErrIs(t.Fatal, err, wantErr)
	})

	t.Run("ErrNotFoundWhenNilOut", func(t *testing.T) {
		ig.Out = nil
		ig.Err = nil

		_, err := sut.Get("")

		assert.ErrIs(t.Fatal, err, db.ErrNotFound)
	})

	t.Run("OK", func(t *testing.T) {
		ig.Out = &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"ID": &types.AttributeValueMemberS{Value: "bob123"},
				"Password": &types.AttributeValueMemberB{
					Value: []byte("password"),
				},
				"IsAdmin": &types.AttributeValueMemberBOOL{Value: true},
				"TeamID":  &types.AttributeValueMemberN{Value: "21"},
			},
		}
		ig.Err = nil

		user, err := sut.Get("")

		assert.Nil(t.Fatal, err)
		assert.Equal(t.Error, user.ID, "bob123")
		assert.Equal(t.Error, string(user.Password), "password")
		assert.True(t.Error, user.IsAdmin)
		assert.Equal(t.Error, user.TeamID, 21)
	})
}
