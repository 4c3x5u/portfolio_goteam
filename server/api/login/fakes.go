package login

// fakeHashComparer is a test fake for HashComparer.
type fakeHashComparer struct {
	inArgA []byte
	inArgB string
	outRes bool
	outErr error
}

// Compare implements the HashComparer interface on fakeHashComparer.
func (f *fakeHashComparer) Compare(hash []byte, plaintext string) (bool, error) {
	f.inArgA, f.inArgB = hash, plaintext
	return f.outRes, f.outErr
}
