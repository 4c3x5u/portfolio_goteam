package task

import (
	"server/api"
	"strconv"
)

// TitleValidator can be used to validate a task title.
type TitleValidator struct{}

// NewTitleValidator creates and returns a new TitleValidator.
func NewTitleValidator() TitleValidator { return TitleValidator{} }

// Validate validates a given task title.
func (v TitleValidator) Validate(title string) error {
	if title == "" {
		return api.ErrStrEmpty
	}
	if len(title) > 50 {
		return api.ErrStrTooLong
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
		return api.ErrStrEmpty
	}
	if _, err := strconv.Atoi(id); err != nil {
		return api.ErrStrNotInt
	}
	return nil
}
