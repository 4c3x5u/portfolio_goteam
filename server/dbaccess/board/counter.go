package board

import "database/sql"

// Counter can be used to get the number of boards that belong to a given team.
type Counter struct{ db *sql.DB }

// NewCounter creates and returns a new Counter.
func NewCounter(db *sql.DB) Counter { return Counter{db: db} }

// Count returns the number of boards in board table whose teamID field's value
// matches the given team ID.
func (c Counter) Count(teamID string) (int, error) {
	var count int
	err := c.db.QueryRow(
		"SELECT COUNT(*) FROM app.board WHERE teamID = $1",
		teamID,
	).Scan(&count)
	return count, err
}
