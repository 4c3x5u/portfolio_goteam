package user

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestItemPutter(t *testing.T) {
	i := User{}
	ip := &db.FakeDynamoDBItemPutter{}
	sut := NewItemPutter(ip)

	t.Run("ErrIDExistsWhenConditonalCheckFailed", func(t *testing.T) {
		wantErr := ErrIDExists
		ip.Err = &smithy.OperationError{
			Err: &types.ConditionalCheckFailedException{},
		}

		err := sut.Put(i)

		assert.ErrIs(t.Fatal, err, wantErr)
	})

	t.Run("DirectErrReturnWhenOtherErr", func(t *testing.T) {
		wantErr := errors.New("failed to put item")
		ip.Err = wantErr

		err := sut.Put(i)

		assert.ErrIs(t.Fatal, err, wantErr)
	})

	t.Run("ErrIsNilWhenOK", func(t *testing.T) {
		ip.Err = nil

		err := sut.Put(i)

		assert.Nil(t.Fatal, err)
	})
}
