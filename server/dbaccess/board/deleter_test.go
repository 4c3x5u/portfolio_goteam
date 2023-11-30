//go:build utest

package board

import (
	"errors"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestDeleter tests the Delete method of Deleter to assert that it sends the
// correct query to the database with the correct arguments, and returns
// whatever error occurs.
func TestDeleter(t *testing.T) {
	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewDeleter(db)

	var (
		sqlSelectColumnIDs = `SELECT id FROM app.\"column\" ` +
			`WHERE boardID = \$1`
		sqlSelectTaskIDs  = `SELECT id FROM app.task WHERE columnID = \$1`
		sqlDeleteSubtasks = `DELETE FROM app.subtask WHERE taskID = \$1`
		sqlDeleteTasks    = `DELETE FROM app.task WHERE columnID = \$1`
		sqlDeleteColumns  = `DELETE FROM app.\"column\" WHERE boardID = \$1`
		sqlDeleteBoard    = `DELETE FROM app.board WHERE id = \$1`
		boardID           = "123"
		errA              = errors.New("an error is occurred")
		errB              = errors.New("another error is occurred")
	)

	for _, c := range []struct {
		name      string
		setUpMock func(sqlmock.Sqlmock)
		wantErrs  []error
	}{
		{
			name: "BeginTxErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errA)
			},
			wantErrs: []error{errA},
		},
		{
			name: "SelectColumnIDsErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "SelectColumnIDsRollbackErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "SelectTaskIDsErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(11).AddRow(12).AddRow(13).AddRow(14),
					)
				mock.ExpectQuery(sqlSelectTaskIDs).
					WithArgs(11).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "SelectTaskIDsRollbackErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(11).AddRow(12).AddRow(13).AddRow(14),
					)
				mock.ExpectQuery(sqlSelectTaskIDs).
					WithArgs(11).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "DeleteSubtasksErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(11).AddRow(12).AddRow(13).AddRow(14),
					)
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectQuery(sqlSelectTaskIDs).
						WithArgs(columnID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id"}).
								AddRow(columnID + 10),
						)
				}
				mock.ExpectExec(sqlDeleteSubtasks).
					WithArgs(21).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "DeleteSubtasksRollbackErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(11).AddRow(12).AddRow(13).AddRow(14),
					)
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectQuery(sqlSelectTaskIDs).
						WithArgs(columnID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id"}).
								AddRow(columnID + 10),
						)
				}
				mock.ExpectExec(sqlDeleteSubtasks).
					WithArgs(21).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "DeleteTasksErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(11).AddRow(12).AddRow(13).AddRow(14),
					)
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectQuery(sqlSelectTaskIDs).
						WithArgs(columnID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id"}).
								AddRow(columnID + 10),
						)
				}
				for taskID := 21; taskID < 25; taskID++ {
					mock.ExpectExec(sqlDeleteSubtasks).
						WithArgs(taskID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				mock.ExpectExec(sqlDeleteTasks).
					WithArgs(11).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "DeleteTasksRollbackErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(11).AddRow(12).AddRow(13).AddRow(14),
					)
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectQuery(sqlSelectTaskIDs).
						WithArgs(columnID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id"}).
								AddRow(columnID + 10),
						)
				}
				for taskID := 21; taskID < 25; taskID++ {
					mock.ExpectExec(sqlDeleteSubtasks).
						WithArgs(taskID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				mock.ExpectExec(sqlDeleteTasks).
					WithArgs(11).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "DeleteColumnsErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(11).AddRow(12).AddRow(13).AddRow(14),
					)
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectQuery(sqlSelectTaskIDs).
						WithArgs(columnID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id"}).
								AddRow(columnID + 10),
						)
				}
				for taskID := 21; taskID < 25; taskID++ {
					mock.ExpectExec(sqlDeleteSubtasks).
						WithArgs(taskID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectExec(sqlDeleteTasks).
						WithArgs(columnID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				mock.ExpectExec(sqlDeleteColumns).
					WithArgs(boardID).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "DeleteColumnsRollbackErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(11).AddRow(12).AddRow(13).AddRow(14),
					)
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectQuery(sqlSelectTaskIDs).
						WithArgs(columnID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id"}).
								AddRow(columnID + 10),
						)
				}
				for taskID := 21; taskID < 25; taskID++ {
					mock.ExpectExec(sqlDeleteSubtasks).
						WithArgs(taskID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectExec(sqlDeleteTasks).
						WithArgs(columnID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				mock.ExpectExec(sqlDeleteColumns).
					WithArgs(boardID).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "DeleteBoardErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(11).AddRow(12).AddRow(13).AddRow(14),
					)
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectQuery(sqlSelectTaskIDs).
						WithArgs(columnID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id"}).
								AddRow(columnID + 10),
						)
				}
				for taskID := 21; taskID < 25; taskID++ {
					mock.ExpectExec(sqlDeleteSubtasks).
						WithArgs(taskID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectExec(sqlDeleteTasks).
						WithArgs(columnID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				mock.ExpectExec(sqlDeleteColumns).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(-1, 4))
				mock.ExpectExec(sqlDeleteBoard).
					WithArgs(boardID).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "DeleteBoardRollbackErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(11).AddRow(12).AddRow(13).AddRow(14),
					)
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectQuery(sqlSelectTaskIDs).
						WithArgs(columnID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id"}).
								AddRow(columnID + 10),
						)
				}
				for taskID := 21; taskID < 25; taskID++ {
					mock.ExpectExec(sqlDeleteSubtasks).
						WithArgs(taskID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectExec(sqlDeleteTasks).
						WithArgs(columnID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				mock.ExpectExec(sqlDeleteColumns).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(-1, 4))
				mock.ExpectExec(sqlDeleteBoard).
					WithArgs(boardID).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "CommitErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(11).AddRow(12).AddRow(13).AddRow(14),
					)
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectQuery(sqlSelectTaskIDs).
						WithArgs(columnID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id"}).
								AddRow(columnID + 10),
						)
				}
				for taskID := 21; taskID < 25; taskID++ {
					mock.ExpectExec(sqlDeleteSubtasks).
						WithArgs(taskID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectExec(sqlDeleteTasks).
						WithArgs(columnID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				mock.ExpectExec(sqlDeleteColumns).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(-1, 4))
				mock.ExpectExec(sqlDeleteBoard).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit().WillReturnError(errA)
			},
			wantErrs: []error{errA},
		},
		{
			name: "Success",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectColumnIDs).
					WithArgs(boardID).
					WillReturnRows(
						sqlmock.NewRows([]string{"id"}).
							AddRow(11).AddRow(12).AddRow(13).AddRow(14),
					)
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectQuery(sqlSelectTaskIDs).
						WithArgs(columnID).
						WillReturnRows(
							sqlmock.NewRows([]string{"id"}).
								AddRow(columnID + 10),
						)
				}
				for taskID := 21; taskID < 25; taskID++ {
					mock.ExpectExec(sqlDeleteSubtasks).
						WithArgs(taskID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				for columnID := 11; columnID < 15; columnID++ {
					mock.ExpectExec(sqlDeleteTasks).
						WithArgs(columnID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				mock.ExpectExec(sqlDeleteColumns).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(-1, 4))
				mock.ExpectExec(sqlDeleteBoard).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErrs: []error{nil},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.setUpMock(mock)

			err := sut.Delete(boardID)

			for _, wantErr := range c.wantErrs {
				assert.SameError(t.Error, err, wantErr)
			}
		})
	}
}
