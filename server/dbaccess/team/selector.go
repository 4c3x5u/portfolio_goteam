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
func (s Selector) Select(id string) (rec Record, err error) {
	err = s.db.QueryRow(
		`SELECT id, inviteCode FROM app.team WHERE id = $1`, id,
	).Scan(&rec.ID, &rec.InviteCode)
	return
}

// SelectorByInvCode can be used to select a team from the databases with a
// given invite code.
type SelectorByInvCode struct{ db *sql.DB }

// NewSelectorByInvCode creates and returns a new SelectorByInvCode.
func NewSelectorByInvCode(db *sql.DB) SelectorByInvCode {
	return SelectorByInvCode{db: db}
}

// Select selects a team from the database with a given code.
func (s SelectorByInvCode) Select(inviteCode string) (rec Record, err error) {
	err = s.db.QueryRow(
		`SELECT id, inviteCode FROM app.team WHERE inviteCode = $1`, inviteCode,
	).Scan(&rec.ID, &rec.InviteCode)
	return
}
