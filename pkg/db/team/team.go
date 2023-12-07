// Package team contains code for working with the team DynamoDB table.
package team

// tableName is the name of the environment variable to retrieve the team
// table's name from.
const tableName = "DYNAMODB_TABLE_TEAM"

// Team defines the team entity - the primary entity of team domain.
type Team struct {
	ID      string   //guid
	Members []string //usernames
	Boards  []Board
}

// Board defines the board entity which a team may own one/many of.
type Board struct {
	ID   string //guid
	Name string
}
