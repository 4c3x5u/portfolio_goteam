package register

// ReqBody defines the request body for the register route.
type ReqBody struct {
	Username Username `json:"username"`
	Password Password `json:"password"`
}

// Validate uses individual field validation logic defined in the validation.go
// file to validate requests sent the register route. It returns an errors
// object if any of the individual validations fail.
func (r *ReqBody) Validate() *ErrsValidation {
	errs := &ErrsValidation{}

	errs.Username = r.Username.Validate()
	errs.Password = r.Password.Validate()

	if len(errs.Username) > 0 || len(errs.Password) > 0 {
		return errs
	}
	return nil
}
