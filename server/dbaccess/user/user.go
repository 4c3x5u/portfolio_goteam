// Package user contains code for working with the user table.
package user

// Record represents a record in the user table.
type Record struct {
	Username string
	Password []byte
	TeamID   int
	IsAdmin  bool
}

// NewRecord creates and returns a new InRecord.
func NewRecord(
	username string, password []byte, teamID int, isAdmin bool,
) Record {
	return Record{
		Username: username,
		Password: password,
		TeamID:   teamID,
		IsAdmin:  isAdmin,
	}
}
