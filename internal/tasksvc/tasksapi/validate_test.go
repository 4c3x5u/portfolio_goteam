//go:build utest

package tasksapi

import (
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/validator"
)

// TestColNoValidator tests the ColNoValidator.Validate method.
func TestColNoValidator(t *testing.T) {
	sut := NewColNoValidator()

	for _, c := range []struct {
		name    string
		colNo   int
		wantErr error
	}{
		{
			name:    "ColNoTooSmall",
			colNo:   -1,
			wantErr: validator.ErrOutOfBounds,
		},
		{
			name:    "ColNoTooBig",
			colNo:   4,
			wantErr: validator.ErrOutOfBounds,
		},
		{
			name:    "Success",
			colNo:   2,
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Validate(c.colNo)
			assert.ErrIs(t.Error, err, c.wantErr)
		})
	}
}
