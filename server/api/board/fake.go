package board

// fakePOSTReqValidator is a test fake for POSTReqValidator.
type fakePOSTReqValidator struct{ OutErrMsg string }

// Validate implements the POSTReqValidator interface on fakePOSTReqValidator.
func (f *fakePOSTReqValidator) Validate(_ POSTReqBody) string {
	return f.OutErrMsg
}

// fakeDELETEReqValidator is a test fake for DELETEReqValidator.
type fakeDELETEReqValidator struct{ OutOK bool }

// Validate implements the DELETEReqValidator interface on
// fakeDELETEReqValidator.
func (f *fakeDELETEReqValidator) Validate(_ string) bool { return f.OutOK }
