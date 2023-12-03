// Package user contains code for working with the user DynamoDB table.
package user

// tableName is the name of the environment variable to retrieve the user
// table's name from.
const tableName = "DYNAMODB_TABLE_TEAM"

// User defines user in the user table.
type User struct {
	ID       string
	Password []byte
	IsAdmin  bool
	TeamID   int
}
