//go:build utest

package column

import (
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
)

// TestIDValidator tests the Validate method of IDValidator to assert
// that it returns the correct error message based on the column ID it's given.
func TestIDValidator(t *testing.T) {
	sut := NewIDValidator()

	for _, c := range []struct {
		name       string
		columnID   string
		wantErrMsg string
	}{
		{
			name:       "Nil",
			columnID:   "",
			wantErrMsg: "Column ID cannot be empty.",
		},
		{
			name:       "NotInt",
			columnID:   "A",
			wantErrMsg: "Column ID must be an integer.",
		},
		{
			name:       "Success",
			columnID:   "12",
			wantErrMsg: "",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Validate(c.columnID)

			if c.wantErrMsg == "" {
				assert.Nil(t.Error, err)
			} else {
				assert.Equal(t.Error, err.Error(), c.wantErrMsg)
			}
		})
	}
}
