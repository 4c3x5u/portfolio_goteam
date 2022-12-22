package register

import (
	"database/sql"
	"errors"

	_ "github.com/lib/pq"
)

// CreatorUser defines the signature for a type that creates a user based on
// a password and a password. This abstraction is for thinning the client that
// the handler depends on by getting it to rely on this rather than a *sql.DB,
// as well as making the handler code more readable.
type CreatorUser interface {
	CreateUser(username, password string) error
}

// CreatorDBUser is a type that is used to create a user in the database
type CreatorDBUser struct {
	db *sql.DB
}

// NewCreatorDBUser is the constructor for CreatorDBUser.
func NewCreatorDBUser(db *sql.DB) *CreatorDBUser {
	return &CreatorDBUser{db: db}
}

// errCreatorUsernameTaken is returned from CreateUser when the username that
// the user picked is already taken by another user.
var errCreatorUsernameTaken = errors.New("CreatorDBUser.CreateUser: username taken")

// CreateUser creates a new user in the database.
//   - The first return value tells whether the username exists (e.g. true if it
//     does).
//   - The second return value is the error returned from the call to create the
//     user in the database.
func (c *CreatorDBUser) CreateUser(username, _ string) error {
	switch err := c.db.
		QueryRow(`SELECT username FROM users WHERE username = $1`, username).
		Err(); err {
	case nil:
		return errCreatorUsernameTaken
	case sql.ErrNoRows:
		return nil
	default:
		return err
	}
}
