// Package db contains code for working with the goteam database.
package db

// Creator describes a type that creates a record in the database.
type Creator[T any] interface{ Create(record T) error }

// Reader describes a type reads a record from the database.
type Reader[T any] interface{ Read(id string) (T, error) }
