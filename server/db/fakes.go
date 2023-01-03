package db

// FakeReaderUser is a test fake for Reader[[]byte].
type FakeReaderUser struct {
	InArg  string
	OutRes *User
	OutErr error
}

// Read implements the Reader[[]byte] interface on FakeReaderUser.
func (f *FakeReaderUser) Read(arg string) (*User, error) {
	f.InArg = arg
	return f.OutRes, f.OutErr
}

// FakeCreatorUser is a test fake for CreatorStrBytes.
type FakeCreatorUser struct {
	InArg  *User
	OutErr error
}

// Create implements the CreatorStrBytes interface on FakeCreatorUser.
func (f *FakeCreatorUser) Create(user *User) error {
	f.InArg = user
	return f.OutErr
}

// FakeCreatorSession is a test fake for CreatorTwoStrTime.
type FakeCreatorSession struct {
	InArg  *Session
	OutErr error
}

// Create implements the CreatorTwoStrTime interface on FakeCreatorSession.
func (f *FakeCreatorSession) Create(session *Session) error {
	f.InArg = session
	return f.OutErr
}
