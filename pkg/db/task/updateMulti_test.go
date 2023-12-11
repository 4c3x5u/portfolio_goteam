//go:build utest

package task

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestMultiUpdater(t *testing.T) {
	tw := &db.FakeDynamoTransactWriter{}
	sut := NewMultiUpdater(tw)

	errA := errors.New("failed to put item")

	for _, c := range []struct {
		name    string
		ipErr   error
		wantErr error
	}{
		{name: "Err", ipErr: errA, wantErr: errA},
		{
			name: "NoItem",
			ipErr: &smithy.OperationError{
				Err: &types.ConditionalCheckFailedException{},
			},
			wantErr: db.ErrNoItem,
		},
		{name: "OK", ipErr: nil, wantErr: nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			tw.Err = c.ipErr

			err := sut.Update(context.Background(), []Task{})

			assert.ErrIs(t.Fatal, err, c.wantErr)
		})
	}
}
