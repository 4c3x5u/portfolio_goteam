// Package assert contains helper functions for commonly used assertions in
// tests.
package assert

import "testing"

func Equal[T comparable](t *testing.T, want, got T) {
	if want != got {
		t.Logf("\nwant: %+v\ngot: %+v", want, got)
		t.Fail()
	}
}

// EqualArr compares two arrays and returns a boolean based on whether or not
// they contain exactly the same items.
func EqualArr[T comparable](t *testing.T, want, got []T) {
	areEqual := func(want []T, got []T) bool {
		if want == nil && got == nil {
			return true
		}
		if len(want) != len(got) {
			return false
		}
		for i := 0; i < len(want); i++ {
			if want[i] != got[i] {
				return false
			}
		}
		return true
	}

	if !areEqual(want, got) {
		t.Logf("\nwant: %+v\ngot: %+v", want, got)
		t.Fail()
	}
}

// Nil asserts that a given "got" is nil.
func Nil(t *testing.T, got any) {
	if got != nil {
		t.Logf("%v+ is not nil", got)
		t.Fail()
	}
}
