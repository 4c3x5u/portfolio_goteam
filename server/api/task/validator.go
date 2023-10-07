package task

import "errors"

var (
	// errTitleEmpty is returned from TitleValidator.Validate when the
	// given title is empty.
	errTitleEmpty = errors.New("title empty")

	// errTitleTooLong is returned from TitleValidator.Validate when the
	// given title is longer than 50 characters.
	errTitleTooLong = errors.New("title too long")
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
	if len(title) > 50 {
		return errTitleTooLong
	}
	return nil
}
