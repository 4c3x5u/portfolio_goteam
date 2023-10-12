package subtask

// FakeSelector is a test fake for dbaccess.Selector[Record].
type FakeSelector struct{ Err error }

// Select implements the dbaccess.Selector[Record] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (Record, error) {
	return Record{}, f.Err
}
