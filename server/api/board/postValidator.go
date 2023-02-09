package board

// POSTReqValidator describes a type that can be used to validate the body of
// POST requests sent to the board route.
type POSTReqValidator interface{ Validate(POSTReqBody) string }

// POSTValidator can be used to validate the body of POST requests sent to the
// board route.
type POSTValidator struct{}

// NewPOSTValidator creates and returns a new POSTValidator.
func NewPOSTValidator() POSTValidator { return POSTValidator{} }

// Validate validates the body of the POST request sent to the board route.
func (v POSTValidator) Validate(reqBody POSTReqBody) string {
	if reqBody.Name == "" {
		return msgNameEmpty
	}
	if len(reqBody.Name) > 35 {
		return msgNameTooLong
	}
	return ""
}

// msgNameEmpty is the error returned from the POSTValidator when the
// received board name is empty.
const msgNameEmpty = "Board name cannot be empty."

// msgNameTooLong is the error returned from the POSTValidator when the
// received board name is too long.
const msgNameTooLong = "Board name cannot be longer than 35 characters."
