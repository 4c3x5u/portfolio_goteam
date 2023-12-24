package boardapi

import (
	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/validator"
)

// NameValidator can be used to validate a board name.
type NameValidator struct{}

// NewNameValidator creates and returns a new NameValidator.
func NewNameValidator() NameValidator { return NameValidator{} }

// Validate validates a given board name.
func (b NameValidator) Validate(boardName string) error {
	if boardName == "" {
		return validator.ErrEmpty
	}
	if len(boardName) > 35 {
		return validator.ErrTooLong
	}
	return nil
}

// IDValidator can be used to validate a board ID.
type IDValidator struct{}

// NewIDValidator creates and returns a new IDValidator.
func NewIDValidator() IDValidator { return IDValidator{} }

// Validate validates a given board ID.
func (i IDValidator) Validate(boardID string) error {
	if boardID == "" {
		return validator.ErrEmpty
	}
	if _, err := uuid.Parse(boardID); err != nil {
		return validator.ErrWrongFormat
	}
	return nil
}
