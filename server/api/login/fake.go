package login

// fakeReqValidator is a test fake for ReqValidator.
type fakeReqValidator struct{ outOK bool }

// Validate implements the ReqValidator interface on fakeReqValidator.
func (f *fakeReqValidator) Validate(_ ReqBody) bool { return f.outOK }

// fakeHashComparer is a test fake for Comparer.
type fakeHashComparer struct{ outErr error }

// Compare implements the Comparer interface on fakeHashComparer.
func (f *fakeHashComparer) Compare(_ []byte, _ string) error { return f.outErr }
