package column

import "database/sql"

// Record represents a record in the column table.
type Record struct {
	ID      int
	BoardID int
	Order   int16
}

// Selector can be used to read records from the column table.
type Selector struct{ db *sql.DB }

// NewSelector creates and returns a new Selector.
func NewSelector(db *sql.DB) Selector { return Selector{db: db} }

// Select selects a record from the board table with the given id.
func (s Selector) Select(id string) (Record, error) {
	var rec Record
	err := s.db.QueryRow(
		`SELECT id, boardID, order FROM app."column" WHERE id = $1`, id,
	).Scan(&rec.ID, &rec.BoardID, &rec.Order)
	return rec, err
}
