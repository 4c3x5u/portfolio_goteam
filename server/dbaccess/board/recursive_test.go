//go:build utest

package board

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
)

// TestRecursiveSelector tests the Select method of RecursiveSelector to assert
// that it sends the correct queries to the database with the correct arguments,
// and returns whatever error occurs alongside the correct data.
func TestRecursiveSelector(t *testing.T) {
	boardID := "21"
	sqlSelectBoard := `SELECT id, name, teamID FROM app.board WHERE id = \$1`
	sqlSelectColumn := `SELECT id, \"order\" FROM app.column ` +
		`WHERE boardID = \$1`
	sqlSelectTask := `SELECT id, title, description, \"order\" FROM app.task ` +
		`WHERE columnID = \$1`
	sqlSelectSubtask := `SELECT id, title, \"order\", isDone ` +
		`FROM app.subtask WHERE taskID = \$1`

	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer func() {
		mock.ExpectClose()
		teardown()
	}()

	sut := NewRecursiveSelector(db)

	t.Run("SelectBoardErr", func(t *testing.T) {
		wantErr := errors.New("error selecting board")
		mock.ExpectQuery(sqlSelectBoard).
			WithArgs(boardID).
			WillReturnError(wantErr)

		_, err := sut.Select(boardID)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("SelectColumnErr", func(t *testing.T) {
		wantErr := errors.New("error selecting column")
		mock.ExpectQuery(sqlSelectBoard).
			WithArgs(boardID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "teamID"}).
					AddRow(1, "board 1", 1),
			)
		mock.ExpectQuery(sqlSelectColumn).
			WithArgs(1).
			WillReturnError(wantErr)

		_, err := sut.Select(boardID)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("SelectTaskErr", func(t *testing.T) {
		wantErr := errors.New("error selecting task")
		mock.ExpectQuery(sqlSelectBoard).
			WithArgs(boardID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "teamID"}).
					AddRow(1, "board 1", 1),
			)
		mock.ExpectQuery(sqlSelectColumn).
			WithArgs(1).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "order"}).
					AddRow(2, 1).AddRow(3, 2).AddRow(4, 3).AddRow(5, 4),
			)
		mock.ExpectQuery(sqlSelectTask).
			WithArgs(2).
			WillReturnError(wantErr)

		_, err := sut.Select(boardID)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("SelectSubtaskErr", func(t *testing.T) {
		wantErr := errors.New("error selecting subtask")
		mock.ExpectQuery(sqlSelectBoard).
			WithArgs(boardID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "teamID"}).
					AddRow(1, "board 1", 1),
			)
		mock.ExpectQuery(sqlSelectColumn).
			WithArgs(1).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "order"}).
					AddRow(2, 1).AddRow(3, 2).AddRow(4, 3).AddRow(5, 4),
			)
		mock.ExpectQuery(sqlSelectTask).
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectQuery(sqlSelectTask).
			WithArgs(3).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectQuery(sqlSelectTask).
			WithArgs(4).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectQuery(sqlSelectTask).
			WithArgs(5).
			WillReturnRows(
				sqlmock.NewRows(
					[]string{"id", "title", "description", "order"},
				).AddRow(6, "task 1", "do things!", 1),
			)
		mock.ExpectQuery(sqlSelectSubtask).
			WithArgs(6).
			WillReturnError(wantErr)

		_, err := sut.Select(boardID)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		mock.ExpectQuery(sqlSelectBoard).
			WithArgs(boardID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "teamID"}).
					AddRow(1, "board 1", 21),
			)
		mock.ExpectQuery(sqlSelectColumn).
			WithArgs(1).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "order"}).
					AddRow(2, 1).AddRow(3, 2).AddRow(4, 3).AddRow(5, 4),
			)
		mock.ExpectQuery(sqlSelectTask).
			WithArgs(2).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectQuery(sqlSelectTask).
			WithArgs(3).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectQuery(sqlSelectTask).
			WithArgs(4).
			WillReturnRows(
				sqlmock.NewRows(
					[]string{"id", "title", "description", "order"},
				).AddRow(6, "task 1", "do things!", 1),
			)
		mock.ExpectQuery(sqlSelectSubtask).
			WithArgs(6).
			WillReturnError(sql.ErrNoRows)
		mock.ExpectQuery(sqlSelectTask).
			WithArgs(5).
			WillReturnRows(
				sqlmock.NewRows(
					[]string{"id", "title", "description", "order"},
				).AddRow(7, "task 2", "do things!", 2),
			)
		mock.ExpectQuery(sqlSelectSubtask).
			WithArgs(7).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "title", "order", "isDone"}).
					AddRow(8, "subtask 1", 1, false).
					AddRow(9, "subtask 2", 2, true),
			)

		res, err := sut.Select(boardID)

		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(1, res.ID); err != nil {
			t.Error(err)
		}
		if err = assert.Equal("board 1", res.Name); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(21, res.TeamID); err != nil {
			t.Error(err)
		}

		columns := res.Columns
		if err = assert.Equal(4, len(columns)); err != nil {
			t.Error(err)
		}
		for i := 0; i < 2; i++ {
			if err = assert.Equal(i+2, columns[i].ID); err != nil {
				t.Error(err)
			}
			if err = assert.Equal(i+1, columns[i].Order); err != nil {
				t.Error(err)
			}
			if err = assert.Equal(0, len(columns[i].Tasks)); err != nil {
				t.Error(err)
			}
		}

		column3 := res.Columns[2]
		if err = assert.Equal(1, len(column3.Tasks)); err != nil {
			t.Error(err)
		}

		task1 := column3.Tasks[0]
		if err = assert.Equal(6, task1.ID); err != nil {
			t.Error(err)
		}
		if err = assert.Equal("task 1", task1.Title); err != nil {
			t.Error(err)
		}
		if err = assert.Equal("do things!", task1.Description); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(1, task1.Order); err != nil {
			t.Error(err)
		}

		column4 := res.Columns[3]
		if err = assert.Equal(1, len(column4.Tasks)); err != nil {
			t.Error(err)
		}

		task2 := column4.Tasks[0]
		if err = assert.Equal(7, task2.ID); err != nil {
			t.Error(err)
		}
		if err = assert.Equal("task 2", task2.Title); err != nil {
			t.Error(err)
		}
		if err = assert.Equal("do things!", task2.Description); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(2, task2.Order); err != nil {
			t.Error(err)
		}

		subtasks := task2.Subtasks
		if err = assert.Equal(2, len(subtasks)); err != nil {
			t.Error(err)
		}

		subtask1 := subtasks[0]
		if err = assert.Equal(8, subtask1.ID); err != nil {
			t.Error(err)
		}
		if err = assert.Equal("subtask 1", subtask1.Title); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(1, subtask1.Order); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(false, subtask1.IsDone); err != nil {
			t.Error(err)
		}

		subtask2 := subtasks[1]
		if err = assert.Equal(9, subtask2.ID); err != nil {
			t.Error(err)
		}
		if err = assert.Equal("subtask 2", subtask2.Title); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(2, subtask2.Order); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(true, subtask2.IsDone); err != nil {
			t.Error(err)
		}

	})
}
