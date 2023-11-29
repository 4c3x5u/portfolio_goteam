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
		assert.Nil(t.Fatal, err)

		assert.Equal(t.Error, res.ID, 1)
		assert.Equal(t.Error, res.Name, "board 1")
		assert.Equal(t.Error, res.TeamID, 21)

		columns := res.Columns
		assert.Equal(t.Error, len(columns), 4)
		for i := 0; i < 2; i++ {
			assert.Equal(t.Error, columns[i].ID, i+2)
			assert.Equal(t.Error, columns[i].Order, i+1)
			assert.Equal(t.Error, len(columns[i].Tasks), 0)
		}

		column3 := res.Columns[2]
		assert.Equal(t.Error, len(column3.Tasks), 1)

		task1 := column3.Tasks[0]
		assert.Equal(t.Error, task1.ID, 6)
		assert.Equal(t.Error, task1.Title, "task 1")
		assert.Equal(t.Error, *task1.Description, "do things!")
		assert.Equal(t.Error, task1.Order, 1)

		column4 := res.Columns[3]
		assert.Equal(t.Error, len(column4.Tasks), 1)

		task2 := column4.Tasks[0]
		assert.Equal(t.Error, task2.ID, 7)
		assert.Equal(t.Error, task2.Title, "task 2")
		assert.Equal(t.Error, *task2.Description, "do things!")
		assert.Equal(t.Error, task2.Order, 2)

		subtasks := task2.Subtasks
		assert.Equal(t.Error, len(subtasks), 2)

		subtask1 := subtasks[0]
		assert.Equal(t.Error, subtask1.ID, 8)
		assert.Equal(t.Error, subtask1.Title, "subtask 1")
		assert.Equal(t.Error, subtask1.Order, 1)
		if err = assert.True(!subtask1.IsDone); err != nil {
			t.Error(err)
		}

		subtask2 := subtasks[1]
		assert.Equal(t.Error, subtask2.ID, 9)
		assert.Equal(t.Error, subtask2.Title, "subtask 2")
		assert.Equal(t.Error, subtask2.Order, 2)
		if err := assert.True(subtask2.IsDone); err != nil {
			t.Error(err)
		}

	})
}
