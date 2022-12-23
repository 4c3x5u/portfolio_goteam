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

// CreatorUser is a type that creates a user in the database with the given
// username and password.
type CreatorUser struct {
	db *sql.DB
}

// NewCreatorUser is the constructor for CreatorUser.
func NewCreatorUser(db *sql.DB) *CreatorUser {
	return &CreatorUser{db: db}
}

// Create creates a user in the database with the given username and password.
func (c *CreatorUser) Create(args ...any) error {
	_, err := c.db.Exec(
		`INSERT INTO users(username, password) VALUES ($1, $2)`,
		args,
	)
	return err
}
