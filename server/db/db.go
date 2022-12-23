package db

// Existor represents a type that checks whether a certain entity exists or not
// based on a string term. It should return true if the said entity exists and
// false if it doesn't, alongside any errors that may occur during the process
// of checking.
type Existor interface {
	Exists(term string) (bool, error)
}

// Creator represents a type that creates an entity in the database based on
// the given args slice. A Creator should have either a query template that
// args are used in or handle args in some other fashion. The error return
// should be used for any errors encountered during creation.
type Creator interface {
	Create(args ...any) error
}
