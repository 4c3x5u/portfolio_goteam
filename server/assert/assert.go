//go:build utest || itest

// Package assert contains simple helper functions for test assertions. Its main
// purpose is to centralise the formatting of the error messages for assertions
// and to provide easy-to-read/use abstractions for commonly used assertions.
package assert

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"testing"
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
	// Order is reversed when passing to errors.Is since it is received as
	// "want" first and "got" second, yet we assert on whether "got" is "want".
	if !errors.Is(errB, errA) {
		return newErr(errA, errB)
	}
	return nil
}

// OnResErr can be used in HTTP tests to assert that a given error message was
// written on the response body's "error" field. It takes in the expected error
// message and returns a function that takes in:
//   - *testing.T to be able to either call Fatal or Error,
//   - *http.Response to read the response body,
//   - *pkgLog.FakeErrorer to match the signature of OnLoggedErr so that it can
//     be used interchangeably with it in table-driven tests.
//
// This two-step function call is necessary to be able to use it effectively in
// table-driven tests.
func OnResErr(
	errMsg string,
) func(*testing.T, *http.Response, string) {
	return func(t *testing.T, res *http.Response, _ string) {
		var resBody map[string]string
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Fatal(err)
		}
		if err := Equal(errMsg, resBody["error"]); err != nil {
			t.Error(err)
		}
	}
}

// OnLoggedErr can be used in HTTP tests to assert that a given error message
// was logged. It takes in the expected error message and returns a function
// that takes in:
//   - *testing.T to be able to either call Fatal or Error,
//   - *http.Response to match the signature of OnResErr so that it can be used
//     interchangeably with it in table-driven tests,
//   - *pkgLog.FakeErrorer to check what error was logged.
//
// This two-step function call is necessary to be able to use it effectively in
// table-driven tests.
func OnLoggedErr(
	errMsg string,
) func(*testing.T, *http.Response, string) {
	return func(t *testing.T, _ *http.Response, loggedErr string) {
		if err := Equal(errMsg, loggedErr); err != nil {
			t.Error(err)
		}
	}
}
