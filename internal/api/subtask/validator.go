package subtask

import (
	"strconv"

	"github.com/kxplxn/goteam/internal/api"
)

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
