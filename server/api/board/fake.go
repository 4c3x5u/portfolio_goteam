package board

// fakeStringValidator is a test fake for StringValidator.
type fakeStringValidator struct{ OutErr error }

// Validate implements the StringValidator interface on fakeStringValidator.
func (f *fakeStringValidator) Validate(_ string) error { return f.OutErr }

// fakeIDValidator is a test fake for IDValidator.
type fakeIDValidator struct{ OutErr error }

// Validate implements the DELETEReqValidator interface on
// fakeIDValidator.
func (f *fakeIDValidator) Validate(_ string) error { return f.OutErr }
