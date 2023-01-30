package login

// fakeValidator is a test fake for Validator.
type fakeValidator struct {
	inReqBody ReqBody
	outOK     bool
}

// Validate implements the Validator interface on fakeValidator.
func (f *fakeValidator) Validate(reqBody ReqBody) bool {
	f.inReqBody = reqBody
	return f.outOK
}

// fakeHashComparer is a test fake for Comparer.
type fakeHashComparer struct {
	inHash      []byte
	inPlaintext string
	outErr      error
}

// Compare implements the Comparer interface on fakeHashComparer.
func (f *fakeHashComparer) Compare(hash []byte, plaintext string) error {
	f.inHash, f.inPlaintext = hash, plaintext
	return f.outErr
}
