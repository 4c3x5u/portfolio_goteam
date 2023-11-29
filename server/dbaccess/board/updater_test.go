//go:build utest

package board

import (
	"database/sql"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"

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
		assertOnErr func(func(...any), error)
	}{
		{
			name: "SqlErrNoRows",
			setUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(sqlUpdateBoard).
					WithArgs(newBoardName, boardID).
					WillReturnError(sql.ErrNoRows)
			},
			assertOnErr: func(errFunc func(...any), err error) {
				if err = assert.SameError(err, sql.ErrNoRows); err != nil {
					errFunc(err)
				}
			},
		},
		{
			name: "NoRowsAffected",
			setUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(sqlUpdateBoard).
					WithArgs(newBoardName, boardID).
					WillReturnResult(sqlmock.NewResult(-1, 0))
			},
			assertOnErr: func(logFunc func(...any), err error) {
				assert.Equal(logFunc, err.Error(), "no rows were affected")
			},
		},
		{
			name: "MoreRowsAffected",
			setUp: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(sqlUpdateBoard).
					WithArgs(newBoardName, boardID).
					WillReturnResult(sqlmock.NewResult(-1, 2))
			},
			assertOnErr: func(logFunc func(...any), err error) {
				assert.Equal(logFunc,
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
			assertOnErr: func(logFunc func(...any), err error) {
				if err := assert.Nil(err); err != nil {
					logFunc(err)
				}
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
