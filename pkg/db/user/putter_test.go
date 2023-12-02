package user

import (
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestPutter(t *testing.T) {
	i := User{}
	ip := &db.FakeDynamoDBItemPutter{}
	sut := NewPutter(ip)

	t.Run("ErrDupKeyWhenConditonalCheckFailed", func(t *testing.T) {
		wantErr := db.ErrDupKey
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

	t.Run("OK", func(t *testing.T) {
		ip.Err = nil

		err := sut.Put(i)

		assert.Nil(t.Fatal, err)
	})
}
