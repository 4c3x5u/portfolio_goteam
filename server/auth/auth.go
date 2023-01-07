package auth

import "time"

// Generator represents a type that can be used to generate a special string for
// a given subject.
type Generator interface {
	Generate(sub string, exp time.Time) (string, error)
}
