package column

import (
	"errors"
	"strconv"
)

// IDValidator can be used to validate a column ID.
type IDValidator struct{}

// NewIDValidator creates and returns a new IDValidator.
func NewIDValidator() IDValidator { return IDValidator{} }

// Validate validates a given column ID.
func (i IDValidator) Validate(columnID string) error {
	if columnID == "" {
		return errors.New("Column ID cannot be empty.")
	}
	if _, err := strconv.Atoi(columnID); err != nil {
		return errors.New("Column ID must be an integer.")
	}
	return nil
}
