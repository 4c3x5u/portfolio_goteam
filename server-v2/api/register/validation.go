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
