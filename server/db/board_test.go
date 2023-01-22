package db

import (
	"errors"
	"testing"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestBoardInserter(t *testing.T) {
	sqlInsertBoard := `INSERT INTO app.board\(name\) VALUES \(\$1\)`
	sqlInsertUserBoard := `INSERT INTO app.user_board\(userID, boardID, ` +
		`isAdmin\) VALUES\(\$1, \$2, \$3\)`

	// Since we're doing pointer comparison for errors in test cases below, we
	// just need a generic error to be returned at different points in code and
	// ansure the expected SQL is executed and the same error is returned.
	wantErr := errors.New("error occured")
	board := NewBoard("bob21", "someboard")

	for _, c := range []struct {
		name      string
		setUpMock func(sqlmock.Sqlmock)
		wantErr   error
	}{
		{
			name: "BeginTxErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(wantErr)
			},
			wantErr: wantErr,
		},
		{
			name: "RollbackInsertBoardErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectExec(sqlInsertBoard).
					WithArgs(board.name).
					WillReturnError(wantErr)
				mock.ExpectRollback()
			},
			wantErr: wantErr,
		},
		{
			name: "RollbackInsertUserBoardErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectExec(sqlInsertBoard).
					WithArgs(board.name).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, board.name, true).
					WillReturnError(wantErr)
				mock.ExpectRollback()
			},
			wantErr: wantErr,
		},
		{
			name: "Success",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.
					ExpectExec(sqlInsertBoard).
					WithArgs(board.name).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.
					ExpectExec(sqlInsertUserBoard).
					WithArgs(board.adminID, board.name, true).
					WillReturnResult(sqlmock.NewResult(0, 1))
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

			if err := assert.Equal(c.wantErr, err); err != nil {
				t.Error(err)
			}
		})
	}
}
