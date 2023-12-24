package log

// FakeErrorer is a test fake for Errorer.
type FakeErrorer struct{ Args []any }

// Log implements the Errorer interface on FakeErrorer. It assigns the message
// passed into it to the InMessage field on the fake instance.
func (f *FakeErrorer) Error(args ...any) { f.Args = args }
