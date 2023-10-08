package column

// FakeSelector is a test fake for dbaccess.Selector[Record].
type FakeSelector struct{ Err error }

// Select implements the dbaccess.Selector[Record] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (Record, error) {
	return Record{}, f.Err
}

// FakeUpdater is a test fake for Updater.
type FakeUpdater struct{ Err error }

// Update implements the dbaccess.Updater[string] interface on FakeUpdater.
func (f *FakeUpdater) Update(_ string, _ []Task) error { return f.Err }
