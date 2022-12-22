package db

import (
	"database/sql"
	_ "github.com/lib/pq"
)

// ExistorUser is a type that checks whether a user with a given username exists
// in the database.
type ExistorUser struct{ db *sql.DB }

// NewExistorUser is the constructor for ExistorUser.
func NewExistorUser(db *sql.DB) *ExistorUser { return &ExistorUser{db: db} }

// Exists checks whether a user exists in the database, returning true if it
// does and false if not. It also returns any errors encountered during database
// query.
func (c *ExistorUser) Exists(username string) (bool, error) {
	switch err := c.db.
		QueryRow(`SELECT username FROM users WHERE username = $1`, username).
		Err(); err {
	case nil:
		return true, nil
	case sql.ErrNoRows:
		return false, nil
	default:
		return false, err
	}
}
