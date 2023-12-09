//go:build utest

package subtask

import (
	"testing"

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/assert"
)

// TestIDValidator tests the Validate method of IDValidator to assert
// that it returns the correct error message based on the board ID it's given.
func TestIDValidator(t *testing.T) {
	sut := NewIDValidator()

	for _, c := range []struct {
		name    string
		id      string
		wantErr error
	}{
		{
			name:    "Empty",
			id:      "",
			wantErr: api.ErrEmpty,
		},
		{
			name:    "NotInt",
			id:      "A",
			wantErr: api.ErrNotInt,
		},
		{
			name:    "Success",
			id:      "1",
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Validate(c.id)
			assert.ErrIs(t.Error, err, c.wantErr)
		})
	}
}
