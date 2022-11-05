package register

// ReqBody defines the request body for the register endpoint.
type ReqBody struct {
	Username Username `json:"username"`
	Password string   `json:"password"`
	Referrer string   `json:"referrer"`
}

// Validate uses individual field validation logic defined in the validation.go
// file to validate requests sent the register endpoint. It returns false and
// an errors object if any of the individual validations fail.
func (r *ReqBody) Validate() *Errs {
	errs := &Errs{}

	// validate username
	if usnErrMsg := r.Username.Validate(); usnErrMsg != "" {
		errs.Username = usnErrMsg
		return errs
	}

	return nil
}
