package column

type FakeSelector struct{ OutErr error }

// Select implements the Selector[Record] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (Record, error) {
	return Record{}, f.OutErr
}
