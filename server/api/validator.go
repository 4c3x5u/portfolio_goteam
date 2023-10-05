package api

// StringValidator describes a type that can be used to validate a string.
type StringValidator interface{ Validate(string) error }
