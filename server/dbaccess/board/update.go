package board

import (
	"database/sql"
	"errors"
)

// Updater can be used to update the name field of record in the board
// table.
type Updater struct{ db *sql.DB }

// NewUpdater is the constructor for Updater.
func NewUpdater(db *sql.DB) Updater { return Updater{db: db} }

// Update updates the name field of a record in the board database with a new
// value.
func (u Updater) Update(id, newName string) error {
	res, err := u.db.Exec(
		"UPDATE app.board SET name = $1 WHERE id = $2", newName, id,
	)
	if err != nil {
		return err
	}
	if rowsAffected, err := res.RowsAffected(); err != nil {
		return err
	} else if rowsAffected == 0 {
		return errors.New("no rows were affected")
	} else if rowsAffected > 1 {
		return errors.New("more than expected rows were affected")
	}
	return nil
}
