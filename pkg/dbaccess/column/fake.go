package column

// FakeUpdater is a test fake for Updater.
type FakeUpdater struct{ Err error }

// Update implements the dbaccess.Updater[string] interface on FakeUpdater.
func (f *FakeUpdater) Update(_ string, _ []Task) error { return f.Err }
