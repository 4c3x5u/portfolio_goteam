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
	ip := &db.FakeItemPutter{}
	sut := NewPutter(ip)

	errA := errors.New("failed to put item")

	for _, c := range []struct {
		name    string
		ipErr   error
		wantErr error
	}{
		{
			name: "DupKey",
			ipErr: &smithy.OperationError{
				Err: &types.ConditionalCheckFailedException{},
			},
			wantErr: db.ErrDupKey,
		},
		{
			name:    "Err",
			ipErr:   errA,
			wantErr: errA,
		},
		{
			name:    "OK",
			ipErr:   nil,
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ip.Err = c.ipErr

			err := sut.Put(User{})

			assert.ErrIs(t.Fatal, err, c.wantErr)
		})
	}
}
