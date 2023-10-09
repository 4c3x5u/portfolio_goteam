package task

import "database/sql"

// Record represents a record in the task table.
type Record struct {
	ID          int
	ColumnID    int
	Title       string
	Description string
	Order       int
}

// Selector can be used to select a record from the task table.
type Selector struct{ db *sql.DB }

// NewSelector creates and returns a new Selector.
func NewSelector(db *sql.DB) Selector { return Selector{db: db} }

// Select selects a record from the task table with the given ID.
func (s Selector) Select(id string) (Record, error) {
	var rec Record
	if err := s.db.QueryRow(
		`SELECT id, columnID, title, description, "order" FROM app.task `+
			`WHERE id = $1`,
		id,
	).Scan(
		&rec.ID, &rec.ColumnID, &rec.Title, &rec.Description, &rec.Order,
	); err != nil {
		return Record{}, err
	}
	return rec, nil
}
