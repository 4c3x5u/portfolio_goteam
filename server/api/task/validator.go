package task

import "errors"

var (
	errTitleEmpty   error = errors.New("title empty")
	errTitleTooLong error = errors.New("title too long")
)

// TitleValidator can be used to validate a task title.
type TitleValidator struct{}

// NewTitleValidator creates and returns a new TitleValidator.
func NewTitleValidator() TitleValidator { return TitleValidator{} }

// Validate validates a given task title.
func (v TitleValidator) Validate(title string) error {
	if title == "" {
		return errTitleEmpty
	}
	return errTitleTooLong
}
