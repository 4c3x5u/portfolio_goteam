package team

import "database/sql"

// Record represents a record in the team table.
type Record struct {
	ID         int
	InviteCode string
}

// Selector can be used to select a team from the databases with a given invite
// code.
type Selector struct{ db *sql.DB }

// NewSelector creates and returns a new Selector.
func NewSelector(db *sql.DB) Selector { return Selector{db: db} }

// Select selects a team from the database with a given code.
func (s Selector) Select(inviteCode string) (Record, error) {
	rec := Record{InviteCode: inviteCode}
	err := s.db.QueryRow(
		`SELECT id FROM app.team WHERE inviteCode = $1`, inviteCode,
	).Scan(&rec.ID)
	return rec, err
}
