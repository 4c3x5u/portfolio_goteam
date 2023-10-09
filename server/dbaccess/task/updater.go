package task

import (
	"context"
	"database/sql"
	"errors"
)

// UpRecord describes the data needed to update a task in the database.
type UpRecord struct {
	Title       string
	Description string
	Subtasks    []Subtask
}

// NewUpRecord creates and returns a new UpRecord.
func NewUpRecord(title, description string, subtasks []Subtask) UpRecord {
	return UpRecord{
		Title:       title,
		Description: description,
		Subtasks:    subtasks,
	}
}

// Subtask defines the Subtasks field of UpRecord.
type Subtask struct {
	Title  string
	Order  int
	IsDone bool
}

// NewSubtask creates and returns a new Subtask.
func NewSubtask(title string, order int, isDone bool) Subtask {
	return Subtask{Title: title, Order: order, IsDone: isDone}
}

// Updater can be used to update a task and its subtasks in the database.
type Updater struct{ db *sql.DB }

// NewUpdater creates and returns a new Updater.
func NewUpdater(db *sql.DB) Updater { return Updater{db: db} }

// Update updates the task in the database with the given ID, as well as
// updating its subtasks, based on data from rec.
func (u Updater) Update(id string, rec UpRecord) error {
	ctx := context.Background()
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Update the task's title and description.
	// No rollback needed because nothing else is touched yet.
	if _, err = tx.Exec(
		"UPDATE app.task SET title = $1, description = $2 WHERE id = $3",
		rec.Title, rec.Description, id,
	); err != nil {
		return err
	}

	// Delete all subtasks associated with this task,
	// because we have the new subtask list at hand anyway.
	if _, err = tx.ExecContext(
		ctx, "DELETE FROM app.subtask WHERE taskID = $1", id,
	); err != nil {
		if rollbackErr := tx.Rollback(); rollbackErr != nil {
			return errors.Join(err, rollbackErr)
		}
		return err
	}

	// Create subtasks for this task based on the Subtasks field of rec.
	for _, subtask := range rec.Subtasks {
		if _, err = tx.ExecContext(
			ctx,
			`INSERT INTO app.subtask(taskID, title, "order", isDone)`+
				`VALUES($1, $2, $3, $4)`,
			id, subtask.Title, subtask.Order, subtask.IsDone,
		); err != nil {
			if rollbackErr := tx.Rollback(); rollbackErr != nil {
				return errors.Join(err, rollbackErr)
			}
			return err
		}
	}

	return tx.Commit()
}
