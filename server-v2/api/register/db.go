package register

import (
	"database/sql"
	_ "github.com/lib/pq"
)

// CreatorUser defines the signature for a type that creates a user based on
// a password and a password.
type CreatorUser interface {
	CreateUser(username, password string) (usernameIsTaken bool, err error)
}

// CreatorDBUser is a type that is used to create a user in the database
type CreatorDBUser struct {
	db *sql.DB
}

// NewCreatorDBUser is the constructor for CreatorDBUser.
func NewCreatorDBUser(db *sql.DB) *CreatorDBUser {
	return &CreatorDBUser{db: db}
}

// CreateUser creates a new user in the database.
func (c *CreatorDBUser) CreateUser(username, _ string) (bool, error) {
	switch err := c.db.
		QueryRow(`SELECT username FROM users WHERE username = $1`, username).
		Err(); err {
	case sql.ErrNoRows:
		return false, nil
	case nil:
		return true, nil
	default:
		return false, err
	}
}
