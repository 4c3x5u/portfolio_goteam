// Package team contains code for accessing the team DynamoDB table.
package team

// Team defines an item in the team table.
type Team struct {
	ID      string   //guid
	Members []string //usernames
	Boards  []Board
}

// Board defines an item in a Team item's Boards.
type Board struct {
	ID   string //guid
	Name string
}
