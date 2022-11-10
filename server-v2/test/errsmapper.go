package test

// ErrsMapper allows errors returned from API routes to be investigated in
// Route.Run through the field name that is set on each RoutCase.
type ErrsMapper interface {
	MapErrs() map[string][]string
}
