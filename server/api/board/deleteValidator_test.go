//go:build utest

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
		wantOK  bool
	}{
		{
			name:    "NoBoardID",
			qParams: url.Values{},
			wantID:  "",
			wantOK:  false,
		},
		{
			name:    "EmptyBoardID",
			qParams: url.Values{"id": []string{""}},
			wantID:  "",
			wantOK:  false,
		},
		{
			name:    "NotInteger",
			qParams: url.Values{"id": []string{"My Board"}},
			wantID:  "",
			wantOK:  false,
		},
		{
			name:    "Valid",
			qParams: url.Values{"id": []string{"12"}},
			wantID:  "12",
			wantOK:  true,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			id, ok := sut.Validate(c.qParams)

			if err := assert.Equal(c.wantOK, ok); err != nil {
				t.Error(err)
			}
			if err := assert.Equal(c.wantID, id); err != nil {
				t.Error(err)
			}
		})
	}
}
