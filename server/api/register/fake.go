package register

// fakeValidator is a test fake for Validator.
type fakeValidator struct{ outErrs ValidationErrs }

// Validate implements the Validator interface on fakeValidator.
func (f *fakeValidator) Validate(_ ReqBody) ValidationErrs { return f.outErrs }

// fakeStringValidator is a test fake for StringValidator.
type fakeStringValidator struct{ outErrs []string }

// Validate implements the StringValidator interface on fakeStringValidator.
func (f *fakeStringValidator) Validate(_ string) []string { return f.outErrs }

// fakeHasher is a test fake for Hasher.
type fakeHasher struct {
	outHash []byte
	outErr  error
}

// Hash implements the Hasher interface on fakeHasher.
func (f *fakeHasher) Hash(_ string) ([]byte, error) {
	return f.outHash, f.outErr
}
