//go:build utest

package tasktable

import (
	"context"
	"errors"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/aws/smithy-go"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestDelete(t *testing.T) {
	idel := &db.FakeDynamoItemDeleter{}
	sut := NewDeleter(idel)

	errA := errors.New("failed")

	for _, c := range []struct {
		name    string
		idelErr error
		wantErr error
	}{
		{name: "Err", idelErr: errA, wantErr: errA},
		{
			name: "Err",
			idelErr: &smithy.OperationError{
				Err: &types.ConditionalCheckFailedException{},
			},
			wantErr: db.ErrNoItem,
		},
		{name: "OK", idelErr: nil, wantErr: nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			idel.Err = c.idelErr

			err := sut.Delete(context.Background(), "", "")

			assert.Equal(t.Fatal, err, c.wantErr)
		})
	}
}
