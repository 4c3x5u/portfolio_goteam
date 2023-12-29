// Package validator contains code reused by validators.
package validator

import (
	"errors"
)

var (
	// ErrEmpty means that the value being validated was empty.
	ErrEmpty = errors.New("empty")

	// ErrTooLong means that the value being validated was too long.
	ErrTooLong = errors.New("too long")

	// ErrWrongFormat means that the format of the value being validated was
	// wrong.
	ErrWrongFormat = errors.New("invalid format")

	// ErrOutOfBounds means that the value baing validated was outside the valid
	// bounds.
	ErrOutOfBounds = errors.New("out of bounds")
)

// String describes a type that can be used to validate a string.
type String interface{ Validate(string) error }

// Int describes a type that can be used to validate an integer.
type Int interface{ Validate(int) error }

// Func describes a function that can be used to validate a value.
type Func[T any] func(T) error
