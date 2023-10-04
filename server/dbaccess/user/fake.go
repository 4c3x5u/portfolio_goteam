package user

// FakeInserter is a test fake for Inserter[Record].
type FakeInserter struct{ OutErr error }

// Insert implements the Inserter[Record] interface on FakeInserter.
func (f *FakeInserter) Insert(_ Record) error { return f.OutErr }

// FakeSelector is a test fake for Selector[Record].
type FakeSelector struct {
	OutRes Record
	OutErr error
}

// Select implements the Selector[Record] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (Record, error) {
	return f.OutRes, f.OutErr
}
