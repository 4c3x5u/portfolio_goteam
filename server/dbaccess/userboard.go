package dbaccess

import "database/sql"

// UserBoardSelector can be used to select a record from the user_board table.
type UserBoardSelector struct{ db *sql.DB }

// NewUserBoardSelector creates and returns a new UserBoardSelector.
func NewUserBoardSelector(db *sql.DB) UserBoardSelector {
	return UserBoardSelector{db: db}
}

// Select selects a record from the user_board table. It only returns the
// isAdmin field since that is the only piece of information required for the
// use cases.
func (s UserBoardSelector) Select(
	username, boardID string,
) (isAdmin bool, err error) {
	err = s.db.QueryRow(
		"SELECT isAdmin FROM app.user_board WHERE username = $1 "+
			"AND boardID = $2",
		username,
		boardID,
	).Scan(&isAdmin)
	return
}

// UserBoardCounter can be used to count the number of boards in the database
// that a certain user is the admin to.
type UserBoardCounter struct{ db *sql.DB }

// NewUserBoardCounter creates and returns a new UserBoardCounter.
func NewUserBoardCounter(db *sql.DB) UserBoardCounter {
	return UserBoardCounter{db: db}
}

// Count counts the number of boards in the database that the user with the
// given username is the admin to.
func (c UserBoardCounter) Count(username string) (count int, err error) {
	err = c.db.QueryRow(
		"SELECT COUNT(*) FROM app.user_board "+
			"WHERE username = $1 AND isAdmin = TRUE",
		username,
	).Scan(&count)
	return
}
