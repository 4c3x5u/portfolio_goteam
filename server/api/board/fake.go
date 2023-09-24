package board

// fakeStringValidator is a test fake for StringValidator.
type fakeStringValidator struct{ OutErr error }

// Validate implements the StringValidator interface on fakeStringValidator.
func (f *fakeStringValidator) Validate(_ string) error { return f.OutErr }
