package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// User represents a row in the users table inside the database.
type User struct {
	Username string
	Password []byte
}

// NewUser is the constructor for User.
func NewUser(username string, password []byte) *User {
	return &User{Username: username, Password: password}
}

// CreatorUser is a type that creates a user in the database with the given
// Username and password.
type CreatorUser struct{ db *sql.DB }

// NewCreatorUser is the constructor for CreatorUser.
func NewCreatorUser(db *sql.DB) *CreatorUser { return &CreatorUser{db: db} }

// Create creates a user in the database with the given Username and password.
func (c *CreatorUser) Create(user *User) error {
	_, err := c.db.Exec(
		"INSERT INTO users(username, password) VALUES ($1, $2) "+
			"ON CONFLICT (username) DO NOTHING",
		user.Username, string(user.Password),
	)
	return err
}

// ReaderUser is a type that can be used to read users (i.e. records from users
// table) from the database.
type ReaderUser struct{ db *sql.DB }

// NewReaderUser is the constructor for ReaderUser.
func NewReaderUser(db *sql.DB) *ReaderUser { return &ReaderUser{db: db} }

// Read uses the username of a user to read their password from the database.
func (r *ReaderUser) Read(username string) (*User, error) {
	user := NewUser("", []byte{})
	err := r.db.
		QueryRow(`SELECT username, password FROM users WHERE username = $1`, username).
		Scan(&user.Username, &user.Password)
	return user, err
}
