package api

import "errors"

var (
	// ErrStrEmpty is the error returned from StringValidator or IntValidator
	// when the value passed in is empty.
	ErrStrEmpty = errors.New("string empty")

	// ErrStrTooLong is the error returned from StringValidator when the value
	// passed in exceeds the maximum length allowed.
	ErrStrTooLong = errors.New("string too long")

	// ErrStrNotInt is the error returned from StringValidator when the value
	// passed in is expected to contain an integer only but does not.
	ErrStrNotInt = errors.New("string not an integer")
)

// StringValidator describes a type that can be used to validate a string.
type StringValidator interface{ Validate(string) error }
