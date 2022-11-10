// Package test contains code that is common between certain tests and
// encapsulates each type of test, test suite, and test case into types that
// contains the pieces of data needed to run the said common code.
package test

// ErrsMapper allows errors returned from API routes to be investigated in
// Route.Run via the field name that is set on each RoutCase.
type ErrsMapper interface {
	ErrsMap() map[string][]string
}
