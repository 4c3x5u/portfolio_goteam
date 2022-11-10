package register

// ReqBody defines the request body for the register route.
type ReqBody struct {
	Username Username `json:"username"`
	Password Password `json:"password"`
	Referrer string   `json:"referrer"`
}

// Validate uses individual field validation logic defined in the validation.go
// file to validate requests sent the register route. It returns an errors
// object if any of the individual validations fail.
func (r *ReqBody) Validate() *Errs {
	errs := &Errs{}

	errs.Username = r.Username.Validate()
	errs.Password = r.Password.Validate()

	if errs.Username != nil || errs.Password != nil {
		return errs
	}
	return nil
}
