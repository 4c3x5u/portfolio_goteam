// Package team contains code for working with the team DynamoDB table.
package team

// tableName is the name of the environment variable to retrieve the team
// table's name from.
const tableName = "DYNAMODB_TABLE_TEAM"

// Team defines the team entity - the primary entity of team domain.
type Team struct {
	ID      string   //uuid
	Members []string //usernames
	Boards  []Board
}

// NewTeam creates and returns a new team.
func NewTeam(id string, members []string, boards []Board) Team {
	return Team{ID: id, Members: members, Boards: boards}
}

// Board defines the board entity which a team may own one/many of.
type Board struct {
	ID   string `json:"id"` //uuid
	Name string `json:"name"`
}

// NewBoard creates and returns a new board.
func NewBoard(id, name string) Board { return Board{ID: id, Name: name} }
