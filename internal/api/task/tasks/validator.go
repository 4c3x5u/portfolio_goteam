package tasks

import "github.com/kxplxn/goteam/internal/api"

// ColNoValidator can be used to validate a task's column number.
type ColNoValidator struct{}

// NewColNoValidator creates and returns a new ColNoValidator.
func NewColNoValidator() ColNoValidator { return ColNoValidator{} }

// Validate validates a task's column number.
func (v ColNoValidator) Validate(number int) error {
	if number < 0 || number > 3 {
		return api.ErrOutOfBounds
	}
	return nil
}
