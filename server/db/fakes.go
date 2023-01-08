package db

// FakeCreatorUser is a test fake for Creator[*User].
type FakeCreatorUser struct {
	InArg  *User
	OutErr error
}

// Create implements the Creator[*User] interface on FakeCreatorUser.
func (f *FakeCreatorUser) Create(user *User) error {
	f.InArg = user
	return f.OutErr
}

// FakeReaderUser is a test fake for Reader[*User].
type FakeReaderUser struct {
	InArg  string
	OutRes *User
	OutErr error
}

// Read implements the Reader[*User] interface on FakeReaderUser.
func (f *FakeReaderUser) Read(arg string) (*User, error) {
	f.InArg = arg
	return f.OutRes, f.OutErr
}
