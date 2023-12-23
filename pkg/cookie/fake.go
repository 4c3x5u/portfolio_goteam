package cookie

import (
	"net/http"
)

// FakeDecode has a Func method that can be used as a test fake for DecodeFunc.
type FakeDecoder[T any] struct {
	Res T
	Err error
}

// Func discards the input parameters and returns FakeDecode's Res and
// Err field values.
func (f *FakeDecoder[T]) Decode(http.Cookie) (T, error) {
	return f.Res, f.Err
}

// FakeEncodeAuth has a Func method that can be used as a test fake for
// DecodeInvite.
type FakeEncoder[T any] struct {
	Res http.Cookie
	Err error
}

// Func discards the input parameters and returns FakeEncodeAuth's Res and
// Err field values.
func (f *FakeEncoder[T]) Encode(T) (http.Cookie, error) {
	return f.Res, f.Err
}
