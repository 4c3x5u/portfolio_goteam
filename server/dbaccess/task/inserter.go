package task

import (
	"context"
	"database/sql"
)

// Task describes the data needed to insert a task into the database. It doesn't
// represent the final record in the task table.
type Task struct{}

// NewTask creates and returns a new Task.
func NewTask() Task { return Task{} }

// Inserter can be used to create a new record in the task table.
type Inserter struct{ db *sql.DB }

// NewInserter creates and returns a new Inserter.
func NewInserter(db *sql.DB) Inserter { return Inserter{db: db} }

// Insert creates a new record in the user table.
func (i Inserter) Insert(_ Task) error {
	ctx := context.Background()
	_, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	return nil
}
