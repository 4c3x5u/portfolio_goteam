//go:build utest

package teamtbl

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestInserter(t *testing.T) {
	ip := &db.FakeDynamoItemPutter{}
	sut := NewInserter(ip)

	errA := errors.New("failed to create item")

	for _, c := range []struct {
		name    string
		ipErr   error
		wantErr error
	}{
		{name: "Err", ipErr: errA, wantErr: errA},
		{
			name: "DupKey",
			ipErr: &smithy.OperationError{
				Err: &types.ConditionalCheckFailedException{},
			},
			wantErr: db.ErrDupKey,
		},
		{name: "OK", ipErr: nil, wantErr: nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			ip.Err = c.ipErr

			err := sut.Insert(context.Background(), Team{})

			assert.ErrIs(t.Fatal, err, c.wantErr)
		})
	}
}
