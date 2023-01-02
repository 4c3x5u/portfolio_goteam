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

// FakeCreatorStrBytes is a test fake for CreatorStrBytes.
type FakeCreatorStrBytes struct {
	InArgA string
	InArgB []byte
	OutErr error
}

// Create implements the CreatorStrBytes interface on FakeCreatorStrBytes.
func (f *FakeCreatorStrBytes) Create(argA string, argB []byte) error {
	f.InArgA, f.InArgB = argA, argB
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
