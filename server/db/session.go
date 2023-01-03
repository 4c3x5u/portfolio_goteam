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

// CreatorSession is a type that can be used to insert a new sesison into the
// database.
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

// UpserterSession is a type that can be used to insert or update a session in
// the database.
type UpserterSession struct{ db *sql.DB }

// NewUpserterSession is the constructor for UpserterSession.
func NewUpserterSession(db *sql.DB) *UpserterSession {
	return &UpserterSession{db: db}
}

// Upsert inserts or updates a session in the database. It updates if a session
// with the same username already exists in the database, and inserts if not.
func (u *UpserterSession) Upsert(session *Session) error {
	_, err := u.db.Exec(
		`INSERT INTO sessions(id, username, expiry) VALUES ($1, $2, $3) ON CONFLICT (username) DO UPDATE SET expiry = $3`,
		session.ID, session.Username, session.Expiry.String(),
	)
	return err
}
