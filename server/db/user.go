package db

import (
	"database/sql"

	_ "github.com/lib/pq"
)

// User represents a record in the users table.
type User struct {
	Username string
	Password []byte
}

// NewUser creates and returns a new user.
func NewUser(username string, password []byte) User {
	return User{Username: username, Password: password}
}

// UserCreator can be used to create a new record in the users table.
type UserCreator struct{ db *sql.DB }

// NewUserCreator creates and returns a new user creator.
func NewUserCreator(db *sql.DB) UserCreator { return UserCreator{db: db} }

// Create creates a new record in the users table.
func (c UserCreator) Create(user User) error {
	_, err := c.db.Exec(
		"INSERT INTO users(username, password) VALUES ($1, $2)",
		user.Username, string(user.Password),
	)
	return err
}

// UserReader can be used to read records from the users table.
type UserReader struct{ db *sql.DB }

// NewUserReader creates and returns a new user reader.
func NewUserReader(db *sql.DB) UserReader { return UserReader{db: db} }

// Read reads a record from the users table with the given username.
func (r UserReader) Read(username string) (User, error) {
	user := NewUser("", []byte{})
	err := r.db.
		QueryRow(`SELECT username, password FROM users WHERE username = $1`, username).
		Scan(&user.Username, &user.Password)
	return user, err
}
