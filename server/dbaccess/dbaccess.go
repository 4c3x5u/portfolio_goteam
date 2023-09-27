// Package dbaccess contains code for accessing the GoTeam! database.
package dbaccess

// Inserter describes a type that inserts a record into a table in the database.
type Inserter[T any] interface{ Insert(record T) error }

// Selector describes a type that selects a record from a table in the database.
type Selector[T any] interface{ Select(id string) (T, error) }

// Deleter describes a type that deletes a record from a table in the database.
type Deleter interface{ Delete(id string) error }

// Counter describes a type counts records in the database.
type Counter interface{ Count(id string) (int, error) }

// Updater defines a type that updates a field of a record in the database with
// a new value.
type Updater interface {
	Update(id string, newValue string) error
}

// RelSelector describes a type that selects a record from a many-to-many
// relationship table from the database. It uses the two IDs - one for each
// entity that is subject to the relationship.
type RelSelector[T any] interface {
	Select(idA, idB string) (T, error)
}
