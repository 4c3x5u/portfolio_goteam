package token

import "time"

// FakeDecodeInvite has a Func method that can be used as a test fake for
// DecodeInvite.
type FakeDecodeInvite struct {
	Decoded Invite
	Err     error
}

// Func discards the input parameters and returns FakeDecodeInvite's Inv and
// Err field values.
func (f *FakeDecodeInvite) Func(_ string) (Invite, error) {
	return f.Decoded, f.Err
}

// FakeEncodeAuth has a Func method that can be used as a test fake for
// DecodeInvite.
type FakeEncodeAuth struct {
	Encoded string
	Err     error
}

// Func discards the input parameters and returns FakeEncodeAuth's Inv and
// Err field values.
func (f *FakeEncodeAuth) Func(_ time.Time, _ Auth) (string, error) {
	return f.Encoded, f.Err
}
