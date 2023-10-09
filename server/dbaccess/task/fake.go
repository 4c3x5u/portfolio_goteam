package task

// FakeInserter is a test fake for dbaccess.Inserter[Task].
type FakeInserter struct{ Err error }

// Insert implements the dbaccess.Inserter[Task] interface on FakeInserter.
func (f *FakeInserter) Insert(_ Task) error { return f.Err }

// FakeSelector is a test fake for dbaccess.Selector[Record].
type FakeSelector struct{ Err error }

// Select implements the dbaccess.Selector[Record] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (Record, error) {
	return Record{}, f.Err
}
