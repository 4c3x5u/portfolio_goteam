package board

import "errors"

// POSTValidator is the ReqValidator for POSTHandler.
type POSTValidator struct{}

// NewPOSTValidator creates and returns a new POSTValidator.
func NewPOSTValidator() POSTValidator { return POSTValidator{} }

// Validate validates the request body sent to the login route.
func (v POSTValidator) Validate(reqBody POSTReqBody) error {
	if reqBody.Name == "" {
		return errNameEmpty
	}
	if len(reqBody.Name) > 35 {
		return errNameTooLong
	}
	return nil
}

var (
	// errNameEmpty is the error returned from the POSTValidator when the
	// received board name is empty.
	errNameEmpty = errors.New("Board name cannot be empty.")

	// errNameTooLong is the error returned from the POSTValidator when the
	// received board name is too long.
	errNameTooLong = errors.New(
		"Board name cannot be longer than 35 characters.",
	)
)
