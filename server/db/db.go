// Package db contains code for working with the goteam database.
package db

// Inserter describes a type that inserts a record into the database.
type Inserter[T any] interface{ Insert(record T) error }

// Selector describes a type that selects a record from the database.
type Selector[T any] interface{ Select(id string) (T, error) }

// RelSelector describes a type that selects a record from a many-to-many
// relationship table from the database. It uses the two IDs - one for each
// entity that is subject to the relationship.
type RelSelector[T any] interface {
	Select(idA, idB string) (T, error)
}

// Counter describes a type counts records in the database.
type Counter interface{ Count(id string) int }
