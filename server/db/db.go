package db

import "time"

// Existor represents a type that checks whether a certain entity exists or not
// based on a string argument.
type Existor interface {
	Exists(id string) (bool, error)
}

// CreatorStrBytes represents a type that creates an entity in the database
// based on a string and a []byte argument.
type CreatorStrBytes interface {
	Create(arg1 string, arg2 []byte) error
}

// CreatorTwoStrTime represents a type that creates an entity in the database
// based on three arguments – two strings and one time.Time.
type CreatorTwoStrTime interface {
	Create(arg1, arg2 string, arg3 time.Time) error
}

// Deleter represents a type that deletes a given entity from the database using
// a string argument.
type Deleter interface {
	Delete(id string) error
}
