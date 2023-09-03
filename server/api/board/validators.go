package board

import (
	"errors"
	"strconv"
)

// StrValidator describes a type that can be used to validate a string.
type StrValidator interface{ Validate(string) error }

// NameValidator can be used to validate a board name.
type NameValidator struct{}

// NewNameValidator creates and returns a new NameValidator.
func NewNameValidator() NameValidator { return NameValidator{} }

// Validate validates a given board name.
func (b NameValidator) Validate(boardName string) error {
	if boardName == "" {
		return errors.New("Board name cannot be empty.")
	}
	if len(boardName) > 35 {
		return errors.New("Board name cannot be longer than 35 characters.")
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
		return errors.New("Board ID cannot be empty.")
	}
	if _, err := strconv.Atoi(boardID); err != nil {
		return errors.New("Board ID must be an integer.")
	}
	return nil
}
