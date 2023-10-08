package userboard

// FakeSelector is a test fake for dbaccess.RelSelector[bool].
type FakeSelector struct {
	IsAdmin bool
	Err     error
}

// Select implements the RelSelector[bool] interface on FakeSelector.
func (f *FakeSelector) Select(_, _ string) (bool, error) {
	return f.IsAdmin, f.Err
}
