package register

// Req is the request body type for the register endpoint.
type Req struct {
	Username Username `json:"username"`
	Password string   `json:"password"`
	Referrer string   `json:"referrer"`
}

// IsValid performs input validation checks on the register request and returns
// an error if any fails.
func (r *Req) IsValid() (bool, *Errs) {
	errs := &Errs{}

	if usnIsValid, usnErrMsg := r.Username.IsValid(); !usnIsValid {
		errs.Username = append(errs.Username, usnErrMsg)
		return false, errs
	}

	return true, nil
}
