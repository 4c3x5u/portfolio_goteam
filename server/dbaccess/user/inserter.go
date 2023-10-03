package user

import "database/sql"

// Inserter can be used to create a new record in the user table.
type Inserter struct{ db *sql.DB }

// NewInserter creates and returns a new Inserter.
func NewInserter(db *sql.DB) Inserter { return Inserter{db: db} }

// Insert creates a new record in the user table.
func (i Inserter) Insert(user User) error {
	_, err := i.db.Exec(
		`INSERT INTO app."user"(username, password) VALUES ($1, $2)`,
		user.Username, string(user.Password),
	)
	return err
}
