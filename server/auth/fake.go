package auth

import "time"

// FakeTokenGenerator is a test fake for TokenGenerator.
type FakeTokenGenerator struct {
	AuthToken string
	Err       error
}

// Generate implements the CookieGenerator interface on FakeTokenGenerator.
func (f *FakeTokenGenerator) Generate(_ string, _ time.Time) (string, error) {
	return f.AuthToken, f.Err
}

// FakeTokenValidator is a test fake for Validator.
type FakeTokenValidator struct{ Sub string }

// Validate implements the Validator interface on FakeTokenValidator.
func (f *FakeTokenValidator) Validate(_ string) string { return f.Sub }

// FakeHeaderReader is a test fake for HeaderReader.
type FakeHeaderReader struct{ AuthToken string }

// Read implements the HeaderReader interface on FakeHeaderReader.
func (f *FakeHeaderReader) Read(_ string) string { return f.AuthToken }
