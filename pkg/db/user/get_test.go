//go:build utest

package user

import (
	"errors"
	"strconv"
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
		ID:       "bob123",
		Password: []byte("p4ssw0rd"),
		IsAdmin:  true,
		TeamID:   21,
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
			igOut:    nil,
			igErr:    nil,
			wantUser: nil,
			wantErr:  db.ErrNoItem,
		},
		{
			name: "OK",
			igOut: &dynamodb.GetItemOutput{
				Item: map[string]types.AttributeValue{
					"ID": &types.AttributeValueMemberS{Value: userA.ID},
					"Password": &types.AttributeValueMemberB{
						Value: userA.Password,
					},
					"IsAdmin": &types.AttributeValueMemberBOOL{
						Value: userA.IsAdmin,
					},
					"TeamID": &types.AttributeValueMemberN{
						Value: strconv.Itoa(userA.TeamID),
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

			user, err := sut.Get("")

			assert.Equal(t.Fatal, err, c.wantErr)
			if c.wantUser != nil {
				assert.Equal(t.Error, user.ID, c.wantUser.ID)
				assert.AllEqual(t.Error, user.Password, c.wantUser.Password)
				assert.True(t.Error, c.wantUser.IsAdmin)
				assert.Equal(t.Error, user.TeamID, c.wantUser.TeamID)
			}
		})
	}
}
