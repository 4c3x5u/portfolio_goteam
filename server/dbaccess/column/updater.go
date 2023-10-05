package column

import "database/sql"

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
	_, err := u.db.Exec(
		"UPDATE app.task SET columnID = $1 AND order = $2 WHERE id = $3",
		columnID, tasks[0].Order, tasks[0].ID,
	)
	return err
}
