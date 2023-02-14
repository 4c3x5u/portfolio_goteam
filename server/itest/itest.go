//go:build itest

// Package itest contains integration tests for the application. This module is
// intended to be run independently of any actual project code. Therefore, some
// code from the other modules are copied over for convenience.
package itest

import "fmt"

// serverURL is the url that is used to send requests to the server running in
// the test container. It is used during setup in main_test.go/TestMain.
const serverURL = "http://localhost:8081"

// dbConnStr is the connection string for the test db.
const dbConnStr = "postgres://itestuser:itestpwd@localhost:5433/itestdb" +
	"?sslmode=disable"

// assertEqual asserts that two given values are equal.
func assertEqual(want, got any) error {
	if want != got {
		return newAssertionErr(want, got)
	}
	return nil
}

// assertEqualArr asserts that two given arrays are the same by comparing their
// children.
func assertEqualArr[T comparable](want, got []T) error {
	if want == nil && got == nil {
		return nil
	}
	if len(want) != len(got) {
		return newAssertionErr(want, got)
	}
	for i := 0; i < len(want); i++ {
		if want[i] != got[i] {
			return newAssertionErr(want, got)
		}
	}
	return nil
}

// newAssertionErr formats, creates, and returns an assertion error.
func newAssertionErr(want, got any) error {
	return fmt.Errorf("\nwant: %+v\ngot: %+v", want, got)
}
