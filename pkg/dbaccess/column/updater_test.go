//go:build utest

package column

import (
	"database/sql"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/dbaccess"
)

// TestUpdater tests the Update method of Updater to assert that it
// sends the correct query to the database with the correct arguments, and
// returns whatever error occurs.
func TestUpdater(t *testing.T) {
	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewUpdater(db)

	var (
		columnID = "1"
		tasks    = []Task{
			{ID: 1, Order: 1},
			{ID: 2, Order: 2},
			{ID: 3, Order: 3},
		}
		sqlUpdateTask = `UPDATE app.task ` +
			`SET columnID = \$1, "order" = \$2 WHERE id = \$3`
		errA = errors.New("an error occurred")
		errB = errors.New("another error occurred")
	)

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
			name: "ExecErr",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(sqlUpdateTask).
					WithArgs(columnID, tasks[0].Order, tasks[0].ID).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "RollbackErr",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(sqlUpdateTask).
					WithArgs(columnID, tasks[0].Order, tasks[0].ID).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "NoRowsAffected",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(sqlUpdateTask).
					WithArgs(columnID, tasks[0].Order, tasks[0].ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback()
			},
			wantErrs: []error{sql.ErrNoRows},
		},
		{
			name: "NoRowsAffectedRollbackErr",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectExec(sqlUpdateTask).
					WithArgs(columnID, tasks[0].Order, tasks[0].ID).
					WillReturnResult(sqlmock.NewResult(0, 0))
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{sql.ErrNoRows, errB},
		},
		{
			name: "CommitErr",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				for _, task := range tasks {
					mock.ExpectExec(sqlUpdateTask).
						WithArgs(columnID, task.Order, task.ID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				mock.ExpectCommit().WillReturnError(errA)
			},
			wantErrs: []error{errA},
		},
		{
			name: "Success",
			mockFunc: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				for _, task := range tasks {
					mock.ExpectExec(sqlUpdateTask).
						WithArgs(columnID, task.Order, task.ID).
						WillReturnResult(sqlmock.NewResult(0, 1))
				}
				mock.ExpectCommit()
			},
			wantErrs: []error{nil},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.mockFunc(mock)
			err := sut.Update(columnID, tasks)
			for _, wantErr := range c.wantErrs {
				assert.ErrIs(t.Error, err, wantErr)
			}
		})
	}
}
