package login

// RequestValidator describes a type that can be used to validate requests sent
// to the login route.
type RequestValidator interface{ Validate(ReqBody) bool }

// Validator is the RequestValidator for the login route.
type Validator struct{}

func NewValidator() Validator { return Validator{} }

// Validate validates the request body sent to the login route.
func (v Validator) Validate(reqBody ReqBody) bool {
	if reqBody.Username == "" || reqBody.Password == "" {
		return false
	}
	return true
}
