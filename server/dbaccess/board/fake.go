package board

// FakeInserter is a test fake for Inserter[InRecord].
type FakeInserter struct{ Err error }

// Insert implements the Inserter[InRecord] interface on FakeInserter.
func (f *FakeInserter) Insert(_ InRecord) error { return f.Err }

// FakeSelector is a test fake for Selector[Record].
type FakeSelector struct {
	Err   error
	Board Record
}

// Select implements the Selector[Record] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (Record, error) {
	return f.Board, f.Err
}

// FakeRecursiveSelector is a test fake for Selector[RecursiveRecord].
type FakeRecursiveSelector struct {
	Err    error
	Record RecursiveRecord
}

// Select implements the Selector[RecursiveRecord] interface on
// FakeRecursiveSelector.
func (f *FakeRecursiveSelector) Select(_ string) (RecursiveRecord, error) {
	return f.Record, f.Err
}
