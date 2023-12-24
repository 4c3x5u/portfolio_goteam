package loginapi

// fakeReqValidator is a test fake for ReqValidator.
type fakeReqValidator struct{ isValid bool }

// Validate implements the ReqValidator interface on fakeReqValidator.
func (f *fakeReqValidator) Validate(_ PostReq) bool { return f.isValid }

// fakeHashComparer is a test fake for Comparator.
type fakeHashComparer struct{ err error }

// Compare implements the Comparator interface on fakeHashComparer.
func (f *fakeHashComparer) Compare(_ []byte, _ string) error { return f.err }
