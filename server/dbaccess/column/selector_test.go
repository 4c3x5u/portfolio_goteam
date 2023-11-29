//go:build utest

package column

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestSelector tests the Select method of Selector to assert that it sends the
// correct query to the database with the correct arguments, and returns all
// data or whatever error occurs.
func TestSelector(t *testing.T) {
	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewSelector(db)

	const (
		sqlSelectBoard = `SELECT id, boardID, "order" ` +
			`FROM app."column" WHERE id = \$1`
		columnID                    = "42"
		existingColumnBoardID       = 21
		existingColumnOrder   int16 = 2
	)

	mock.ExpectQuery(sqlSelectBoard).
		WithArgs(columnID).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery(sqlSelectBoard).
		WithArgs(columnID).
		WillReturnRows(
			sqlmock.NewRows([]string{"id", "boardID", "order"}).
				AddRow(columnID, existingColumnBoardID, existingColumnOrder),
		)

	column, err := sut.Select(columnID)
	if err := assert.SameError(err, sql.ErrNoRows); err != nil {
		t.Error(err)
	}

	column, err = sut.Select(columnID)
	assert.Nil(t.Fatal, err)

	assert.Equal(t.Error, strconv.Itoa(column.ID), columnID)
	assert.Equal(t.Error, column.BoardID, existingColumnBoardID)
	assert.Equal(t.Error, column.Order, existingColumnOrder)
}
