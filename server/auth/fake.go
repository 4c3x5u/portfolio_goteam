package auth

import (
	"net/http"
	"time"
)

// FakeCookieGenerator is a test fake for CookieGenerator.
type FakeCookieGenerator struct {
	InSub  string
	InExp  time.Time
	OutRes *http.Cookie
	OutErr error
}

// Generate implements the CookieGenerator interface on FakeCookieGenerator.
func (f *FakeCookieGenerator) Generate(sub string, exp time.Time) (*http.Cookie, error) {
	f.InSub, f.InExp = sub, exp
	return f.OutRes, f.OutErr
}

// FakeTokenValidator is a test fake for TokenValidator.
type FakeTokenValidator struct {
	InToken string
	OutSub  string
	OutErr  error
}

// Validate implements the TokenValidator interface on FakeTokenValidator.
func (f FakeTokenValidator) Validate(token string) (string, error) {
	f.InToken = token
	return f.OutSub, f.OutErr
}
