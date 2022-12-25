package db

// Existor represents a type that checks whether a certain entity exists or not
// based on a string id. It should return true if the said entity exists and
// false if it doesn't, alongside any errors that may occur during the process
// of checking.
type Existor interface {
	Exists(id string) (bool, error)
}

// Creator represents a type that creates an entity in the database based on
// the given args slice. A Creator should have either a query template that
// args are used in or handle args in some other fashion. The error return
// should be used for any errors encountered during creation.
type Creator interface {
	Create(args ...any) error
}

// Deleter represents a type that deletes a given entity from the database using
// a string ID. It should return an error if the ID doesn't match an entity or
// if other errors occur during database query.
type Deleter interface {
	Delete(id string) error
}
