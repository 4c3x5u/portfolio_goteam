package dbaccess

// FakeCloser is a test fake for Closer.
type FakeCloser struct{}

// Close implements the Closer interface on FakeCloser.
func (c *FakeCloser) Close() {}

// FakeCounter is a test fake for Counter.
type FakeCounter struct {
	OutRes int
	OutErr error
}

// Count implements the Counter interface on FakeCounter.
func (f *FakeCounter) Count(_ string) (int, error) { return f.OutRes, f.OutErr }

// FakeDeleter is a test fake for Deleter.
type FakeDeleter struct{ OutErr error }

// Delete implements the Deleter interface on FakeDeleter.
func (f *FakeDeleter) Delete(_ string) error { return f.OutErr }

// FakeUpdater is a test fake for Updater.
type FakeUpdater struct{ OutErr error }

// Update implements the Updater interface on FakeUpdater.
func (f *FakeUpdater) Update(_, _ string) error { return f.OutErr }
