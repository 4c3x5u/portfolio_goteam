package board

// FakeSelector is a test fake for Selector[Board].
type FakeSelector struct{ OutErr error }

// Select implements the Selector[Board] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (Board, error) {
	return Board{}, f.OutErr
}

// FakeInserter is a test fake for Inserter[InBoard].
type FakeInserter struct{ OutErr error }

// Insert implements the Inserter[InBoard] interface on FakeInserter.
func (f *FakeInserter) Insert(_ InBoard) error { return f.OutErr }
