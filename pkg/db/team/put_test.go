//go:build utest

package team

import (
	"context"
	"errors"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestPutter(t *testing.T) {
	ip := &db.FakeDynamoDBPutter{}
	sut := NewPutter(ip)

	errA := errors.New("failed to put item")

	for _, c := range []struct {
		name    string
		ipErr   error
		wantErr error
	}{
		{name: "Err", ipErr: errA, wantErr: errA},
		{name: "OK", ipErr: nil, wantErr: nil},
	} {
		t.Run(c.name, func(t *testing.T) {
			ip.Err = c.ipErr

			err := sut.Put(context.Background(), Team{})

			assert.ErrIs(t.Fatal, err, c.wantErr)
		})
	}
}
