package board

import (
	"database/sql"
	"errors"
)

// RecursiveRecord can be used for recursive board data, which means board data
// and data from each column that belong to that board, each task that belong
// to those columns, and each subtask that belong to those tasks.
type RecursiveRecord struct {
	ID      int
	Name    string
	TeamID  int
	Columns []Column
}

// Column encapsulates the data for each column in RecursiveRecord.
type Column struct {
	ID    int
	Order int
	Tasks []Task
}

// Task encapsulates the data for each task in Column.
type Task struct {
	ID          int
	Title       string
	Description *string
	Order       int
	Subtasks    []Subtask
}

// Subtask encapsulates the data for each subtask in Task.
type Subtask struct {
	ID     int
	Title  string
	Order  int
	IsDone bool
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
func (r RecursiveSelector) Select(id string) (RecursiveRecord, error) {
	var res RecursiveRecord

	// Select board.
	err := r.db.QueryRow(
		"SELECT id, name, teamID FROM app.board WHERE id = $1", id,
	).Scan(&res.ID, &res.Name, &res.TeamID)
	if err != nil {
		return RecursiveRecord{}, err
	}

	// Select each column that belongs to the board.
	columnRows, err := r.db.Query(
		`SELECT id, "order" FROM app.column WHERE boardID = $1`, res.ID,
	)
	if err != nil {
		return RecursiveRecord{}, err
	}
	for columnRows.Next() {
		var col Column
		if err = columnRows.Scan(&col.ID, &col.Order); err != nil {
			return RecursiveRecord{}, err
		}

		// Select each task for each column.
		taskRows, err := r.db.Query(
			`SELECT id, title, description, "order" FROM app.task `+
				`WHERE columnID = $1`,
			col.ID,
		)
		if errors.Is(err, sql.ErrNoRows) {
			res.Columns = append(res.Columns, col)
			continue
		}
		if err != nil {
			return RecursiveRecord{}, err
		}

		for taskRows.Next() {
			var task Task
			if err = taskRows.Scan(
				&task.ID, &task.Title, &task.Description, &task.Order,
			); err != nil {
				return RecursiveRecord{}, err
			}

			subtaskRows, err := r.db.Query(
				`SELECT id, title, "order", isDone FROM app.subtask `+
					`WHERE taskID = $1`,
				task.ID,
			)
			if errors.Is(err, sql.ErrNoRows) {
				col.Tasks = append(col.Tasks, task)
				continue
			}
			if err != nil {
				return RecursiveRecord{}, err
			}

			for subtaskRows.Next() {
				var subtask Subtask
				if err = subtaskRows.Scan(
					&subtask.ID,
					&subtask.Title,
					&subtask.Order,
					&subtask.IsDone,
				); err != nil {
					return RecursiveRecord{}, err
				}

				task.Subtasks = append(task.Subtasks, subtask)
			}

			col.Tasks = append(col.Tasks, task)
		}

		res.Columns = append(res.Columns, col)
	}

	return res, nil
}
