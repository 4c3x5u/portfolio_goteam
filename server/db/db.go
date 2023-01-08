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
