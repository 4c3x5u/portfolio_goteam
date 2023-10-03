//go:build utest

package board

import (
	"errors"
	"server/db"
	"testing"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestDeleter tests the Delete method of Deleter to assert that it sends the
// correct query to the database with the correct arguments, and returns
// whatever error occurs.
func TestDeleter(t *testing.T) {
	const (
		sqlDeleteUserBoards = `DELETE FROM app.user_board WHERE boardID = \$1`
		sqlDeleteColumns    = `DELETE FROM app.\"column\" WHERE boardID = \$1`
		sqlDeleteBoard      = `DELETE FROM app.board WHERE id = \$1`
		boardID             = "123"
	)

	errA := errors.New("an error is occured")
	errB := errors.New("another error is occured")

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
			name: "DeleteUserBoardRollbackErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectExec(sqlDeleteUserBoards).
					WithArgs(boardID).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "DeleteUserBoardErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectExec(sqlDeleteUserBoards).
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
				mock.
					ExpectExec(sqlDeleteUserBoards).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectExec(sqlDeleteColumns).
					WithArgs(boardID).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "DeleteColumnsErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectExec(sqlDeleteUserBoards).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectExec(sqlDeleteColumns).
					WithArgs(boardID).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "DeleteBoardErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectExec(sqlDeleteUserBoards).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectExec(sqlDeleteColumns).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(-1, 4))
				mock.
					ExpectExec(sqlDeleteBoard).
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
				mock.
					ExpectExec(sqlDeleteUserBoards).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectExec(sqlDeleteColumns).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(-1, 4))
				mock.
					ExpectExec(sqlDeleteBoard).
					WithArgs(boardID).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "CommitErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectExec(sqlDeleteUserBoards).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectExec(sqlDeleteColumns).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(-1, 4))
				mock.
					ExpectExec(sqlDeleteBoard).
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
				mock.
					ExpectExec(sqlDeleteUserBoards).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectExec(sqlDeleteColumns).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(-1, 4))
				mock.
					ExpectExec(sqlDeleteBoard).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErrs: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			db, mock, teardown := db.SetUpDBTest(t)
			defer teardown()
			c.setUpMock(mock)
			sut := NewDeleter(db)

			err := sut.Delete(boardID)

			for _, wantErr := range c.wantErrs {
				if assertErr := assert.SameError(
					wantErr, err,
				); assertErr != nil {
					t.Error(assertErr)
				}
			}
		})
	}
}
