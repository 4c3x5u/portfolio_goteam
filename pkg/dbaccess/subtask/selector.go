package subtask

import "database/sql"

// Record defines the subtask data returned from Selector.Select.
type Record struct{ TaskID int }

// Selector can be used to select a record from the task table.
type Selector struct{ db *sql.DB }

// NewSelector creates and returns a new Selector.
func NewSelector(db *sql.DB) Selector { return Selector{db: db} }

// Select selects a record from the task table with the given ID.
func (s Selector) Select(id string) (Record, error) {
	var rec Record
	if err := s.db.QueryRow(
		`SELECT taskID FROM app.subtask WHERE id = $1`, id,
	).Scan(&rec.TaskID); err != nil {
		return Record{}, err
	}
	return rec, nil
}
