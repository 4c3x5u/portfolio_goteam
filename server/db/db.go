// Package db contains code for working with the goteam database.
package db

// Inserter describes a type that inserts a record into the database.
type Inserter[T any] interface{ Insert(record T) error }

// Selector describes a type selects a record from the database.
type Selector[T any] interface{ Select(id string) (T, error) }

// Counter describes a type counts records in the database.
type Counter interface{ Count(id string) int }
