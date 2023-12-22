//go:build utest

package board

import (
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
)

func TestNameValidator(t *testing.T) {
	sut := NewNameValidator()

	for _, c := range []struct {
		name      string
		boardName string
		wantErr   error
	}{
		{
			name:      "Empty",
			boardName: "",
			wantErr:   ErrEmpty,
		},
		{
			name:      "TooLong",
			boardName: "boardyboardsyboardkyboardishboardxyz",
			wantErr:   ErrTooLong,
		},
		{
			name:      "OK",
			boardName: "My Board",
			wantErr:   nil,
		},
	} {
		err := sut.Validate(c.boardName)

		assert.ErrIs(t.Error, err, c.wantErr)
	}
}

func TestIDValidator(t *testing.T) {
	sut := NewIDValidator()

	for _, c := range []struct {
		name    string
		boardID string
		wantErr error
	}{
		{
			name:    "Empty",
			boardID: "",
			wantErr: ErrEmpty,
		},
		{
			name:    "NotUUID",
			boardID: "21",
			wantErr: ErrNotUUID,
		},
		{
			name:    "Success",
			boardID: "97377e55-5a2a-4172-bf5d-354b40aa2735",
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Validate(c.boardID)

			assert.ErrIs(t.Error, err, c.wantErr)
		})
	}
}
