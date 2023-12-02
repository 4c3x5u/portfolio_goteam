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

func TestGetterByTeamID(t *testing.T) {
	ig := &db.FakeDynamoDBItemGetter{}
	sut := NewGetterByTeamID(ig)

	t.Run("Err", func(t *testing.T) {
		wantErr := errors.New("failed to get item")
		ig.Err = wantErr

		_, err := sut.Get("")

		assert.ErrIs(t.Fatal, err, wantErr)
	})

	t.Run("ErrNoItem", func(t *testing.T) {
		ig.Out = nil
		ig.Err = nil

		_, err := sut.Get("")

		assert.ErrIs(t.Fatal, err, db.ErrNoItem)
	})

	t.Run("OK", func(t *testing.T) {
		ig.Out = &dynamodb.GetItemOutput{
			Item: map[string]types.AttributeValue{
				"ID":       &types.AttributeValueMemberS{Value: "bob124"},
				"Password": &types.AttributeValueMemberB{Value: []byte("pwd")},
				"IsAdmin":  &types.AttributeValueMemberBOOL{Value: false},
				"TeamID":   &types.AttributeValueMemberN{Value: "22"},
			},
		}
		ig.Err = nil

		user, err := sut.Get("")

		assert.Nil(t.Fatal, err)
		assert.Equal(t.Error, user.ID, "bob124")
		assert.Equal(t.Error, string(user.Password), "pwd")
		assert.True(t.Error, !user.IsAdmin)
		assert.Equal(t.Error, user.TeamID, 22)
	})
}
