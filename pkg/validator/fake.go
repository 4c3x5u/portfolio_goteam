// go:build utest || itest

package validator

// FakeString is a test fake for String.
type FakeString struct{ Err error }

// Validate discards the given param and returns the fake's Err field value.
func (f FakeString) Validate(string) error { return f.Err }
