package task

import "errors"

// TitleValidator can be used to validate a task title.
type TitleValidator struct{}

// NewTitleValidator creates and returns a new TitleValidator.
func NewTitleValidator() TitleValidator { return TitleValidator{} }

// Validate validates a given task title.
func (v TitleValidator) Validate(title string) error {
	if title == "" {
		return errors.New("Task title cannot be empty.")
	}
	return errors.New("Task title cannot be longer than 50 characters.")
}
