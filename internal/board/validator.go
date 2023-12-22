package board

import (
	"errors"
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
