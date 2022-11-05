package register

import "regexp"

// Username defines the username field of a register request.
type Username string

// Validate applies username validation rules to the Username string and returns
// the error message if any fails.
func (u *Username) Validate() string {
	if *u == "" {
		return "Username cannot be empty."
	}
	if len(*u) < 5 {
		return "Username cannot be shorter than 5 characters."
	}
	if len(*u) > 15 {
		return "Username cannot be longer than 15 characters."
	}
	if match, _ := regexp.MatchString("[^A-Za-z0-9]+", string(*u)); match {
		return "Username can contain only letters and digits."
	}
	if match, _ := regexp.MatchString("(^\\d)", string(*u)); match {
		return "Username can start only with a letter."
	}
	return ""
}

// Password defines the password field of a register request.
type Password string

// Validate applies password validation rules to the Password string and returns
// the error message if any fails.
func (p *Password) Validate() string {
	if *p == "" {
		return "Password cannot be empty."
	}
	if len(*p) < 8 {
		return "Password cannot be shorter than 5 characters."
	}
	if len(*p) > 64 {
		return "Password cannot be longer than 64 characters."
	}
	if match, _ := regexp.MatchString("[a-z]", string(*p)); !match {
		return "Password must contain a lowercase letter (a-z)."
	}
	if match, _ := regexp.MatchString("[A-Z]", string(*p)); !match {
		return "Password must contain an uppercase letter (A-Z)."
	}
	if match, _ := regexp.MatchString("[0-9]", string(*p)); !match {
		return "Password must contain a digit (0-9)."
	}
	return ""
}
