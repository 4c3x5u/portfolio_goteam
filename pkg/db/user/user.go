// Package user contains code for accessing the user DynamoDB table.
package user

// User defines an item in the user table.
type User struct {
	ID       string
	Password []byte
	IsAdmin  bool
	TeamID   int
}
