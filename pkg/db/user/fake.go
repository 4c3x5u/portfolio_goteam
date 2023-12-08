//go:build utest

package user

import "context"

// FakePutter is a test fake for Putter.
type FakePutter struct{ Err error }

// Put discards params and returns FakePutter.Err.
func (f *FakePutter) Put(_ context.Context, _ User) error { return f.Err }

// FakeGetter is a test fake for Getter.
type FakeGetter struct {
	User User
	Err  error
}

// Get discards params and returns FakeGetter.User and FakeGetter.Err.
func (f *FakeGetter) Get(ctx context.Context, username string) (User, error) {
	return f.User, f.Err
}
