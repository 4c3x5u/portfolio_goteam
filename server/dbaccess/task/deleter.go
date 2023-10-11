package task

import "database/sql"

// Deleter can be used to delete a record from the task table.
type Deleter struct{ db *sql.DB }

// NewDeleter creates and returns a new Deleter.
func NewDeleter(db *sql.DB) Deleter { return Deleter{db: db} }

// Delete deletes a record from the task table with the given ID.
func (d Deleter) Delete(id string) error {
	_, err := d.db.Exec("DELETE FROM app.task WHERE id = $1", id)
	return err
}
