package task

// FakeInserter is a test fake for dbaccess.Inserter[Task].
type FakeInserter struct{ Err error }

// Insert implements the dbaccess.Inserter[Task] interface on FakeInserter.
func (f *FakeInserter) Insert(_ Task) error { return f.Err }
