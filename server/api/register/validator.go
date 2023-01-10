package register

import "regexp"

const (
	usnEmpty       = "Username cannot be empty."
	usnTooShort    = "Username cannot be shorter than 5 characters."
	usnTooLong     = "Username cannot be longer than 15 characters."
	usnInvalidChar = "Username can contain only letters (a-z/A-Z) and digits (0-9)."
	usnDigitStart  = "Username can start only with a letter (a-z/A-Z)."

	pwdEmpty     = "Password cannot be empty."
	pwdTooShort  = "Password cannot be shorter than 8 characters."
	pwdTooLong   = "Password cannot be longer than 64 characters."
	pwdNoLower   = "Password must contain a lowercase letter (a-z)."
	pwdNoUpper   = "Password must contain an uppercase letter (A-Z)."
	pwdNoDigit   = "Password must contain a digit (0-9)."
	pwdNoSpecial = "Password must contain one of the following special characters: " +
		"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ _ ` { | } ~."
	pwdHasSpace = "Password cannot contain spaces."
	pwdNonASCII = "Password can contain only letters (a-z/A-Z), digits (0-9), " +
		"and the following special characters: " +
		"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ _ ` { | } ~."
)

// Validator describes a type that validates a request body and returns
// validation errors that occur.
type Validator interface {
	Validate(req ReqBody) (errs ValidationErrs)
}

// RequestValidator is the request validator for the register route.
type RequestValidator struct {
	UsernameValidator StringValidator
	PasswordValidator StringValidator
}

// NewRequestValidator is the constructor for RequestValidator.
func NewRequestValidator(
	usernameValidator, passwordValidator StringValidator,
) RequestValidator {
	return RequestValidator{
		UsernameValidator: usernameValidator,
		PasswordValidator: passwordValidator,
	}
}

// Validate uses UsernameValidator and PasswordValidator to validate requests
// sent the register route. It returns an errors object if any of the individual
// validations fail. It implements the Validator interface on the
// RequestValidator struct.
func (v RequestValidator) Validate(req ReqBody) ValidationErrs {
	return ValidationErrs{
		Username: v.UsernameValidator.Validate(req.Username),
		Password: v.PasswordValidator.Validate(req.Password),
	}
}

// StringValidator describes a type that validates a string input and returns a
// string slice containing validation error messages.
type StringValidator interface {
	Validate(string) (errs []string)
}

// UsernameValidator is the password field validator for the register route.
type UsernameValidator struct{}

// NewUsernameValidator creates and returns a new username validator.
func NewUsernameValidator() UsernameValidator { return UsernameValidator{} }

// Validate applies password validation rules to the Username string and returns
// the error message if any fails.
func (v UsernameValidator) Validate(username string) (errs []string) {
	if username == "" {
		errs = append(errs, "Username cannot be empty.")
		// if password empty, further validation is pointless – return errors
		return
	} else if len([]rune(username)) < 5 {
		errs = append(errs, "Username cannot be shorter than 5 characters.")
	} else if len([]rune(username)) > 15 {
		errs = append(errs, "Username cannot be longer than 15 characters.")
	}

	if match, _ := regexp.MatchString("[^A-Za-z0-9]+", username); match {
		errs = append(errs, "Username can contain only letters (a-z/A-Z) and digits (0-9).")
	}
	if match, _ := regexp.MatchString("(^\\d)", username); match {
		errs = append(errs, "Username can start only with a letter (a-z/A-Z).")
	}

	return
}

// PasswordValidator is the password field validator for the register route.
type PasswordValidator struct{}

// NewPasswordValidator is the constructor for PasswordValidator.
func NewPasswordValidator() PasswordValidator { return PasswordValidator{} }

// Validate applies password validation rules to the Password string and returns
// the error message if any fails.
func (v PasswordValidator) Validate(password string) (errs []string) {
	if password == "" {
		errs = append(errs, "Password cannot be empty.")
		// if password empty, further validation is pointless
		return
	} else if len([]rune(password)) < 8 {
		errs = append(errs, "Password cannot be shorter than 8 characters.")
	} else if len([]rune(password)) > 64 {
		errs = append(errs, "Password cannot be longer than 64 characters.")
	}

	if match, _ := regexp.MatchString("[a-z]", password); !match {
		errs = append(errs, "Password must contain a lowercase letter (a-z).")
	}
	if match, _ := regexp.MatchString("[A-Z]", password); !match {
		errs = append(errs, "Password must contain an uppercase letter (A-Z).")
	}
	if match, _ := regexp.MatchString("[0-9]", password); !match {
		errs = append(errs, "Password must contain a digit (0-9).")
	}
	if match, _ := regexp.MatchString("[!\"#$%&'()*+,-./:;<=>?[\\]^_`{|}~]", password); !match {
		errs = append(errs, "Password must contain one of the following special characters: "+
			"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ _ ` { | } ~.")
	}
	if match, _ := regexp.MatchString("\\s", password); match {
		errs = append(errs, "Password cannot contain spaces.")
	}
	if match, _ := regexp.MatchString("[^\\x00-\\x7F]", password); match {
		errs = append(errs, "Password can contain only letters (a-z/A-Z), digits (0-9), "+
			"and the following special characters: "+
			"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ _ ` { | } ~.",
		)
	}

	return
}
