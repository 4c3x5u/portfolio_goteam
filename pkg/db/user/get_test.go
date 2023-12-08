//go:build utest

package user

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestGetter(t *testing.T) {
	ig := &db.FakeItemGetter{}
	sut := NewGetter(ig)

	userA := User{
		Username: "bob123",
		Password: []byte("p4ssw0rd"),
		IsAdmin:  true,
		TeamID:   "afeadc4a-68b0-4c33-9e83-4648d20ff26a",
	}
	errA := errors.New("failed to get item")

	for _, c := range []struct {
		name     string
		igOut    *dynamodb.GetItemOutput
		igErr    error
		wantUser *User
		wantErr  error
	}{
		{
			name:     "Err",
			igOut:    nil,
			igErr:    errA,
			wantUser: nil,
			wantErr:  errA,
		},
		{
			name:     "NoItem",
			igOut:    &dynamodb.GetItemOutput{Item: nil},
			igErr:    nil,
			wantUser: nil,
			wantErr:  db.ErrNoItem,
		},
		{
			name: "OK",
			igOut: &dynamodb.GetItemOutput{
				Item: map[string]types.AttributeValue{
					"Username": &types.AttributeValueMemberS{Value: userA.Username},
					"Password": &types.AttributeValueMemberB{
						Value: userA.Password,
					},
					"IsAdmin": &types.AttributeValueMemberBOOL{
						Value: userA.IsAdmin,
					},
					"TeamID": &types.AttributeValueMemberS{
						Value: userA.TeamID,
					},
				},
			},
			igErr:    nil,
			wantUser: &userA,
			wantErr:  nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ig.Out = c.igOut
			ig.Err = c.igErr

			user, err := sut.Get(context.Background(), "")

			assert.Equal(t.Fatal, err, c.wantErr)
			if c.wantUser != nil {
				assert.Equal(t.Error, user.Username, c.wantUser.Username)
				assert.AllEqual(t.Error, user.Password, c.wantUser.Password)
				assert.True(t.Error, c.wantUser.IsAdmin)
				assert.Equal(t.Error, user.TeamID, c.wantUser.TeamID)
			}
		})
	}
}
