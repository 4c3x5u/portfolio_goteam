package login

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
