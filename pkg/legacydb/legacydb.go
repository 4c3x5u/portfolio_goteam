// Package dbaccess contains code for working with the goteam database.
package legacydb

// Inserter describes a type that inserts a record into a table in the database.
type Inserter[T any] interface{ Insert(T) error }

// Selector describes a type that selects a record from a table in the database.
type Selector[T any] interface{ Select(string) (T, error) }

// Deleter describes a type that deletes a record from a table in the database.
type Deleter interface{ Delete(string) error }

// Counter describes a type counts records in the database.
type Counter interface{ Count(string) (int, error) }

// Updater defines a type that updates a field of a record in the database with
// a new value.
type Updater[T any] interface{ Update(string, T) error }
