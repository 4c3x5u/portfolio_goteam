package db

import (
	"database/sql"
	"log"
)

// Closer describes a type that can be used to close something.
type Closer interface{ Close() }

// ConnCloser can be used to close database connections.
type ConnCloser struct{ db *sql.DB }

// NewConnCloser creates and returns a new ConnCloser.
func NewConnCloser(db *sql.DB) ConnCloser { return ConnCloser{db: db} }

// Close closes the database connections associated with the ConnCloser.
func (c ConnCloser) Close() {
	if err := c.db.Close(); err != nil {
		log.Println("Connection could not be closed.\n  ERR:", err)
	}
}
