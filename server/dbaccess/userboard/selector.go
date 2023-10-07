package userboard

import "database/sql"

// Selector can be used to select a record from the user_board table.
type Selector struct{ db *sql.DB }

// NewSelector creates and returns a new Selector.
func NewSelector(db *sql.DB) Selector {
	return Selector{db: db}
}

// Select selects a record from the user_board table. It only returns the
// isAdmin field since that is the only piece of information required for the
// use cases.
func (s Selector) Select(username, boardID string) (isAdmin bool, err error) {
	err = s.db.QueryRow(
		"SELECT isAdmin FROM app.user_board WHERE username = $1 "+
			"AND boardID = $2",
		username,
		boardID,
	).Scan(&isAdmin)
	return
}
