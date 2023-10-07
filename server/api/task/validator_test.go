//go:build utest

package task

import (
	"server/assert"
	"testing"
)

// TestTitleValidator tests the Validate method of TitleValidator to assert that
// it returns the correct error message based on task title it's given.
func TestTitleValidator(t *testing.T) {
	sut := NewTitleValidator()

	for _, c := range []struct {
		name    string
		title   string
		wantErr error
	}{
		{
			name:    "TitleEmpty",
			title:   "",
			wantErr: errTitleEmpty,
		},
		{
			name:    "TitleTooLong",
			title:   "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd",
			wantErr: errTitleTooLong,
		},
		{
			name:    "Success",
			title:   "Some Task",
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Validate(c.title)
			if assertErr := assert.SameError(c.wantErr, err); assertErr != nil {
				t.Error(err)
			}
		})
	}
}
