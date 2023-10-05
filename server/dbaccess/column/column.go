// Package column contains code for working with the column table.
package column

type Record struct {
	ID      int
	BoardID int
	Order   int16
}

// Task contains data needed to move a task from one column to another.
type Task struct {
	ID    int
	Order int
}
