package task

import (
	"context"
	"database/sql"
	"errors"
)

// Task describes the data needed to insert a task into the database. It doesn't
// represent the final record in the task table.
type Task struct {
	columnID      int
	title         string
	description   string
	subtaskTitles []string
}

// NewTask creates and returns a new Task.
func NewTask(
	columnID int, title string, description string, subtaskTitles []string,
) Task {
	return Task{
		columnID:      columnID,
		title:         title,
		description:   description,
		subtaskTitles: subtaskTitles,
	}
}

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

	// Get the task with the highest order that is associated with the same
	// column that the new task is.
	var highestOrder int
	if err = tx.QueryRowContext(
		ctx,
		`SELECT "order" FROM app.task WHERE columnID = $1 `+
			`ORDER BY "order" DESC LIMIT 1`,
		task.columnID,
	).Scan(&highestOrder); errors.Is(err, sql.ErrNoRows) {
		highestOrder = 0
	} else if err != nil {
		return err
	}

	// Insert a record into task table with the given data and the order of 1
	// higher than the highest order found in this table.
	_, err = tx.ExecContext(
		ctx,
		`INSERT INTO app.task(columnID, title, description, "order")`+
			`VALUES ($1, $2, $3, $4)`,
		task.columnID, task.title, task.description, highestOrder+1,
	)
	if err != nil {
		return err
	}

	return nil
}
