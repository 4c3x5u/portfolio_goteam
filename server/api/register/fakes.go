package register

// fakeValidator is a test fake for Validator
type fakeValidator struct {
	inReqBody *ReqBody
	outErrs   *Errs
}

// Validate implements the Validator interface on fakeValidator.
// struct.
func (f *fakeValidator) Validate(reqBody *ReqBody) *Errs {
	f.inReqBody = reqBody
	return f.outErrs
}

// fakeValidatorStr is a test fake for ValidatorStr.
type fakeValidatorStr struct {
	inArg   string
	outErrs []string
}

// Validate implements the ValidatorStr interface on fakeValidatorStr.
func (f *fakeValidatorStr) Validate(arg string) (errs []string) {
	f.inArg = arg
	return f.outErrs
}

// fakeHasherPwd is a test fake for Hasher.
type fakeHasherPwd struct {
	inPlaintext string
	outHash     []byte
	outErr      error
}

// Hash implements the Hasher interface on fakeHasherPwd.
func (f *fakeHasherPwd) Hash(plaintext string) ([]byte, error) {
	f.inPlaintext = plaintext
	return f.outHash, f.outErr
}
