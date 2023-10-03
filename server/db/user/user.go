// Package user contains code for working with the user table that is inside
// the goteam database.
package user

// User represents a record in the user table.
type User struct {
	Username string
	Password []byte
}

// NewUser creates and returns a new User.
func NewUser(username string, password []byte) User {
	return User{Username: username, Password: password}
}
