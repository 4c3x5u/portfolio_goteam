//go:build utest

package board

import (
	"errors"
	"testing"

	"server/assert"
	"server/dbaccess"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestInserter tests the Insert method of Inserter to assert that it sends the
// correct queries to the database with the correct arguments, and returns
// whatever error occurs.
func TestInserter(t *testing.T) {
	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewInserter(db)

	var (
		sqlInsertBoard     = `INSERT INTO app.board\(name\) VALUES \(\$1\)`
		sqlInsertUserBoard = `INSERT INTO app.user_board\(username, boardID, ` +
			`isAdmin\) VALUES\(\$1, \$2, TRUE\)`
		sqlInsertColumn = `INSERT INTO app.\"column\"\(boardID, "order\"\) ` +
			`VALUES \(\$1, \$2\)`
		errA  = errors.New("an error occurred")
		errB  = errors.New("another error occurred")
		board = NewBoard("someboard", "bob123")
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
					WithArgs(board.name).
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
					WithArgs(board.name).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		{
			name: "InsertUserBoardErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlInsertBoard).
					WithArgs(board.name).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, 1).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "InsertUserBoardRollbackErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlInsertBoard).
					WithArgs(board.name).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, 1).
					WillReturnError(errA)
				mock.ExpectRollback().WillReturnError(errB)
			},
			wantErrs: []error{errA, errB},
		},
		// 4 columns in total are inserted after the board creation. Below are
		// error cases for each of these insert operations.
		{
			name: "InsertColumnErr#1",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlInsertBoard).
					WithArgs(board.name).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 1).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "InsertColumnErr#2",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlInsertBoard).
					WithArgs(board.name).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 2).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "InsertColumnErr#3",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlInsertBoard).
					WithArgs(board.name).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 2).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectExec(sqlInsertColumn).
					WithArgs(1, 3).
					WillReturnError(errA)
				mock.ExpectRollback()
			},
			wantErrs: []error{errA},
		},
		{
			name: "InsertColumnErr#4",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlInsertBoard).
					WithArgs(board.name).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
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
					WithArgs(board.name).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
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
					WithArgs(board.name).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
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
					WithArgs(board.name).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
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
				if assertErr := assert.SameError(
					wantErr, err,
				); assertErr != nil {
					t.Error(assertErr)
				}
			}
		})
	}
}
