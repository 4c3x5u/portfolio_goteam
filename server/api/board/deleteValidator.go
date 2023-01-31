package board

import (
	"net/url"
	"strconv"
)

// POSTValidator describes a type that can be used to validate the URL query
// parameters sent to the board route on DELETE requests.
type DELETEReqValidator interface {
	Validate(url.Values) (string, bool)
}

// POSTValidator can be used to validate the URL query parameters sent to the
// board route on DELETE requests.
type DELETEValidator struct{}

// NewDELETEValidator creates and returns a new DELETEValidator.
func NewDELETEValidator() DELETEValidator { return DELETEValidator{} }

// Validate validates the URL query parameters sent to the board route on DELETE
// requests.
func (v DELETEValidator) Validate(qParams url.Values) (string, bool) {
	idStr := qParams.Get("id")
	if idStr == "" {
		return "", false
	}
	if _, err := strconv.Atoi(idStr); err != nil {
		return "", false
	}
	return idStr, true
}
