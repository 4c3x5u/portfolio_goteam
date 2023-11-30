package task

// FakeInserter is a test fake for dbaccess.Inserter[InRecord].
type FakeInserter struct{ Err error }

// Insert implements the dbaccess.Inserter[InRecord] interface on FakeInserter.
func (f *FakeInserter) Insert(_ InRecord) error { return f.Err }

// FakeSelector is a test fake for dbaccess.Selector[Record].
type FakeSelector struct{ Err error }

// Select implements the dbaccess.Selector[Record] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (Record, error) {
	return Record{}, f.Err
}

// FakeUpdater is a test fake for dbaccess.Updater[UpRecord].
type FakeUpdater struct{ Err error }

// Update implements the dbaccess.Updater[UpRecord] interface on FakeUpdater.
func (f *FakeUpdater) Update(string, UpRecord) error { return f.Err }
