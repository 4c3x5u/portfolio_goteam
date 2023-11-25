//go:build utest

package task

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
)

// TestDeleter tests the Delete method of Deleter to assert that it sends the
// correct query to the database with the correct arguments, and returns
// whatever error occurs.
func TestDeleter(t *testing.T) {
	sqlDeleteTask := `DELETE FROM app.task WHERE id = \$1`
	sqlDeleteSubtask := `DELETE FROM app.subtask WHERE taskID = \$1`
	id := "3"
	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()
	sut := NewDeleter(db)

	t.Run("BeginErr", func(t *testing.T) {
		wantErr := sql.ErrConnDone
		mock.ExpectBegin().WillReturnError(wantErr)

		err := sut.Delete(id)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("DeleteSubtaskErr", func(t *testing.T) {
		wantErr := sql.ErrConnDone
		mock.ExpectBegin()
		mock.ExpectExec(sqlDeleteSubtask).
			WithArgs(id).
			WillReturnError(wantErr)

		err := sut.Delete(id)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("DeleteTaskErr", func(t *testing.T) {
		wantErr := sql.ErrNoRows
		mock.ExpectBegin()
		mock.ExpectExec(sqlDeleteSubtask).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(sqlDeleteTask).
			WithArgs(id).
			WillReturnError(wantErr)
		mock.ExpectRollback()

		err := sut.Delete(id)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("RollbackErr", func(t *testing.T) {
		wantErr := sql.ErrNoRows
		wantRollbackErr := sql.ErrConnDone
		mock.ExpectBegin()
		mock.ExpectExec(sqlDeleteSubtask).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(sqlDeleteTask).
			WithArgs(id).
			WillReturnError(wantErr)
		mock.ExpectRollback().WillReturnError(wantRollbackErr)

		err := sut.Delete(id)

		if assertErr := assert.SameError(wantErr, err); assertErr != nil {
			t.Error(assertErr)
		}
		if assertErr := assert.SameError(
			wantRollbackErr, err,
		); assertErr != nil {
			t.Error(assertErr)
		}
	})

	t.Run("CommitErr", func(t *testing.T) {
		wantErr := sql.ErrConnDone
		mock.ExpectBegin()
		mock.ExpectExec(sqlDeleteSubtask).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(sqlDeleteTask).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit().WillReturnError(wantErr)

		err := sut.Delete(id)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(sqlDeleteSubtask).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 0))
		mock.ExpectExec(sqlDeleteTask).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := sut.Delete(id)

		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
	})
}
