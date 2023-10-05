//go:build utest

package board

import (
	"database/sql"
	"testing"

	"server/assert"
	"server/dbaccess"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestUpdater tests the Update method of Updater to assert that it
// sends the correct query to the database with the correct arguments, and
// returns whatever error occurs.
func TestUpdater(t *testing.T) {
	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewUpdater(db)

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
				mock.ExpectExec(sqlUpdateBoard).
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
				mock.ExpectExec(sqlUpdateBoard).
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
				mock.ExpectExec(sqlUpdateBoard).
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
				mock.ExpectExec(sqlUpdateBoard).
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
