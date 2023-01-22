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

	board := NewBoard("bob21", "someboard")

	t.Run("BeginTxErr", func(t *testing.T) {
		db, mock, teardown := setUpDBMock(t)
		defer teardown()
		sut := NewBoardInserter(db)

		wantErr := errors.New("tx begin error")
		mock.ExpectBegin().WillReturnError(wantErr)

		err := sut.Insert(board)

		if err := assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("RollbackInsertBoardErr", func(t *testing.T) {
		db, mock, teardown := setUpDBMock(t)
		defer teardown()
		sut := NewBoardInserter(db)

		wantErr := errors.New("insert board error")
		mock.ExpectBegin()
		mock.
			ExpectExec(sqlInsertBoard).
			WithArgs(board.name).
			WillReturnError(wantErr)
		mock.ExpectRollback()

		err := sut.Insert(board)

		if err := assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("RollbackInsertUserBoardErr", func(t *testing.T) {
		db, mock, teardown := setUpDBMock(t)
		defer teardown()
		sut := NewBoardInserter(db)

		wantErr := errors.New("insert userBoard error")
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

		err := sut.Insert(board)

		if err := assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		db, mock, teardown := setUpDBMock(t)
		defer teardown()
		sut := NewBoardInserter(db)

		mock.ExpectBegin()
		mock.
			ExpectExec(sqlInsertBoard).
			WithArgs(board.name).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.
			ExpectExec(sqlInsertUserBoard).
			WithArgs(board.adminID, board.name, true).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := sut.Insert(board)

		if err != nil {
			t.Error(err)
		}
	})
}
