package db

// Board represents a record in the board table.
type Board struct {
	name    string
	adminID string
}

// NewBoard creates and returns a new Board.
func NewBoard(name string, adminID string) Board {
	return Board{name: name, adminID: adminID}
}
