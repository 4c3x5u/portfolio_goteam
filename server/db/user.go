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

// UserCreator is a type that creates a user in the database with the given
// Username and password.
type UserCreator struct{ db *sql.DB }

// NewUserCreator is the constructor for UserCreator.
func NewUserCreator(db *sql.DB) *UserCreator { return &UserCreator{db: db} }

// Create creates a user in the database with the given Username and password.
func (c *UserCreator) Create(user *User) error {
	_, err := c.db.Exec(
		"INSERT INTO users(username, password) VALUES ($1, $2)",
		user.Username, string(user.Password),
	)
	return err
}

// UserReader is a type that can be used to read users (i.e. records from users
// table) from the database.
type UserReader struct{ db *sql.DB }

// NewUserReader is the constructor for UserReader.
func NewUserReader(db *sql.DB) *UserReader { return &UserReader{db: db} }

// Read uses the username of a user to read their password from the database.
func (r *UserReader) Read(username string) (*User, error) {
	user := NewUser("", []byte{})
	err := r.db.
		QueryRow(`SELECT username, password FROM users WHERE username = $1`, username).
		Scan(&user.Username, &user.Password)
	return user, err
}
