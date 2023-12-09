//go:build utest

package task

import "context"

// FakePutter is a test fake for Putter.
type FakePutter struct{ Err error }

// Put discards params and returns FakePutter.Err.
func (f *FakePutter) Put(_ context.Context, _ Task) error { return f.Err }
