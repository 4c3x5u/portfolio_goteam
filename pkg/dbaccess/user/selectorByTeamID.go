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
func (s SelectorByTeamID) Select(teamID string) (records []Record, err error) {
	rows, err := s.db.Query(
		`SELECT username, isAdmin FROM app."user" WHERE teamID = $1`,
		teamID,
	)
	if err != nil {
		return
	}

	for rows.Next() {
		var rec Record
		if err = rows.Scan(&rec.Username, &rec.IsAdmin); err != nil {
			return
		}
		records = append(records, rec)
	}

	return
}
