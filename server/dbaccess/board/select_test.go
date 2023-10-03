package board

import (
	"database/sql"
	"strconv"
	"testing"

	"server/assert"
	"server/dbaccess"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestSelector tests the Select method of Selector to assert that it sends the
// correct query to the database with the correct arguments, and returns
// whatever error occurs.
func TestSelector(t *testing.T) {
	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewSelector(db)

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
