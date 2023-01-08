package db

// FakeUserCreator is a test fake for Creator[User].
type FakeUserCreator struct {
	InArg  User
	OutErr error
}

// Create implements the Creator[User] interface on FakeUserCreator.
func (f *FakeUserCreator) Create(user User) error {
	f.InArg = user
	return f.OutErr
}

// FakeUserReader is a test fake for Reader[User].
type FakeUserReader struct {
	InArg  string
	OutRes User
	OutErr error
}

// Read implements the Reader[User] interface on FakeUserReader.
func (f *FakeUserReader) Read(arg string) (User, error) {
	f.InArg = arg
	return f.OutRes, f.OutErr
}
