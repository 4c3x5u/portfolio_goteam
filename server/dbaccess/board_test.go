//go:build utest

package dbaccess

import (
	"database/sql"
	"errors"
	"strconv"
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
	board := NewInBoard("someboard", "bob123")

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

			if err = assert.SameError(c.wantErr, err); err != nil {
				t.Error(err)
			}
		})
	}
}

// TestBoardUpdater tests the Update method of BoardUpdater to assert that it
// sends the correct query to the database with the correct arguments, and
// returns whatever error occurs.
func TestBoardUpdater(t *testing.T) {
	db, mock, teardown := setUpDBTest(t)
	defer teardown()

	sut := NewBoardUpdater(db)

	const (
		sqlUpdateBoard = "UPDATE app.board SET name = \\$1 WHERE id = \\$2"
		boardID        = "21"
		newBoardName   = "Board B"
	)

	for _, c := range []struct {
		name        string
		setUp       func(sqlmock.Sqlmock)
		assertOnErr func(error) error
	}{
		{
			name: "SqlErrNoRows",
			setUp: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(sqlUpdateBoard).
					WithArgs(newBoardName, boardID).
					WillReturnError(sql.ErrNoRows)
			},
			assertOnErr: func(err error) error {
				return assert.SameError(err, sql.ErrNoRows)
			},
		},
		{
			name: "NoRowsAffected",
			setUp: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(sqlUpdateBoard).
					WithArgs(newBoardName, boardID).
					WillReturnResult(sqlmock.NewResult(-1, 0))
			},
			assertOnErr: func(err error) error {
				return assert.Equal(err.Error(), "no rows were affected")
			},
		},
		{
			name: "MoreRowsAffected",
			setUp: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(sqlUpdateBoard).
					WithArgs(newBoardName, boardID).
					WillReturnResult(sqlmock.NewResult(-1, 2))
			},
			assertOnErr: func(err error) error {
				return assert.Equal(
					err.Error(), "more than expected rows were affected",
				)
			},
		},
		{
			name: "Success",
			setUp: func(mock sqlmock.Sqlmock) {
				mock.
					ExpectExec(sqlUpdateBoard).
					WithArgs(newBoardName, boardID).
					WillReturnResult(sqlmock.NewResult(21, 1))
			},
			assertOnErr: func(err error) error { return assert.Nil(err) },
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.setUp(mock)
			err := sut.Update(boardID, newBoardName)
			if assertErr := c.assertOnErr(err); assertErr != nil {
				t.Error(assertErr)
			}
		})
	}
}

// TestBoardSelector tests the Select method of BoardSelector to assert that it
// sends the correct query to the database with the correct arguments, and
// returns whatever error occurs.
func TestBoardSelector(t *testing.T) {
	db, mock, teardown := setUpDBTest(t)
	defer teardown()

	sut := NewBoardSelector(db)

	const (
		sqlSelectBoard    = "SELECT id, name FROM app.board WHERE id = \\$1"
		boardID           = "21"
		existingBoardName = "Board A"
	)

	mock.
		ExpectQuery(sqlSelectBoard).
		WithArgs(boardID).
		WillReturnError(sql.ErrNoRows)

	mock.
		ExpectQuery(sqlSelectBoard).
		WithArgs(boardID).
		WillReturnRows(
			sqlmock.
				NewRows([]string{"id", "name"}).
				AddRow(boardID, existingBoardName),
		)

	board, err := sut.Select(boardID)
	if err := assert.SameError(err, sql.ErrNoRows); err != nil {
		t.Error(err)
	}

	board, err = sut.Select(boardID)
	if err = assert.Nil(err); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(boardID, strconv.Itoa(board.id)); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(existingBoardName, board.name); err != nil {
		t.Error(err)
	}
}

// TestBoardDeleter tests the Delete method of BoardDeleter to assert that it
// sends the correct query to the database with the correct arguments, and
// returns whatever error occurs.
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
