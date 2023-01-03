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

// ReaderSession is a type that can be used to read sesions from the database.
type ReaderSession struct{ db *sql.DB }

// NewReaderSession is the constructor for ReaderSession.
func NewReaderSession(db *sql.DB) *ReaderSession {
	return &ReaderSession{db: db}
}

// Read uses the username of a user to read the session associated with that
// user from the database.
func (r *ReaderSession) Read(username string) (*Session, error) {
	session := NewSession("", "", time.Now())
	err := r.db.
		QueryRow(`SELECT id, username, expiry FROM sessions WHERE username = $1`, username).
		Scan(&session.ID, &session.Username, &session.Expiry)
	return session, err
}
