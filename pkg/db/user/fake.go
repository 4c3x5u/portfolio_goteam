//go:build utest

package user

import "context"

// FakePutter is a test fake for Putter.
type FakePutter struct{ Err error }

// Put discards params and returns FakePutter.Err.
func (f FakePutter) Put(_ context.Context, _ User) error { return f.Err }
