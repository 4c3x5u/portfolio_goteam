//go:build utest

package task

import (
	"testing"

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/assert"
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
			wantErr: api.ErrEmpty,
		},
		{
			name:    "TitleTooLong",
			title:   "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd",
			wantErr: api.ErrTooLong,
		},
		{
			name:    "Success",
			title:   "Some Task",
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Validate(c.title)
			assert.ErrIs(t.Error, err, c.wantErr)
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
