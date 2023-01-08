package register

// fakeValidator is a test fake for Validator.
type fakeValidator struct {
	inReqBody ReqBody
	outErrs   ValidationErrs
}

// Validate implements the Validator interface on fakeValidator.
// struct.
func (f *fakeValidator) Validate(reqBody ReqBody) ValidationErrs {
	f.inReqBody = reqBody
	return f.outErrs
}

// fakeStringValidator is a test fake for StringValidator.
type fakeStringValidator struct {
	inArg   string
	outErrs []string
}

// Validate implements the StringValidator interface on fakeStringValidator.
func (f *fakeStringValidator) Validate(arg string) (errs []string) {
	f.inArg = arg
	return f.outErrs
}

// fakeHasher is a test fake for Hasher.
type fakeHasher struct {
	inPlaintext string
	outHash     []byte
	outErr      error
}

// Hash implements the Hasher interface on fakeHasher.
func (f *fakeHasher) Hash(plaintext string) ([]byte, error) {
	f.inPlaintext = plaintext
	return f.outHash, f.outErr
}
