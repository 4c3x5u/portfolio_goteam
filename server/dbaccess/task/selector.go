package task

// Record represents a record in the task table.
type Record struct{}

// Selector can be used to select a record from the task table.
type Selector struct{}

// NewSelector creates and returns a new Selector.
func NewSelector() Selector { return Selector{} }

// Select selects a record from the task table with the given ID.
func (s Selector) Select(string) (Record, error) { return Record{}, nil }
