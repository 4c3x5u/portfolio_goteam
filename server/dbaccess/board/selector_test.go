//go:build utest

package board

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
)

// TestSelector tests the Select method of Selector to assert that it sends the
// correct query to the database with the correct arguments, and returns all
// data or whatever error occurs.
func TestSelector(t *testing.T) {
	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewSelector(db)

	const (
		sqlSelectBoard = `SELECT id, name, teamID FROM app.board ` +
			`WHERE id = \$1`
		boardID           = "21"
		existingBoardName = "Board A"
		teamID            = 21
	)

	mock.ExpectQuery(sqlSelectBoard).
		WithArgs(boardID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery(sqlSelectBoard).
		WithArgs(boardID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "name", "teamID"}).
				AddRow(boardID, existingBoardName, teamID),
		)

	board, err := sut.Select(boardID)
	if err := assert.SameError(sql.ErrNoRows, err); err != nil {
		t.Error(err)
	}

	board, err = sut.Select(boardID)
	if err = assert.Nil(err); err != nil {
		t.Error(err)
	}
	assert.Equal(t.Error, strconv.Itoa(board.ID), boardID)
	assert.Equal(t.Error, board.Name, existingBoardName)
	assert.Equal(t.Error, board.TeamID, teamID)
}
