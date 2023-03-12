//go:build utest

package db

import (
	"errors"
	"testing"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestBoardInserter tests the Insert method of BoardInserter to assert that it
// sends the correct queries to the database with the correct arguments, and
// returns whatever error occurs.
func TestBoardInserter(t *testing.T) {
	const (
		sqlInsertBoard     = `INSERT INTO app.board\(name\) VALUES \(\$1\)`
		sqlInsertUserBoard = `INSERT INTO app.user_board\(username, boardID, ` +
			`isAdmin\) VALUES\(\$1, \$2, TRUE\)`
	)

	// Since we're doing pointer comparison for errors in test cases below, we
	// just need a generic error to be returned at different points in code and
	// ansure the expected SQL is executed and the same error is returned.
	someErr := errors.New("some error occured")
	board := NewBoard("someboard", "bob123")

	for _, c := range []struct {
		name      string
		setUpMock func(sqlmock.Sqlmock)
		wantErr   error
	}{
		{
			name: "BeginTxErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(someErr)
			},
			wantErr: someErr,
		},
		{
			name: "InsertBoardErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectQuery(sqlInsertBoard).
					WithArgs(board.name).
					WillReturnError(someErr)
				mock.ExpectRollback()
			},
			wantErr: someErr,
		},
		{
			name: "InsertUserBoardErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectQuery(sqlInsertBoard).
					WithArgs(board.name).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.
					ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, 1).
					WillReturnError(someErr)
				mock.ExpectRollback()
			},
			wantErr: someErr,
		},
		{
			name: "Success",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectQuery(sqlInsertBoard).
					WithArgs(board.name).
					WillReturnRows(
						sqlmock.NewRows([]string{"boardID"}).AddRow(1),
					)
				mock.
					ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, 1).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			db, mock, teardown := setUpDBTest(t)
			defer teardown()
			c.setUpMock(mock)
			sut := NewBoardInserter(db)

			err := sut.Insert(board)

			if err = assert.Equal(c.wantErr, err); err != nil {
				t.Error(err)
			}
		})
	}
}

func TestBoardDeleter(t *testing.T) {
	const (
		sqlDeleteRel   = `DELETE FROM app.user_board WHERE boardID = \$1`
		sqlDeleteBoard = `DELETE FROM app.board WHERE id = \$1`
		boardID        = "123"
	)

	someErr := errors.New("some error occured")

	for _, c := range []struct {
		name      string
		setUpMock func(sqlmock.Sqlmock)
		wantErr   error
	}{
		{
			name: "BeginTxErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(someErr)
			},
			wantErr: someErr,
		},
		{
			name: "DeleteRelsErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectExec(sqlDeleteRel).
					WithArgs(boardID).
					WillReturnError(someErr)
				mock.ExpectRollback()
			},
			wantErr: someErr,
		},
		{
			name: "DeleteBoardErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectExec(sqlDeleteRel).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectExec(sqlDeleteBoard).
					WithArgs(boardID).
					WillReturnError(someErr)
				mock.ExpectRollback()
			},
			wantErr: someErr,
		},
		{
			name: "Success",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectExec(sqlDeleteRel).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectExec(sqlDeleteBoard).
					WithArgs(boardID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			},
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			db, mock, teardown := setUpDBTest(t)
			defer teardown()
			c.setUpMock(mock)
			sut := NewBoardDeleter(db)

			err := sut.Delete(boardID)
			if err = assert.Equal(c.wantErr, err); err != nil {
				t.Error(err)
			}
		})
	}
}

// TestWrapRollbackErr tests the wrapRollbackErr helper function to assert
// that it constructs a sensible error string when called with two different
// errors.
func TestWrapRollbackErr(t *testing.T) {
	err := errors.New("something went wrong")
	rollbackErr := errors.New("rollback error")
	wantErrStr := "multiple errors occured:" +
		"\n  (0) err: " + err.Error() +
		"\n  (1) rollbackErr: " + rollbackErr.Error()

	got := wrapRollbackErr(err, rollbackErr)

	if err := assert.Equal(wantErrStr, got.Error()); err != nil {
		t.Error(err)
	}
}
