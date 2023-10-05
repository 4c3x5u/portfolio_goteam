package userboard

// FakeSelector is a test fake for dbaccess.RelSelector[bool].
type FakeSelector struct {
	OutIsAdmin bool
	OutErr     error
}

// Select implements the RelSelector[bool] interface on FakeSelector.
func (f *FakeSelector) Select(_, _ string) (bool, error) {
	return f.OutIsAdmin, f.OutErr
}
