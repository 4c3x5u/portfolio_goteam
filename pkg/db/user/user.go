// Package user contains code for working with the user DynamoDB table.
package user

// tableName is the name of the environment variable to retrieve the user
// table's name from.
const tableName = "DYNAMODB_TABLE_USER"

// User defines the user entity - the primary and only entity of user domain.
type User struct {
	Username string
	Password []byte
	IsAdmin  bool
	TeamID   string
}
