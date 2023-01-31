package db

import (
	"database/sql"

	"server/log"
)

// Closer describes a type that can be used to close something.
type Closer interface{ Close() }

// ConnCloser can be used to close database connections.
type ConnCloser struct {
	db     *sql.DB
	logger log.Logger
}

// NewConnCloser creates and returns a new ConnCloser.
func NewConnCloser(db *sql.DB, logger log.Logger) ConnCloser {
	return ConnCloser{db: db, logger: logger}
}

// Close closes the database connections associated with the ConnCloser.
func (c ConnCloser) Close() {
	if err := c.db.Close(); err != nil {
		c.logger.Log(log.LevelError, err.Error())
	}
}
