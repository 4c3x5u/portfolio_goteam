package register

import "regexp"

// Username defines the username field of a register request.
type Username string

// IsValid applies username validation rules to the username string.
func (u *Username) IsValid() (bool, string) {
	// username too short
	if len(*u) < 5 {
		return false, "Username cannot be shorter than 5 characters."
	}

	// username too long
	if len(*u) > 15 {
		return false, "Username cannot be longer than 15 characters."
	}

	// username contains invalid characters
	if match, _ := regexp.MatchString("[^A-Za-z0-9]+", string(*u)); match {
		return false, "Username can contain only letters and digits."
	}

	// username starts with a digit
	if match, _ := regexp.MatchString("(^\\d)", string(*u)); match {
		return false, "Username can start only with a letter."
	}

	return true, ""
}
