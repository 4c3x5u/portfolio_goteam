package token

import "time"

// FakeGenerator is a test fake for Generator.
type FakeGenerator struct {
	InSub  string
	InExp  time.Time
	OutRes string
	OutErr error
}

// Generate implements the Generator interface on *FakeGenerator.
func (f *FakeGenerator) Generate(sub string, exp time.Time) (string, error) {
	f.InSub, f.InExp = sub, exp
	return f.OutRes, f.OutErr
}
