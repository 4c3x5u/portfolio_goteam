package db

import "time"

// FakeReaderBytes is a test fake for Reader[[]byte].
type FakeReaderBytes struct {
	InArg  string
	OutRes []byte
	OutErr error
}

// Read implements the Reader[[]byte] interface on FakeReaderBytes.
func (f *FakeReaderBytes) Read(arg string) ([]byte, error) {
	f.InArg = arg
	return f.OutRes, f.OutErr
}

// FakeExistor is a test fake for Existor.
type FakeExistor struct {
	InArg     string
	OutExists bool
	OutErr    error
}

// Exists implements the Existor interface on FakeExistor.
func (f *FakeExistor) Exists(arg string) (bool, error) {
	f.InArg = arg
	return f.OutExists, f.OutErr
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
