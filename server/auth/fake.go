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
	OutErr  error
}

// Validate implements the Validator interface on FakeTokenValidator.
func (f FakeTokenValidator) Validate(token string) (string, error) {
	f.InToken = token
	return f.OutSub, f.OutErr
}
