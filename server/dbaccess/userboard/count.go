package userboard

import "database/sql"

// Counter can be used to count the number of boards in the database
// that a certain user is the admin to.
type Counter struct{ db *sql.DB }

// NewCounter creates and returns a new Counter.
func NewCounter(db *sql.DB) Counter {
	return Counter{db: db}
}

// Count counts the number of boards in the database that the user with the
// given username is the admin to.
func (c Counter) Count(username string) (count int, err error) {
	err = c.db.QueryRow(
		"SELECT COUNT(*) FROM app.user_board "+
			"WHERE username = $1 AND isAdmin = TRUE",
		username,
	).Scan(&count)
	return
}
