//go:build utest

package board

import (
	"testing"

	"github.com/kxplxn/goteam/server/assert"
)

// TestNameValidator tests the Validate method of NameValidator to assert that
// it returns the correct error message based on board name it's given.
func TestNameValidator(t *testing.T) {
	sut := NewNameValidator()

	msgEmpty := "Board name cannot be empty."

	for _, c := range []struct {
		name       string
		boardName  string
		wantErrMsg string
	}{
		{
			name:       "Empty",
			boardName:  "",
			wantErrMsg: msgEmpty,
		},
		{
			name:       "TooLong",
			boardName:  "boardyboardsyboardkyboardishboardxyz",
			wantErrMsg: "Board name cannot be longer than 35 characters.",
		},
		{
			name:       "OK",
			boardName:  "My Board",
			wantErrMsg: "",
		},
	} {
		err := sut.Validate(c.boardName)

		if c.wantErrMsg == "" {
			if err = assert.Nil(err); err != nil {
				t.Error(err)
			}
		} else {
			assert.Equal(t.Error, err.Error(), c.wantErrMsg)
		}
	}
}

// TestIDValidator tests the Validate method of IDValidator to assert
// that it returns the correct error message based on the board ID it's given.
func TestIDValidator(t *testing.T) {
	sut := NewIDValidator()

	for _, c := range []struct {
		name       string
		boardID    string
		wantErrMsg string
	}{
		{
			name:       "Nil",
			boardID:    "",
			wantErrMsg: "Board ID cannot be empty.",
		},
		{
			name:       "NotInt",
			boardID:    "My Board",
			wantErrMsg: "Board ID must be an integer.",
		},
		{
			name:       "Success",
			boardID:    "12",
			wantErrMsg: "",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Validate(c.boardID)

			if c.wantErrMsg == "" {
				if err = assert.Nil(err); err != nil {
					t.Error(err)
				}
			} else {
				assert.Equal(t.Error, err.Error(), c.wantErrMsg)
			}
		})
	}
}
