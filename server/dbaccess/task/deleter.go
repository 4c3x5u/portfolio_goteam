package task

import (
	"context"
	"database/sql"
	"errors"
)

// Deleter can be used to delete a record from the task table.
type Deleter struct{ db *sql.DB }

// NewDeleter creates and returns a new Deleter.
func NewDeleter(db *sql.DB) Deleter { return Deleter{db: db} }

// Delete deletes a record from the task table with the given ID.
func (d Deleter) Delete(id string) error {
	ctx := context.Background()
	tx, err := d.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback() }()

	// Delete all subtasks for the task.
	_, err = tx.ExecContext(
		ctx, "DELETE FROM app.subtask WHERE taskID = $1", id,
	)
	if err != nil {
		return err
	}

	// Delete the task.
	_, err = tx.ExecContext(ctx, "DELETE FROM app.task WHERE id = $1", id)
	if err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	return tx.Commit()
}
