package team

// FakeSelector is a test fake for Selector[Record].
type FakeSelector struct {
	Rec Record
	Err error
}

// Select implements the Selector[Record] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (Record, error) { return f.Rec, f.Err }
