package register

import "regexp"

// ValidationErrs defines the validation errors returned in POSTResp.
type ValidationErrs struct {
	Username []string `json:"username,omitempty"`
	Password []string `json:"password,omitempty"`
	TeamID   string   `json:"teamID,omitEmpty"`
}

// Any checks whether there are any validation errors within the ValidationErrors.
func (e ValidationErrs) Any() bool {
	return len(e.Username) > 0 || len(e.Password) > 0
}

// ReqValidator describes a type that validates a request body and returns
// validation errors that occur.
type ReqValidator interface{ Validate(PostReq) ValidationErrs }

// UserValidator is the ReqValidator for the register route.
type UserValidator struct {
	UsernameValidator StrValidator
	PasswordValidator StrValidator
	TeamIDValidator   StrValidator
}

// NewUserValidator creates and returns a new UserValidator.
func NewUserValidator(usnVdtor, pwdVdtor StrValidator) UserValidator {
	return UserValidator{
		UsernameValidator: usnVdtor,
		PasswordValidator: pwdVdtor,
	}
}

// Validate uses UsernameValidator and PasswordValidator to validate requests
// sent the register route. It returns an errors object if any of the individual
// validations fail. It implements the UserValidator interface on the
// ReqValidator struct.
func (v UserValidator) Validate(req PostReq) ValidationErrs {
	errs := ValidationErrs{
		Username: v.UsernameValidator.Validate(req.Username),
		Password: v.PasswordValidator.Validate(req.Password),
	}
	// team ID can be empty
	if req.TeamID != "" {
		errsTID := v.TeamIDValidator.Validate(req.TeamID)
		if len(errsTID) > 0 {
			errs.TeamID = errsTID[0]
		}
	}
	return errs
}

// StrValidator describes a type that validates a string arg and returns a
// string slice containing validation error messages.
type StrValidator interface{ Validate(string) (errs []string) }

// IDValidator is the ID field validator for POST register requests.
type IDValidator struct{}

// NewUsernameValidator creates and returns a new username validator.
func NewUsernameValidator() IDValidator { return IDValidator{} }

// Validate applies user ID validation rules to the ID string and returns the
// error message if any fails.
func (v IDValidator) Validate(id string) (errs []string) {
	if id == "" {
		errs = append(errs, "Username cannot be empty.")
		// if password empty, further validation is pointless – return errors
		return
	} else if len([]rune(id)) < 5 {
		errs = append(errs, "Username cannot be shorter than 5 characters.")
	} else if len([]rune(id)) > 15 {
		errs = append(errs, "Username cannot be longer than 15 characters.")
	}

	if match, _ := regexp.MatchString("[^A-Za-z0-9]+", id); match {
		errs = append(
			errs,
			"Username can contain only letters (a-z/A-Z) and digits (0-9).",
		)
	}
	if match, _ := regexp.MatchString("(^\\d)", id); match {
		errs = append(errs, "Username can start only with a letter (a-z/A-Z).")
	}

	return
}

// PwdValidator is the password field validator for the register route.
type PwdValidator struct{}

// NewPasswordValidator creates and returns a new PasswordValidator.
func NewPasswordValidator() PwdValidator { return PwdValidator{} }

// Validate applies password validation rules to the Password string and returns
// the error message if any fails.
func (v PwdValidator) Validate(pwd string) (errs []string) {
	if pwd == "" {
		errs = append(errs, "Password cannot be empty.")
		// if password empty, further validation is pointless
		return
	} else if len([]rune(pwd)) < 8 {
		errs = append(errs, "Password cannot be shorter than 8 characters.")
	} else if len([]rune(pwd)) > 64 {
		errs = append(errs, "Password cannot be longer than 64 characters.")
	}

	if match, _ := regexp.MatchString("[a-z]", pwd); !match {
		errs = append(errs, "Password must contain a lowercase letter (a-z).")
	}
	if match, _ := regexp.MatchString("[A-Z]", pwd); !match {
		errs = append(errs, "Password must contain an uppercase letter (A-Z).")
	}
	if match, _ := regexp.MatchString("[0-9]", pwd); !match {
		errs = append(errs, "Password must contain a digit (0-9).")
	}
	if match, _ := regexp.MatchString(
		"[!\"#$%&'()*+,-./:;<=>?[\\]^_`{|}~]", pwd,
	); !match {
		errs = append(
			errs,
			"Password must contain one of the following special characters: "+
				"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ _ ` { | "+
				"} ~.",
		)
	}
	if match, _ := regexp.MatchString("\\s", pwd); match {
		errs = append(errs, "Password cannot contain spaces.")
	}
	if match, _ := regexp.MatchString("[^\\x00-\\x7F]", pwd); match {
		errs = append(
			errs,
			"Password can contain only letters (a-z/A-Z), digits (0-9), and "+
				"the following special characters: ! \" # $ % & ' ( ) * + , "+
				"- . / : ; < = > ? [ \\ ] ^ _ ` { | } ~.",
		)
	}

	return
}
