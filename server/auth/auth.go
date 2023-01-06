package auth

// Generator represents a type that can be used to generate a special string for
// a given subject.
type Generator interface {
	Generate(idSubject string) (string, error)
}
