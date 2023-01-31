package register

// fakeValidator is a test fake for Validator.
type fakeValidator struct {
	inReqBody ReqBody
	outErrs   ValidationErrs
}

// Validate implements the Validator interface on fakeValidator.
// struct. It assigns the parameters passed into it to their corresponding In...
// fields on the fake instance and returns its Out.. fields as per function
// signature.
func (f *fakeValidator) Validate(reqBody ReqBody) ValidationErrs {
	f.inReqBody = reqBody
	return f.outErrs
}

// fakeStringValidator is a test fake for StringValidator.
type fakeStringValidator struct {
	inVal   string
	outErrs []string
}

// Validate implements the StringValidator interface on fakeStringValidator.
// It assigns the parameters passed into it to their corresponding In... fields
// on the fake instance and returns its Out.. fields as per function signature.
func (f *fakeStringValidator) Validate(val string) (errs []string) {
	f.inVal = val
	return f.outErrs
}

// fakeHasher is a test fake for Hasher.
type fakeHasher struct {
	inPlaintext string
	outHash     []byte
	outErr      error
}

// Hash implements the Hasher interface on fakeHasher. It assigns the parameters
// passed into it to their corresponding In... fields on the fake instance and
// returns its Out.. fields as per function signature.
func (f *fakeHasher) Hash(plaintext string) ([]byte, error) {
	f.inPlaintext = plaintext
	return f.outHash, f.outErr
}
