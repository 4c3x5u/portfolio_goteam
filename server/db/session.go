package db

import (
	"database/sql"
	"time"
)

// Session represents a record in the sessions table inside the database.
type Session struct {
	ID       string
	Username string
	Expiry   time.Time
}

// NewSession is the constructor for Session.
func NewSession(id, username string, expiry time.Time) *Session {
	return &Session{ID: id, Username: username, Expiry: expiry}
}

// CreatorSession inserts a new sesison to the database.
type CreatorSession struct{ db *sql.DB }

// NewCreatorSession is the constructor for CreatorSession.
func NewCreatorSession(db *sql.DB) *CreatorSession {
	return &CreatorSession{db: db}
}

// Create creates a session in the database with the given id and username
// strings, and expiry time.
func (c *CreatorSession) Create(session *Session) error {
	_, err := c.db.Exec(
		`INSERT INTO sessions(id, username, expiry) VALUES ($1, $2, $3)`,
		session.ID, session.Username, session.Expiry.String(),
	)
	return err
}

type UpserterSession struct{ db *sql.DB }

func NewUpserterSession(db *sql.DB) *UpserterSession {
	return &UpserterSession{db: db}
}

func (u *UpserterSession) Upsert(session *Session) error {
	_, err := u.db.Exec(
		`INSERT INTO sessions(id, username, expiry) VALUES ($1, $2, $3) ON CONFLICT username DO UPDATE SET expiry = $3`,
		session.ID, session.Username, session.Expiry.String(),
	)
	return err
}
