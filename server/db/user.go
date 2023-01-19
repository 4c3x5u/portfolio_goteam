package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// User represents a record in the user table.
type User struct {
	Username string
	Password []byte
}

// NewUser creates and returns a new User.
func NewUser(username string, password []byte) User {
	return User{Username: username, Password: password}
}

// UserInserter can be used to create a new record in the user table.
type UserInserter struct{ db *sql.DB }

// NewUserInserter creates and returns a new UserInserter.
func NewUserInserter(db *sql.DB) UserInserter { return UserInserter{db: db} }

// Insert creates a new record in the user table.
func (c UserInserter) Insert(user User) error {
	_, err := c.db.Exec(
		`INSERT INTO app."user"(username, password) VALUES ($1, $2)`,
		user.Username, string(user.Password),
	)
	return err
}

// UserSelector can be used to read records from the user table.
type UserSelector struct{ db *sql.DB }

// NewUserSelector creates and returns a new user reader.
func NewUserSelector(db *sql.DB) UserSelector { return UserSelector{db: db} }

// Select selects a record from the user table with the given username.
func (r UserSelector) Select(username string) (user User, err error) {
	err = r.db.QueryRow(
		`SELECT username, password FROM app."user" WHERE username = $1`,
		username,
	).Scan(&user.Username, &user.Password)
	return
}
