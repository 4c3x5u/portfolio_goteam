package register

// fakeUserValidator is a test fake for UserValidator.
type fakeUserValidator struct{ validationErrs ValidationErrors }

// Validate implements the UserValidator interface on fakeUserValidator.
func (f *fakeUserValidator) Validate(_ ReqBody) ValidationErrors {
	return f.validationErrs
}

// fakeStringValidator is a test fake for StringValidator.
type fakeStringValidator struct{ errs []string }

// Validate implements the StringValidator interface on fakeStringValidator.
func (f *fakeStringValidator) Validate(_ string) []string { return f.errs }

// fakeHasher is a test fake for Hasher.
type fakeHasher struct {
	hash []byte
	err  error
}

// Hash implements the Hasher interface on fakeHasher.
func (f *fakeHasher) Hash(_ string) ([]byte, error) {
	return f.hash, f.err
}
