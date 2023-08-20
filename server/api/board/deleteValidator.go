package board

import (
	"strconv"
)

// DELETEReqValidator describes a type that can be used to validate the URL query
// parameters sent to the board route on DELETE requests.
type DELETEReqValidator interface{ Validate(string) bool }

// DELETEValidator can be used to validate the URL query parameters sent to the
// board route on DELETE requests.
type DELETEValidator struct{}

// NewDELETEValidator creates and returns a new DELETEValidator.
func NewDELETEValidator() DELETEValidator { return DELETEValidator{} }

// Validate validates the URL query parameters sent to the board route on DELETE
// requests.
func (v DELETEValidator) Validate(id string) bool {
	if id == "" {
		return false
	}
	if _, err := strconv.Atoi(id); err != nil {
		return false
	}
	return true
}
