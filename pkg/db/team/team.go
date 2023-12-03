// Package team contains code for working with the team DynamoDB table.
package team

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
