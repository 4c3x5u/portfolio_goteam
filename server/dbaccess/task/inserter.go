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
	defer tx.Rollback()

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
	// higher than the highest order found in this table. Return task ID to use
	// is subtask insertions.
	var taskID int
	if err = tx.QueryRowContext(
		ctx,
		`INSERT INTO app.task(columnID, title, description, "order")`+
			`VALUES ($1, $2, $3, $4) RETURNING id`,
		task.columnID, task.title, task.description, highestOrder+1,
	).Scan(&taskID); err != nil {
		return err
	}

	// Insert subtasks for the task into the database.
	for j, subtaskTitle := range task.subtaskTitles {
		order := j + 1
		_, err = tx.ExecContext(
			ctx,
			`INSERT INTO app.subtask(taskID, title, "order", isDone) `+
				`VALUES($1, $2, $3, $4)`,
			taskID, subtaskTitle, order, false,
		)
		if err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Join(err, rollbackErr)
			}
			return err
		}
	}

	return tx.Commit()
}
