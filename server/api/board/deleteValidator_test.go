package board

import (
	"net/url"
	"testing"

	"server/assert"
)

// TestDELETEValidator tests the Validate method of DELETEValidator to assert
// that it returns the correct error based on the URL query parameters.
func TestDELETEValidator(t *testing.T) {
	sut := NewDELETEValidator()

	for _, c := range []struct {
		name    string
		boardID url.Values
		wantErr error
	}{
		{
			name:    "NoBoardID",
			boardID: url.Values{},
			wantErr: errEmptyBoardID,
		},
		{
			name:    "EmptyBoardID",
			boardID: url.Values{"id": []string{""}},
			wantErr: errEmptyBoardID,
		},
		{
			name:    "Valid",
			boardID: url.Values{"id": []string{"12"}},
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Validate(c.boardID)

			if err = assert.Equal(c.wantErr, err); err != nil {
				t.Error(err)
			}
		})
	}
}
