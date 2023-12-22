//go:build utest

package team

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestBoardDeleter(t *testing.T) {
	igetput := &db.FakeDynamoItemGetPutter{}
	sut := NewBoardDeleter(igetput)

	errA := errors.New("failed")
	itemA := map[string]types.AttributeValue{
		"Boards": &types.AttributeValueMemberL{
			Value: []types.AttributeValue{
				&types.AttributeValueMemberM{
					Value: map[string]types.AttributeValue{
						"ID": &types.AttributeValueMemberS{
							Value: "boardID",
						},
						"Name": &types.AttributeValueMemberS{
							Value: "boardName",
						},
					},
				},
			},
		},
	}

	for _, c := range []struct {
		name       string
		errGetItem error
		outGetItem *dynamodb.GetItemOutput
		errPutItem error
		wantErr    error
	}{
		{
			name:       "ErrGetItem",
			errGetItem: errA,
			outGetItem: nil,
			errPutItem: nil,
			wantErr:    errA,
		},
		{
			name:       "ErrNoItemTeam",
			errGetItem: nil,
			outGetItem: &dynamodb.GetItemOutput{Item: nil},
			errPutItem: nil,
			wantErr:    db.ErrNoItem,
		},
		{
			name:       "ErrNoItemBoard",
			errGetItem: nil,
			outGetItem: &dynamodb.GetItemOutput{
				Item: map[string]types.AttributeValue{
					"Boards": &types.AttributeValueMemberL{
						Value: []types.AttributeValue{},
					},
				},
			},
			errPutItem: nil,
			wantErr:    db.ErrNoItem,
		},
		{
			name:       "ErrPutItem",
			errGetItem: nil,
			outGetItem: &dynamodb.GetItemOutput{Item: itemA},
			errPutItem: errA,
			wantErr:    errA,
		},
		{
			name:       "OK",
			errGetItem: nil,
			outGetItem: &dynamodb.GetItemOutput{Item: itemA},
			errPutItem: nil,
			wantErr:    nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			igetput.ErrGet = c.errGetItem
			igetput.OutGet = c.outGetItem
			igetput.ErrPut = c.errPutItem

			err := sut.Delete(context.Background(), "", "boardID")

			assert.Equal(t.Fatal, err, c.wantErr)
		})
	}
}
