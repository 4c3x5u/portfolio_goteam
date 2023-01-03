// Package db contains code for working with the goteam database.
package db

// Reader represents a type reads a record from the database.
type Reader[T any] interface {
	Read(id string) (T, error)
}

// Creator represents a type that inserts a record to the database.
type Creator[T any] interface {
	Create(record T) error
}

// Deleter represents a type that deletes a record from the database.
type Deleter interface {
	Delete(id string) error
}
