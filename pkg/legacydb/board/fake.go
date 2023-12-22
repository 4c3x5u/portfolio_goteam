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

// FakeSelectorByTeamID is a test fake for dbaccess.Selector[[]Record].
type FakeSelectorByTeamID struct {
	Recs []Record
	Err  error
}

// Select implements the dbaccess.Selector[[]Record] interface on FakeSelector.
func (f *FakeSelectorByTeamID) Select(_ string) ([]Record, error) {
	return f.Recs, f.Err
}
