package board

import "net/url"

// fakePOSTReqValidator is a test fake for POSTReqValidator.
type fakePOSTReqValidator struct{ OutErrMsg string }

// Validate implements the POSTReqValidator interface on fakePOSTReqValidator.
func (f *fakePOSTReqValidator) Validate(_ POSTReqBody) string {
	return f.OutErrMsg
}

// fakeDELETEReqValidator is a test fake for DELETEReqValidator.
type fakeDELETEReqValidator struct {
	OutID string
	OutOK bool
}

// Validate implements the DELETEReqValidator interface on
// fakeDELETEReqValidator.
func (f *fakeDELETEReqValidator) Validate(_ url.Values) (string, bool) {
	return f.OutID, f.OutOK
}
