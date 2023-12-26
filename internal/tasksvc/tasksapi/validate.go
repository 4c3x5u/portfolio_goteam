package tasksapi

import (
	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/validator"
)

// ColNoValidator can be used to validate a task's column number.
type ColNoValidator struct{}

// NewColNoValidator creates and returns a new ColNoValidator.
func NewColNoValidator() ColNoValidator { return ColNoValidator{} }

// Validate validates a task's column number.
func (v ColNoValidator) Validate(number int) error {
	if number < 0 || number > 3 {
		return validator.ErrOutOfBounds
	}
	return nil
}

// BoardIDValidator can be used to validate a board ID.
type BoardIDValidator struct{}

// NewBoardIDValidator creates and returns a new BoardIDValidator.
func NewBoardIDValidator() BoardIDValidator { return BoardIDValidator{} }

// Validate validates a given board ID.
func (i BoardIDValidator) Validate(boardID string) error {
	if boardID == "" {
		return validator.ErrEmpty
	}
	if _, err := uuid.Parse(boardID); err != nil {
		return validator.ErrWrongFormat
	}
	return nil
}
