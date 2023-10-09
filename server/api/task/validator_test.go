//go:build utest

package task

import (
	"server/api"
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
			wantErr: api.ErrValueEmpty,
		},
		{
			name:    "TitleTooLong",
			title:   "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd",
			wantErr: api.ErrValueTooLong,
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
			wantErr: api.ErrValueEmpty,
		},
		{
			name:    "Success",
			id:      "12",
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Validate(c.id)
			if assertErr := assert.SameError(c.wantErr, err); assertErr != nil {
				t.Error(assertErr)
			}
		})
	}
}
