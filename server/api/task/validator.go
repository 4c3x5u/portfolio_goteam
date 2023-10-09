package task

import (
	"server/api"
)

// TitleValidator can be used to validate a task title.
type TitleValidator struct{}

// NewTitleValidator creates and returns a new TitleValidator.
func NewTitleValidator() TitleValidator { return TitleValidator{} }

// Validate validates a given task title.
func (v TitleValidator) Validate(title string) error {
	if title == "" {
		return api.ErrValueEmpty
	}
	if len(title) > 50 {
		return api.ErrValueTooLong
	}
	return nil
}

// IDValidator can be used to validate a board ID.
type IDValidator struct{}

// NewIDValidator creates and returns a new IDValidator.
func NewIDValidator() IDValidator { return IDValidator{} }

// Validate validates a given board ID.
func (i IDValidator) Validate(id string) error {
	if id == "" {
		return api.ErrValueEmpty
	}
	return nil
}
