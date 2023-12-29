package taskapi

import (
	"errors"

	"github.com/google/uuid"

	"github.com/kxplxn/goteam/pkg/validator"
)

// TitleValidator can be used to validate a task title.
type TitleValidator struct{}

// NewTitleValidator creates and returns a new TitleValidator.
func NewTitleValidator() TitleValidator { return TitleValidator{} }

// Validate validates a given task title.
func (v TitleValidator) Validate(title string) error {
	if title == "" {
		return validator.ErrEmpty
	}
	if len(title) > 50 {
		return validator.ErrTooLong
	}
	return nil
}

// ColNoValidator can be used to validate a task's column number.
type ColNoValidator struct{}

// ValidatePostReq validates a given PostReq.
func ValidatePostReq(req PostReq) error {
	if req.BoardID == "" {
		return errBoardIDEmpty
	}
	if _, err := uuid.Parse(req.BoardID); err != nil {
		return errParseBoardID
	}
	if req.ColNo < 1 || req.ColNo > 4 {
		return errColNoOutOfBounds
	}
	if req.Title == "" {
		return errTitleEmpty
	}
	if len(req.Title) > 50 {
		return errTitleTooLong
	}
	if len(req.Description) > 500 {
		return errDescTooLong
	}
	for _, st := range req.Subtasks {
		if st.Title == "" {
			return errSubtaskTitleEmpty
		}
		if len(st.Title) > 50 {
			return errSubtaskTitleTooLong
		}
	}
	if req.Order < 0 {
		return errOrderNegative
	}
	return nil
}

var (
	// errBoardIDEmpty is returned when a board ID is empty.
	errBoardIDEmpty = errors.New("board id is empty")

	// errParseBoardID is returned when a board ID cannot be parsed.
	errParseBoardID = errors.New("could not parse board id")

	// errColNoOutOfBounds is returned when a column number is out of bounds.
	errColNoOutOfBounds = errors.New("column number is out of bounds")

	// errTitleEmpty is returned when a task title is empty.
	errTitleEmpty = errors.New("title is empty")

	// errTitleTooLong is returned when a task title is too long.
	errTitleTooLong = errors.New("title is too long")

	// errTitleEmpty is returned when a task description is too long.
	errDescTooLong = errors.New("description is too long")

	// errSubtaskTitleEmpty is returned when a subtask is empty.
	errSubtaskTitleEmpty = errors.New("subtask is empty")

	// errSubtaskTitleTooLong is returned when a subtask is too long.
	errSubtaskTitleTooLong = errors.New("subtask is too long")

	// errOrderNegative is returned when the order is negative.
	errOrderNegative = errors.New("order is negative")
)
