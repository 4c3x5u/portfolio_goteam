package db

import "time"

// Reader represents a type reads a record from the database.
type Reader[T any] interface {
	Read(id string) (T, error)
}

// CreatorStrBytes represents a type that creates a record in the database
// based on a string and a []byte argument.
type CreatorStrBytes interface {
	Create(arg1 string, arg2 []byte) error
}

// CreatorTwoStrTime represents a type that creates a record in the database
// based on three arguments – two strings and one time.Time.
type CreatorTwoStrTime interface {
	Create(arg1, arg2 string, arg3 time.Time) error
}

// Deleter represents a type that deletes a given entity from the database using
// a string argument.
type Deleter interface {
	Delete(id string) error
}
