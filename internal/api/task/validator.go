package task

import (
	"strconv"

	"github.com/kxplxn/goteam/internal/api"
)

// TitleValidator can be used to validate a task title.
type TitleValidator struct{}

// NewTitleValidator creates and returns a new TitleValidator.
func NewTitleValidator() TitleValidator { return TitleValidator{} }

// Validate validates a given task title.
func (v TitleValidator) Validate(title string) error {
	if title == "" {
		return api.ErrEmpty
	}
	if len(title) > 50 {
		return api.ErrTooLong
	}
	return nil
}

// IDValidator can be used to validate a task ID.
type IDValidator struct{}

// NewIDValidator creates and returns a new IDValidator.
func NewIDValidator() IDValidator { return IDValidator{} }

// Validate validates a given task ID.
func (i IDValidator) Validate(id string) error {
	if id == "" {
		return api.ErrEmpty
	}
	if _, err := strconv.Atoi(id); err != nil {
		return api.ErrNotInt
	}
	return nil
}

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
