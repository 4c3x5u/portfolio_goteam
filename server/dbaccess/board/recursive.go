package board

import "database/sql"

// RecursiveBoard can be used for recursive board data, which means board data
// and data from each column that belong to that board, each task that belong
// to those columns, and each subtask that belong to those tasks.
type RecursiveBoard struct {
	ID   int
	Name string
}

// RecursiveSelector can be used to select a record from the board table, as well
// as all the columns that belong to the board, all the tasks that belong to
// those columns, and all the subtasks that belong to those tasks.
type RecursiveSelector struct{ db *sql.DB }

// NewRecursiveSelector creates and returns a new RecursiveSelector.
func NewRecursiveSelector(db *sql.DB) RecursiveSelector {
	return RecursiveSelector{db: db}
}

// Select selects a record from the board table, as well as all the columns that
// belong to the board, all the tasks that belong to those columns, and all the
// subtasks that belong to those tasks.
func (r RecursiveSelector) Select(id string) (RecursiveBoard, error) {
	var res RecursiveBoard

	err := r.db.QueryRow(
		"SELECT id, name FROM app.board WHERE id = $1", id,
	).Scan(&res.ID, &res.Name)
	if err != nil {
		return RecursiveBoard{}, err
	}

	return res, nil
}
