package board

import "database/sql"

// SelectorByTeamID can be used to select a record from the board table that
// matches a team ID.
type SelectorByTeamID struct{ db *sql.DB }

// NewSelectorByTeamID creates and returns a new SelectorByTeamID.
func NewSelectorByTeamID(db *sql.DB) SelectorByTeamID {
	return SelectorByTeamID{db: db}
}

// Select selects a record from the board table that matches the given team ID.
func (s SelectorByTeamID) Select(teamID string) (rec Record, err error) {
	_, err = s.db.Query(
		"SELECT id, name FROM app.board WHERE teamID = $1",
		teamID,
	)
	return
}
