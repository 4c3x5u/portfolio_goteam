// Package assert contains simple helper functions for test assertions. Its main
// purpose is to centralise the formatting of the error messages for assertions
// and to provide easy-to-read/use abstractions for commonly used assertions.
package assert

import (
	"errors"
	"fmt"
)

// newErr formats, creates, and returns an assertion error.
func newErr(want, got any) error {
	return fmt.Errorf("\nwant: %+v\ngot: %+v", want, got)
}

// Equal asserts that two given values are equal.
func Equal(want, got any) error {
	if want != got {
		return newErr(want, got)
	}
	return nil
}

// EqualArr asserts that two given arrays are the same by comparing their
// children.
func EqualArr[T comparable](want, got []T) error {
	if want == nil && got == nil {
		return nil
	}
	if len(want) != len(got) {
		return newErr(want, got)
	}
	for i := 0; i < len(want); i++ {
		if want[i] != got[i] {
			return newErr(want, got)
		}
	}
	return nil
}

// Nil asserts that a given value is nil.
func Nil(got any) error {
	if got != nil {
		return newErr("<nil>", got)
	}
	return nil
}

// True asserts that a given boolean value is true.
func True(got bool) error {
	if !got {
		return newErr("true", got)
	}
	return nil
}

// SameError asserts that the given two errors are the same.
func SameError(errA, errB error) error {
	if !errors.Is(errA, errB) {
		return newErr(errA, errB)
	}
	return nil
}
