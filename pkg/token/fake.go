package token

import (
	"time"
)

// FakeDecode has a Func method that can be used as a test fake for DecodeFunc.
type FakeDecode[T any] struct {
	Res T
	Err error
}

// Func discards the input parameters and returns FakeDecode's Res and
// Err field values.
func (f *FakeDecode[T]) Func(_ string) (T, error) {
	return f.Res, f.Err
}

// FakeEncodeAuth has a Func method that can be used as a test fake for
// DecodeInvite.
type FakeEncode[T any] struct {
	Res string
	Err error
}

// Func discards the input parameters and returns FakeEncodeAuth's Res and
// Err field values.
func (f *FakeEncode[T]) Func(_ time.Time, _ T) (string, error) {
	return f.Res, f.Err
}
