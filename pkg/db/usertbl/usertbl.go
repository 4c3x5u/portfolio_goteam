// Package usertbl contains code to interact with the user table in DynamoDB.
package usertbl

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

// NewUser creates and returns a new User,
func NewUser(
	username string, password []byte, isAdmin bool, teamID string,
) User {
	return User{
		Username: username,
		Password: password,
		IsAdmin:  isAdmin,
		TeamID:   teamID,
	}
}
