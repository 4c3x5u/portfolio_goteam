package login

// fakeComparer is a test fake for Comparer.
type fakeComparer struct {
	inArgA []byte
	inArgB string
	outRes bool
	outErr error
}

// Compare implements the Comparer interface on fakeComparer.
func (f *fakeComparer) Compare(hash []byte, plaintext string) (bool, error) {
	f.inArgA, f.inArgB = hash, plaintext
	return f.outRes, f.outErr
}
