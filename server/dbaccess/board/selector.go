package board

import "database/sql"

// Record represents a record in the board table.
type Record struct {
	ID     int
	Name   string
	TeamID int
}

// Selector can be used to read records from the board table.
type Selector struct{ db *sql.DB }

// NewSelector creates and returns a new Selector.
func NewSelector(db *sql.DB) Selector { return Selector{db: db} }

// Select selects a record from the board table with the given id.
func (s Selector) Select(id string) (Record, error) {
	var board Record
	err := s.db.
		QueryRow(`SELECT id, name, teamID FROM app.board WHERE id = $1`, id).
		Scan(&board.ID, &board.Name, &board.TeamID)
	return board, err
}
