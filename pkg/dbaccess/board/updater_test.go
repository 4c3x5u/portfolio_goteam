//go:build utest

package board

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/dbaccess"
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
		assertOnErr func(func(...any), error)
	}{
		{
			name: "SqlErrNoRows",
			setUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(sqlUpdateBoard).
					WithArgs(newBoardName, boardID).
					WillReturnError(sql.ErrNoRows)
			},
			assertOnErr: func(logErr func(...any), err error) {
				assert.ErrIs(logErr, err, sql.ErrNoRows)
			},
		},
		{
			name: "NoRowsAffected",
			setUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(sqlUpdateBoard).
					WithArgs(newBoardName, boardID).
					WillReturnResult(sqlmock.NewResult(-1, 0))
			},
			assertOnErr: func(logErr func(...any), err error) {
				assert.Equal(logErr, err.Error(), "no rows were affected")
			},
		},
		{
			name: "MoreRowsAffected",
			setUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(sqlUpdateBoard).
					WithArgs(newBoardName, boardID).
					WillReturnResult(sqlmock.NewResult(-1, 2))
			},
			assertOnErr: func(logErr func(...any), err error) {
				assert.Equal(logErr,
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
			assertOnErr: func(logErr func(...any), err error) {
				assert.Nil(t.Error, err)
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.setUp(mock)
			err := sut.Update(boardID, newBoardName)
			c.assertOnErr(t.Error, err)
		})
	}
}
