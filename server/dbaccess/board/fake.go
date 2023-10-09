package board

// FakeSelector is a test fake for Selector[Record].
type FakeSelector struct{ Err error }

// Select implements the Selector[Record] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (Record, error) {
	return Record{}, f.Err
}

// FakeInserter is a test fake for Inserter[InRecord].
type FakeInserter struct{ Err error }

// Insert implements the Inserter[InRecord] interface on FakeInserter.
func (f *FakeInserter) Insert(_ InRecord) error { return f.Err }
