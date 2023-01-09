package cookie

import (
	"net/http"
	"time"
)

// FakeAuthGenerator is a test fake for AuthGenerator.
type FakeAuthGenerator struct {
	InSub  string
	InExp  time.Time
	OutRes *http.Cookie
	OutErr error
}

// Generate implements the AuthGenerator interface on FakeAuthGenerator.
func (f *FakeAuthGenerator) Generate(sub string, exp time.Time) (*http.Cookie, error) {
	f.InSub, f.InExp = sub, exp
	return f.OutRes, f.OutErr
}
