package task

// FakeInserter is a test fake for dbaccess.Inserter[InRecord].
type FakeInserter struct{ Err error }

// Insert implements the dbaccess.Inserter[InRecord] interface on FakeInserter.
func (f *FakeInserter) Insert(_ InRecord) error { return f.Err }

// FakeUpdater is a test fake for dbaccess.Updater[UpRecord].
type FakeUpdater struct{ Err error }

// Update implements the dbaccess.Updater[UpRecord] interface on FakeUpdater.
func (f *FakeUpdater) Update(string, UpRecord) error { return f.Err }
