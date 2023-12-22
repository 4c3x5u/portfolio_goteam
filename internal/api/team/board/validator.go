package board

import (
	"errors"

	"github.com/google/uuid"
)

var (
	// ErrEmpty is returned from Validate methods when the given input is empty.
	ErrEmpty = errors.New("input is empty")

	// ErrTooLong is returned from NameValidator.Validate when the given
	// board name is longer than 35 characters.
	ErrTooLong = errors.New("input is too long")

	// ErrNotUUID is returned from IDValidator.Validate when the given board ID
	// is not a valid UUID.
	ErrNotUUID = errors.New("input is not a UUID")
)

// NameValidator can be used to validate a board name.
type NameValidator struct{}

// NewNameValidator creates and returns a new NameValidator.
func NewNameValidator() NameValidator { return NameValidator{} }

// Validate validates a given board name.
func (b NameValidator) Validate(boardName string) error {
	if boardName == "" {
		return ErrEmpty
	}
	if len(boardName) > 35 {
		return ErrTooLong
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
		return ErrEmpty
	}
	if _, err := uuid.Parse(boardID); err != nil {
		return ErrNotUUID
	}
	return nil
}
