package board

import "database/sql"

// Board represents a record in the board table.
type Board struct {
	id   int
	name string
}

// Selector can be used to read records from the board table.
type Selector struct{ db *sql.DB }

// NewSelector creates and returns a new Selector.
func NewSelector(db *sql.DB) Selector { return Selector{db: db} }

// Select selects a record from the board table with the given id.
func (s Selector) Select(id string) (Board, error) {
	var board Board
	err := s.db.
		QueryRow(`SELECT id, name FROM app.board WHERE id = $1`, id).
		Scan(&board.id, &board.name)
	return board, err
}
