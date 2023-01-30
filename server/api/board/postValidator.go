package board

import "errors"

// POSTReqValidator describes a type that can be used to validate the body of
// POST requests sent to the board route.
type POSTReqValidator interface{ Validate(POSTReqBody) error }

// POSTValidator can be used to validate the body of POST requests sent to the
// board route.
type POSTValidator struct{}

// NewPOSTValidator creates and returns a new POSTValidator.
func NewPOSTValidator() POSTValidator { return POSTValidator{} }

// Validate validates the body of the POST request sent to the board route.
func (v POSTValidator) Validate(reqBody POSTReqBody) error {
	if reqBody.Name == "" {
		return errNameEmpty
	}
	if len(reqBody.Name) > 35 {
		return errNameTooLong
	}
	return nil
}

// errNameEmpty is the error returned from the POSTValidator when the
// received board name is empty.
var errNameEmpty = errors.New("Board name cannot be empty.")

// errNameTooLong is the error returned from the POSTValidator when the
// received board name is too long.
var errNameTooLong = errors.New(
	"Board name cannot be longer than 35 characters.",
)
