package column

import (
	"context"
	"database/sql"
	"errors"
)

// Task contains data needed to reorder tasks within a column or move tasks
// from one column to another.
type Task struct {
	ID    int
	Order int
}

// Updater can be used to reorder tasks within a column or move tasks from one
// column to another.
type Updater struct{ db *sql.DB }

// NewUpdater creates and returns a new Updater.
func NewUpdater(db *sql.DB) Updater { return Updater{db: db} }

// Update sets the column ID and order of each task to the passed-in column ID
// and the task order respectively, thus reordering the tasks within the column
// with the given ID or moving tasks from one column to another.
func (u Updater) Update(columnID string, tasks []Task) error {
	// Begin a transaction so that no tasks end up with the same order in case
	// updates fail halfway though
	ctx := context.Background()
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer tx.Rollback()

	// Go through each task and update their column ID and order based on the
	// columnID and tasks received.
	for _, task := range tasks {
		res, execErr := tx.ExecContext(
			ctx,
			`UPDATE app.task SET columnID = $1, "order" = $2 `+
				`WHERE id = $3`,
			columnID, task.Order, task.ID,
		)
		if execErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Join(execErr, rollbackErr)
			}
			return execErr
		}
		rowsAffected, rowsAffectedErr := res.RowsAffected()
		if rowsAffectedErr != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Join(rowsAffectedErr, rollbackErr)
			}
			return rowsAffectedErr
		}
		// If no rows were affected when execErr was nil, it means task not
		// found? todo: confirm
		if rowsAffected != int64(1) {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Join(sql.ErrNoRows, rollbackErr)
			}
			return sql.ErrNoRows
		}
	}

	// All went well, commit transaction and return err if occurs.
	return tx.Commit()
}
