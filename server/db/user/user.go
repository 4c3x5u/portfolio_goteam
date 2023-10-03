// Package user contains code for working with the user table that is inside
// the goteam database.
package user

// Record represents a record in the user table.
type Record struct {
	Username string
	Password []byte
}

// NewRecord creates and returns a new Record.
func NewRecord(username string, password []byte) Record {
	return Record{Username: username, Password: password}
}
