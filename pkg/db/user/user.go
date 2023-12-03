// Package user contains code for working with the user DynamoDB table.
package user

// User defines user in the user table.
type User struct {
	ID       string
	Password []byte
	IsAdmin  bool
	TeamID   int
}
