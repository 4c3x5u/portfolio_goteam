package user

// FakeInserter is a test fake for Inserter[User].
type FakeInserter struct{ OutErr error }

// Insert implements the Inserter[User] interface on FakeInserter.
func (f *FakeInserter) Insert(_ User) error { return f.OutErr }

// FakeSelector is a test fake for Selector[User].
type FakeSelector struct {
	OutRes User
	OutErr error
}

// Select implements the Selector[User] interface on FakeSelector.
func (f *FakeSelector) Select(_ string) (User, error) {
	return f.OutRes, f.OutErr
}
