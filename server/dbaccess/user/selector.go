package user

import "database/sql"

// Selector can be used to read records from the user table.
type Selector struct{ db *sql.DB }

// NewSelector creates and returns a new Selector.
func NewSelector(db *sql.DB) Selector { return Selector{db: db} }

// Select selects a record from the user table with the given username.
func (s Selector) Select(username string) (rec Record, err error) {
	err = s.db.QueryRow(
		`SELECT password, teamID, isAdmin FROM app."user" WHERE username = $1`,
		username,
	).Scan(&rec.Password, &rec.TeamID, &rec.IsAdmin)
	rec.Username = username
	return
}
