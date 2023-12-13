package legacydb

// FakeCounter is a test fake for Counter.
type FakeCounter struct {
	BoardCount int
	Err        error
}

// Count implements the Counter interface on FakeCounter.
func (f *FakeCounter) Count(_ string) (int, error) { return f.BoardCount, f.Err }

// FakeDeleter is a test fake for Deleter.
type FakeDeleter struct{ Err error }

// Delete implements the Deleter interface on FakeDeleter.
func (f *FakeDeleter) Delete(_ string) error { return f.Err }

// FakeUpdater is a test fake for Updater.
type FakeUpdater struct{ Err error }

// Update implements the Updater interface on FakeUpdater.
func (f *FakeUpdater) Update(_, _ string) error { return f.Err }
