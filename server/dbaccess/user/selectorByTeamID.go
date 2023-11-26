package user

import "database/sql"

// SelectorByTeamID can be used to retrieve all records from the user table that
// matches a team ID.
type SelectorByTeamID struct{ db *sql.DB }

// NewSelectorByTeamID creates and returns a new SelectorByTeamID.
func NewSelectorByTeamID(db *sql.DB) SelectorByTeamID {
	return SelectorByTeamID{db: db}
}

// Select selects all records from the user table that matches the team ID.
func (s SelectorByTeamID) Select(teamID string) ([]Record, error) {
	_, err := s.db.Query(
		`SELECT username, isAdmin FROM app."user" WHERE teamID = $1`,
		teamID,
	)
	if err != nil {
		return nil, err
	}

	return nil, nil
}
