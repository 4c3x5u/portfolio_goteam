//go:build utest

package task

import (
	"errors"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestUpdater tests the Update method of Updater to assert that it sends the
// correct queries to the database with the correct arguments, and returns
// whatever error occurs.
func TestUpdater(t *testing.T) {
	const sqlUpdateTask = `UPDATE app.task ` +
		`SET title = \$1, description = \$2 WHERE id = \$3`
	const sqlDeleteSubtasks = `DELETE FROM app.subtask WHERE taskID = \$1`
	const sqlInsertSubtask = `INSERT ` +
		`INTO app.subtask\(taskID, title, \"order\", isDone\)` +
		`VALUES\(\$1, \$2, \$3, \$4\)`
	id := "1"
	rec := NewUpRecord(
		"Some Task",
		"Do Something",
		[]Subtask{
			NewSubtask("Some Subtask", 0, false),
			NewSubtask("Some Other Subtask", 1, true),
		},
	)
	errA := errors.New("an error occurred")
	errB := errors.New("another error occurred")

	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()
	sut := NewUpdater(db)

	for _, c := range []struct {
		name     string
		mockFunc func(sqlmock.Sqlmock)
		wantErrs []error
	}{
		{
			name: "BeginErr",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(errA)
			},
			wantErrs: []error{errA},
		},
		{
			name: "UpdateTaskErr",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(sqlUpdateTask).
					WithArgs(rec.Title, rec.Description, id).
					WillReturnError(errA)
			},
			wantErrs: []error{errA},
		},
		{
			name: "DeleteSubtasksErr",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(sqlUpdateTask).
					WithArgs(rec.Title, rec.Description, id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(sqlDeleteSubtasks).
					WithArgs(id).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "DeleteSubtasksRollbackErr",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(sqlUpdateTask).
					WithArgs(rec.Title, rec.Description, id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(sqlDeleteSubtasks).
					WithArgs(id).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "InsertSubtaskErr",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(sqlUpdateTask).
					WithArgs(rec.Title, rec.Description, id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(sqlDeleteSubtasks).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(sqlInsertSubtask).WithArgs(
					id,
					rec.Subtasks[0].Title,
					rec.Subtasks[0].Order,
					rec.Subtasks[0].IsDone,
				).WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "InsertSubtaskRollbackErr",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(sqlUpdateTask).
					WithArgs(rec.Title, rec.Description, id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(sqlDeleteSubtasks).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectExec(sqlInsertSubtask).WithArgs(
					id,
					rec.Subtasks[0].Title,
					rec.Subtasks[0].Order,
					rec.Subtasks[0].IsDone,
				).WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "CommitErr",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(sqlUpdateTask).
					WithArgs(rec.Title, rec.Description, id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(sqlDeleteSubtasks).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 0))
				for i, subtask := range rec.Subtasks {
					mock.ExpectExec(sqlInsertSubtask).WithArgs(
						id, subtask.Title, subtask.Order, subtask.IsDone,
					).WillReturnResult(sqlmock.NewResult(int64(i), 1))
				}
				mock.ExpectCommit().WillReturnError(errA)
			},
			wantErrs: []error{errA},
		},
		{
			name: "Success",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(sqlUpdateTask).
					WithArgs(rec.Title, rec.Description, id).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(sqlDeleteSubtasks).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 0))
				for i, subtask := range rec.Subtasks {
					mock.ExpectExec(sqlInsertSubtask).WithArgs(
						id, subtask.Title, subtask.Order, subtask.IsDone,
					).WillReturnResult(sqlmock.NewResult(int64(i), 1))
				}
				mock.ExpectCommit()
			},
			wantErrs: []error{nil},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc(mock)
			err := sut.Update(id, rec)
			for _, wantErr := range c.wantErrs {
				if assertErr := assert.SameError(wantErr, err); assertErr != nil {
					t.Error(assertErr)
				}
			}
		})
	}
}
