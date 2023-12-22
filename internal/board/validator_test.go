//go:build utest

package board

import (
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
)

func TestNameValidator(t *testing.T) {
	sut := NewNameValidator()

	for _, c := range []struct {
		name      string
		boardName string
		wantErr   error
	}{
		{
			name:      "Empty",
			boardName: "",
			wantErr:   ErrEmpty,
		},
		{
			name:      "TooLong",
			boardName: "boardyboardsyboardkyboardishboardxyz",
			wantErr:   ErrTooLong,
		},
		{
			name:      "OK",
			boardName: "My Board",
			wantErr:   nil,
		},
	} {
		err := sut.Validate(c.boardName)

		assert.ErrIs(t.Error, err, c.wantErr)
	}
}
