package board

import "net/url"

// fakePOSTReqValidator is a test fake for POSTReqValidator.
type fakePOSTReqValidator struct {
	InReqBody POSTReqBody
	OutErr    error
}

// Validate implements the POSTReqValidator interface on fakePOSTReqValidator.
// It assigns the parameters passed into it to their corresponding In... fields
// on the fake instance and returns its Out.. fields as per function signature.
func (f *fakePOSTReqValidator) Validate(reqBody POSTReqBody) error {
	f.InReqBody = reqBody
	return f.OutErr
}

// fakeDELETEReqValidator is a test fake for DELETEReqValidator.
type fakeDELETEReqValidator struct {
	InQParams url.Values
	OutID     string
	OutOK     bool
}

// Validate implements the DELETEReqValidator interface on
// fakeDELETEReqValidator. It assigns the parameters passed into it to
// their corresponding In... fields on the fake instance and returns its Out..
// fields as per function signature.
func (f *fakeDELETEReqValidator) Validate(qParams url.Values) (string, bool) {
	f.InQParams = qParams
	return f.OutID, f.OutOK
}
