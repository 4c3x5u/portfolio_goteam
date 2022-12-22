package db

// Existor represents a type that checks whether a certain entity exists or not
// based on a string term. It should return true if the said entity exists and
// false if it doesn't, alongside any errors that may occur during the process
// of checking.
type Existor interface {
	Exists(term string) (bool, error)
}
