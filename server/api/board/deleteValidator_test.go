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
		qParams url.Values
		wantID  string
		wantErr error
	}{
		{
			name:    "NoBoardID",
			qParams: url.Values{},
			wantID:  "",
			wantErr: errEmptyBoardID,
		},
		{
			name:    "EmptyBoardID",
			qParams: url.Values{"id": []string{""}},
			wantID:  "",
			wantErr: errEmptyBoardID,
		},
		{
			name:    "Valid",
			qParams: url.Values{"id": []string{"12"}},
			wantID:  "12",
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			id, err := sut.Validate(c.qParams)

			if err = assert.Equal(c.wantErr, err); err != nil {
				t.Error(err)
			}
			if err = assert.Equal(c.wantID, id); err != nil {
				t.Error(err)
			}
		})
	}
}
