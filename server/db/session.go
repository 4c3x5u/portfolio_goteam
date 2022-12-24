package db

import (
	"database/sql"
	"time"
)

// CreatorSession manages active sessions and provides access to them.
type CreatorSession struct {
	db *sql.DB
}

// NewCreatorSession is the constructor for CreatorSession.
func NewCreatorSession(db *sql.DB) *CreatorSession {
	return &CreatorSession{db: db}
}

// Create creates a new session or renews it if it already exists.
func (c *CreatorSession) Create(args ...any) error {
	err := c.db.
		QueryRow(`SELECT id FROM sessions WHERE id = $1`, args[0]).
		Err()
	if err == nil {
		// session exists
		if _, err := c.db.Exec(
			`DELETE FROM sessions WHERE id = $1;`, args[0],
		); err != nil {
			return err
		}
	}
	if err == nil || err == sql.ErrNoRows {
		// insert session if deleted or doesn't exist
		if _, err := c.db.Exec(
			`INSERT INTO sessions(id, username, expiresAt) VALUES ($1, $2, $3)`,
			args,
		); err != nil {
			return err
		}
	}
	return err
}

// Session holds information about a user session.
type Session struct {
	ID       string
	Username string
	Expiry   time.Time
}

// IsExpired tells us whether the user session is expired.
func (s Session) IsExpired() bool {
	return s.Expiry.Before(time.Now())
}
