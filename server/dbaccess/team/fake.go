package team

// FakeSelector is a test fake for Selector[Record].
type FakeSelector struct{ Err error }

// Select implements the Selector[Record] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (Record, error) {
	return Record{}, f.Err
}
