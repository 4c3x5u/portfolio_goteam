//go:build utest

package team

import (
	"errors"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
)

func TestPutter(t *testing.T) {
	ip := &db.FakeItemPutter{}
	sut := NewPutter(ip)

	t.Run("Err", func(t *testing.T) {
		wantErr := errors.New("failed to put item")
		ip.Out = nil
		ip.Err = wantErr

		err := sut.Put(Team{})

		assert.ErrIs(t.Fatal, err, wantErr)
	})

	t.Run("OK", func(t *testing.T) {
		ip.Err = nil

		err := sut.Put(Team{})

		assert.Nil(t.Fatal, err)
	})
}
