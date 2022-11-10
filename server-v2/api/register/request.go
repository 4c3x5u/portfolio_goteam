package register

// ReqBody defines the request body for the register route.
type ReqBody struct {
	Username Username `json:"username"`
	Password Password `json:"password"`
	Referrer string   `json:"referrer"`
}

// Validate uses individual field validation logic defined in the validation.go
// file to validate requests sent the register route. It returns false and
// an errors object if any of the individual validations fail.
func (r *ReqBody) Validate() *Errs {
	errs := &Errs{}

	// validate username
	if errsUsername := r.Username.Validate(); errsUsername != nil {
		errs.Username = errsUsername
		return errs
	}

	// validate password
	if errsPassword := r.Password.Validate(); errsPassword != nil {
		errs.Password = errsPassword
		return errs
	}

	return nil
}
