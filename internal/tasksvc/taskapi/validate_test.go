//go:build utest

package taskapi

import (
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/validator"
)

// TestTitleValidator tests the TitleValidator.Validate method.
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
			wantErr: validator.ErrEmpty,
		},
		{
			name:    "TitleTooLong",
			title:   "asdqweasdqweasdqweasdqweasdqweasdqweasdqweasdqweasd",
			wantErr: validator.ErrTooLong,
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

// TestColNoValidator tests the ColNoValidator.Validate method.
func TestColNoValidator(t *testing.T) {
	sut := NewColNoValidator()

	for _, c := range []struct {
		name    string
		colNo   int
		wantErr error
	}{
		{
			name:    "ColNoTooSmall",
			colNo:   -1,
			wantErr: validator.ErrOutOfBounds,
		},
		{
			name:    "ColNoTooBig",
			colNo:   4,
			wantErr: validator.ErrOutOfBounds,
		},
		{
			name:    "Success",
			colNo:   2,
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Validate(c.colNo)
			assert.ErrIs(t.Error, err, c.wantErr)
		})
	}
}
