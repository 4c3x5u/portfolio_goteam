package board

import "database/sql"

type Counter struct{ db *sql.DB }

func NewCounter(db *sql.DB) Counter { return Counter{db: db} }

func (c Counter) Count(teamID string) (int, error) {
	var count int
	err := c.db.QueryRow(
		"SELECT COUNT(*) FROM app.board  WHERE teamID = $1",
		teamID,
	).Scan(&count)
	return count, err
}
