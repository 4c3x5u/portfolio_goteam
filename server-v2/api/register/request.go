package register

import "regexp"

// Req is the request body type for the register endpoint.
type Req struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Referrer string `json:"referrer"`
}

// IsValid performs input validation checks on the register request and returns
// an error if any fails.
func (r *Req) IsValid() (bool, *Errs) {
	errs := &Errs{}

	// username too short
	if len(r.Username) < 5 {
		errs.Username = append(errs.Username, "Username cannot be shorter than 5 characters.")
		return true, errs
	}

	// username too long
	if len(r.Username) > 15 {
		errs.Username = append(errs.Username, "Username cannot be longer than 15 characters.")
		return true, errs
	}

	// username contains invalid characters
	if match, _ := regexp.MatchString("[^A-Za-z0-9]+", r.Username); match {
		errs.Username = append(errs.Username, "Username can contain only letters and digits.")
		return true, errs
	}

	// username starts with a digit
	if match, _ := regexp.MatchString("(^\\d)", r.Username); match {
		errs.Username = append(errs.Username, "Username can start only with a letter.")
		return true, errs
	}

	return true, nil
}
