// Package team contains code for working with the team DynamoDB table.
package team

// tableName is the name of the environment variable to retrieve the team
// table's name from.
const tableName = "DYNAMODB_TABLE_TEAM"

// Team defines a team in the team table.
type Team struct {
	ID      string   //guid
	Members []string //usernames
	Boards  []Board
}

// Board defines a board in a team's boards.
type Board struct {
	ID   string //guid
	Name string
}
