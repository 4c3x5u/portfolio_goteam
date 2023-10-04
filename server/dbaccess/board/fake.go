package board

// FakeSelector is a test fake for Selector[Record].
type FakeSelector struct{ OutErr error }

// Select implements the Selector[Record] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (Record, error) {
	return Record{}, f.OutErr
}

// FakeInserter is a test fake for Inserter[Board].
type FakeInserter struct{ OutErr error }

// Insert implements the Inserter[Board] interface on FakeInserter.
func (f *FakeInserter) Insert(_ Board) error { return f.OutErr }
