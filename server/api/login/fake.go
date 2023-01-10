package login

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
