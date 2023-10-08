package task

import (
	"context"
	"database/sql"
)

// Task describes the data needed to insert a task into the database. It doesn't
// represent the final record in the task table.
type Task struct{ columnID int }

// NewTask creates and returns a new Task.
func NewTask(columnID int) Task { return Task{columnID: columnID} }

// Inserter can be used to create a new record in the task table.
type Inserter struct{ db *sql.DB }

// NewInserter creates and returns a new Inserter.
func NewInserter(db *sql.DB) Inserter { return Inserter{db: db} }

// Insert creates a new record in the user table.
func (i Inserter) Insert(task Task) error {
	ctx := context.Background()
	tx, err := i.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	var order int
	if err = tx.QueryRowContext(
		ctx,
		`SELECT "order" FROM app.task WHERE columnID = $1 `+
			`ORDER BY "order" DESC LIMIT 1`,
		task.columnID,
	).Scan(&order); err != nil {
		return err
	}

	return nil
}
