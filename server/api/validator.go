package api

import "errors"

var (
	// ErrValueEmpty is the error returned from StringValidator or IntValidator
	// when the value passed in is empty.
	ErrValueEmpty = errors.New("value empty")

	// ErrValueTooLong is the error returned from StringValidator when the value
	// passed in is too long.
	ErrValueTooLong = errors.New("value too long")
)

// StringValidator describes a type that can be used to validate a string.
type StringValidator interface{ Validate(string) error }
