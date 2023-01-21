package db

import "database/sql"

// UserBoardCounter can be used to count the number of boards in the database
// that a certain user is the admin to.
type UserBoardCounter struct{ db *sql.DB }

// NewUserBoardCounter creates and returns a new UserBoardCounter.
func NewUserBoardCounter(db *sql.DB) UserBoardCounter {
	return UserBoardCounter{db: db}
}

// Count counts the number of boards in the database that the user with the
// given userID is the admin to.
func (c UserBoardCounter) Count(userID string) (count int) {
	c.db.QueryRow(
		"SELECT COUNT(*) FROM app.user_board WHERE userID = $1 AND isAdmin = $2",
		userID,
		true,
	).Scan(&count)
	return
}
