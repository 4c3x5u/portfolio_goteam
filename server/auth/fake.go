package auth

import "time"

// FakeTokenGenerator is a test fake for TokenGenerator.
type FakeTokenGenerator struct {
	InSub  string
	InExp  time.Time
	OutRes string
	OutErr error
}

// Generate implements the CookieGenerator interface on FakeTokenGenerator.
// It assigns the parameters passed into it to their corresponding In... fields
// on the fake instance and returns its Out.. fields as per function signature.
func (f *FakeTokenGenerator) Generate(
	sub string, exp time.Time,
) (string, error) {
	f.InSub, f.InExp = sub, exp
	return f.OutRes, f.OutErr
}

// FakeTokenValidator is a test fake for Validator.
type FakeTokenValidator struct {
	InToken string
	OutSub  string
}

// Validate implements the Validator interface on FakeTokenValidator. It assigns
// the parameters passed into it to their corresponding In... fields on the fake
// instance and returns its Out.. fields as per function signature.
func (f *FakeTokenValidator) Validate(token string) string {
	f.InToken = token
	return f.OutSub
}

// FakeHeaderReader is a test fake for HeaderReader.
type FakeHeaderReader struct {
	InHeaderValue string
	OutToken      string
}

// Read implements the HeaderReader interface on FakeHeaderReader. It assigns
// the parameters passed into it to their corresponding In... fields on the fake
// instance and returns its Out.. fields as per function signature.
func (f *FakeHeaderReader) Read(headerValue string) string {
	f.InHeaderValue = headerValue
	return f.OutToken
}
