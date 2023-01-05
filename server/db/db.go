// Package db contains code for working with the goteam database.
package db

// Creator represents a type that inserts a record to the database.
type Creator[T any] interface {
	Create(record T) error
}

// Reader represents a type reads a record from the database.
type Reader[T any] interface {
	Read(id string) (T, error)
}

// Upserter represents a type that inserts a record to do the database, or if
// the record already exists updates it.
type Upserter[T any] interface {
	Upsert(record T) error
}

// Deleter represents a type that deletes a record from the database.
type Deleter interface {
	Delete(id string) error
}
