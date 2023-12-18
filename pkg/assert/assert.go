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
func newErr(got, want any) error {
	return fmt.Errorf("\ngot: %+v\nwant: %+v", got, want)
}

// Equal asserts that two given values are equal.
func Equal(logErr func(...any), got, want any) {
	if want != got {
		logErr(newErr(got, want))
	}
}

// AllEqual asserts that two given arrays are the same by comparing their
// children.
func AllEqual[T comparable](logErr func(...any), got, want []T) {
	if got == nil && want == nil {
		return
	}
	if len(got) != len(want) {
		logErr(newErr(got, want))
		return
	}
	for i := 0; i < len(want); i++ {
		if got[i] != want[i] {
			logErr(newErr(got, want))
			return
		}
	}
}

// ErrIs asserts that the given two errors are the same.
func ErrIs(logErr func(...any), got, want error) {
	if !errors.Is(got, want) {
		logErr(newErr(got, want))
	}
}

// Nil asserts that a given value is nil.
func Nil(logErr func(...any), got any) {
	if got != nil {
		logErr(newErr(got, "<nil>"))
	}
}

// True asserts that a given boolean value is true.
func True(logErr func(...any), got bool) {
	if !got {
		logErr(newErr(got, "true"))
	}
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
	wantErrMsg string,
) func(*testing.T, *http.Response, string) {
	return func(t *testing.T, res *http.Response, _ string) {
		var resBody map[string]string
		if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
			t.Fatal(err)
		}
		Equal(t.Error, resBody["error"], wantErrMsg)
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
func OnLoggedErr(wantErrMsg string) func(*testing.T, *http.Response, string) {
	return func(t *testing.T, _ *http.Response, loggedErr string) {
		Equal(t.Error, loggedErr, wantErrMsg)
	}
}
