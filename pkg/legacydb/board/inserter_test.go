//go:build utest

package board

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/legacydb"
)

// TestInserter tests the Insert method of Inserter to assert that it sends the
// correct queries to the database with the correct arguments, and returns
// whatever error occurs.
func TestInserter(t *testing.T) {
	db, mock, teardown := legacydb.SetUpDBTest(t)
	defer teardown()

	sut := NewInserter(db)

	var (
		sqlInsertBoard = `INSERT INTO app.board\(name, teamID\) ` +
			`VALUES \(\$1, \$2\)`
		sqlInsertColumn = `INSERT INTO app.\"column\"\(boardID, "order\"\) ` +
			`VALUES \(\$1, \$2\)`
		errA  = errors.New("an error occurred")
		errB  = errors.New("another error occurred")
		board = NewInRecord("someboard", 21)
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
			name: "InsertBoardErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlInsertBoard).
					WithArgs(board.name, board.teamID).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "InsertBoardRollbackErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlInsertBoard).
					WithArgs(board.name, board.teamID).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "InsertColumnErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlInsertBoard).
					WithArgs(board.name, board.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 1).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "InsertColumnRollbackErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlInsertBoard).
					WithArgs(board.name, board.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 1).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "CommitErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlInsertBoard).
					WithArgs(board.name, board.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 3).
					WillReturnResult(sqlmock.NewResult(2, 1))
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 4).
					WillReturnResult(sqlmock.NewResult(3, 1))
				mock.ExpectCommit().WillReturnError(errA)
			},
			wantErrs: []error{errA},
		},
		{
			name: "Success",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlInsertBoard).
					WithArgs(board.name, board.teamID).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 3).
					WillReturnResult(sqlmock.NewResult(2, 1))
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 4).
					WillReturnResult(sqlmock.NewResult(3, 1))
				mock.ExpectCommit()
			},
			wantErrs: []error{nil},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.setUpMock(mock)

			err := sut.Insert(board)

			for _, wantErr := range c.wantErrs {
				assert.ErrIs(t.Error, err, wantErr)
			}
		})
	}
}
