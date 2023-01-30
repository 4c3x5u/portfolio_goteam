package board

// fakeValidator is a test fake for POSTReqValidator.
type fakePOSTReqValidator struct {
	InReqBody POSTReqBody
	OutErr    error
}

// Validate implements the POSTReqValidator interface on fakePOSTReqValidator.
func (f *fakePOSTReqValidator) Validate(reqBody POSTReqBody) error {
	f.InReqBody = reqBody
	return f.OutErr
}
