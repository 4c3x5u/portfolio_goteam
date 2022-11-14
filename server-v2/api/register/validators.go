package register

import "regexp"

// ValidatorReq represents a type that validates a *ReqBody and returns an
// *Errs based on the validation errors that occur.
type ValidatorReq interface {
	Validate(req *ReqBody) (errs *Errs)
}

// Validator is the request validator for the register route.
type Validator struct {
	ValidatorUsername ValidatorStr
	ValidatorPassword ValidatorStr
}

// NewValidator is the constructor for Validator.
func NewValidator(validatorUsername, validatorPassword ValidatorStr) *Validator {
	return &Validator{
		ValidatorUsername: validatorUsername,
		ValidatorPassword: validatorPassword,
	}
}

// Validate uses individual field validation logic defined in the validation.go
// file to validate requests sent the register route. It returns an errors
// object if any of the individual validations fail. It implements the
// ValidatorReq interface on the Validator struct.
func (v *Validator) Validate(req *ReqBody) *Errs {
	errs := &Errs{}
	errs.Username = v.ValidatorUsername.Validate(req.Username)
	errs.Password = v.ValidatorPassword.Validate(req.Password)
	if len(errs.Username) > 0 || len(errs.Password) > 0 {
		return errs
	}
	return nil
}

// ValidatorStr represents a type that validates a string input and returns a
// string slice containing validation error messages.
type ValidatorStr interface {
	Validate(string) (errs []string)
}

// ValidatorUsername is the password field validator for the register route.
type ValidatorUsername struct {
}

// NewValidatorUsername is the constructor for ValidatorUsername.
func NewValidatorUsername() *ValidatorUsername {
	return &ValidatorUsername{}
}

// Validate applies password validation rules to the Username string and returns
// the error message if any fails.
func (v *ValidatorUsername) Validate(username string) (errs []string) {
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

// ValidatorPassword is the password field validator for the register route.
type ValidatorPassword struct {
}

// NewValidatorPassword is the constructor for ValidatorPassword.
func NewValidatorPassword() *ValidatorPassword {
	return &ValidatorPassword{}
}

// Validate applies password validation rules to the Password string and returns
// the error message if any fails.
func (v *ValidatorPassword) Validate(password string) (errs []string) {
	if password == "" {
		errs = append(errs, "Password cannot be empty.")
		// if password empty, further validation is pointless – return wantErrs
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
