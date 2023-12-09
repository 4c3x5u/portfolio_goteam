package api

import "errors"

var (
	// ErrEmpty is the error returned from StringValidator or IntValidator
	// when the value passed in is empty.
	ErrEmpty = errors.New("empty")

	// ErrTooLong is the error returned from StringValidator when the value
	// passed in exceeds the maximum length allowed.
	ErrTooLong = errors.New("too long")

	// ErrNotInt is the error returned from StringValidator when the value
	// passed in is expected to contain an integer only but does not.
	ErrNotInt = errors.New("not integer")

	// ErrOutOfBounds is returned from IntValidator when the value must be
	// within a given range but is not.
	ErrOutOfBounds = errors.New("out of bounds")
)

// StringValidator describes a type that can be used to validate a string.
type StringValidator interface{ Validate(string) error }

// IntValidator describes a type that can be used to validate an integer.
type IntValidator interface{ Validate(int) error }
