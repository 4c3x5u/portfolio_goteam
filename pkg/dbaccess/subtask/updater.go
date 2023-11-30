package subtask

import (
	"database/sql"
)

// UpRecord describes the data needed to update a subtask in the database.
type UpRecord struct{ isDone bool }

// NewUpRecord creates and returns a new UpRecord.
func NewUpRecord(isDone bool) UpRecord { return UpRecord{isDone: isDone} }

// Updater can be used to update a subtask in the database.
type Updater struct{ db *sql.DB }

// NewUpdater creates and returns a new Updater.
func NewUpdater(db *sql.DB) Updater { return Updater{db: db} }

// Update updates the subtask in the database with the given ID.
func (u Updater) Update(id string, rec UpRecord) error {
	_, err := u.db.Exec(
		"UPDATE app.subtask SET isDone = $1 WHERE id = $2",
		rec.isDone, id,
	)
	return err
}
