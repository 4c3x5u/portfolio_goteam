//go:build utest

package cookie

import (
	"net/http"
)

// FakeEncoder is a test fake for Encoder.
type FakeEncoder[T any] struct {
	Res http.Cookie
	Err error
}

// Encode discards the input parameters and returns the FakeEncoder's Res and
// Err field values.
func (f *FakeEncoder[T]) Encode(T) (http.Cookie, error) {
	return f.Res, f.Err
}

// FakeDecoder is a test fake for Decoder.
type FakeDecoder[T any] struct {
	Res T
	Err error
}

// Decode discards the input parameters and returns the FakeDecoder's Res and
// Err field values.
func (f *FakeDecoder[T]) Decode(http.Cookie) (T, error) {
	return f.Res, f.Err
}

// FakeStringDecoder is a test fake for StringDecoder.
type FakeStringDecoder[T any] struct {
	Res T
	Err error
}

// FakeStringDecoder discards the input parameters and returns the
// FakeStringDecoder's Res and Err field values.
func (f *FakeStringDecoder[T]) Decode(string) (T, error) {
	return f.Res, f.Err
}
