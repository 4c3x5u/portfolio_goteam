package board

// fakeStrValidator is a test fake for StrValidator.
type fakeStrValidator struct{ OutErr error }

// Validate implements the StrValidator interface on fakeStrValidator.
func (f *fakeStrValidator) Validate(_ string) error { return f.OutErr }

// fakeIDValidator is a test fake for IDValidator.
type fakeIDValidator struct{ OutErr error }

// Validate implements the DELETEReqValidator interface on
// fakeIDValidator.
func (f *fakeIDValidator) Validate(_ string) error { return f.OutErr }
