//go:build utest

package board

import (
	"testing"

	"server/assert"
)

// TestDELETEValidator tests the Validate method of DELETEValidator to assert
// that it returns the correct error based on the URL query parameters.
func TestDELETEValidator(t *testing.T) {
	sut := NewDELETEValidator()

	for _, c := range []struct {
		name    string
		boardID string
		wantOK  bool
	}{
		{
			name:    "Nil",
			boardID: "",
			wantOK:  false,
		},
		{
			name:    "NotInt",
			boardID: "My Board",
			wantOK:  false,
		},
		{
			name:    "Success",
			boardID: "12",
			wantOK:  true,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			ok := sut.Validate(c.boardID)

			if err := assert.Equal(c.wantOK, ok); err != nil {
				t.Error(err)
			}
		})
	}
}
