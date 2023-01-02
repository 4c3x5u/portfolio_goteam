package db

import "time"

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

// FakeCreatorTwoStrTime is a test fake for CreatorTwoStrTime.
type FakeCreatorTwoStrTime struct {
	InArgA string
	InArgB string
	InArgC time.Time
	OutErr error
}

// Create implements the CreatorTwoStrTime interface on FakeCreatorTwoStrTime.
func (f *FakeCreatorTwoStrTime) Create(argA string, argB string, argC time.Time) error {
	f.InArgA, f.InArgB, f.InArgC = argA, argB, argC
	return f.OutErr
}
