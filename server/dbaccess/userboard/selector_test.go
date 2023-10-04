//go:build utest

package userboard

import (
	"errors"
	"testing"

	"server/assert"
	"server/dbaccess"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestSelector tests the Select method of Selector to assert that it
// sends the correct query to the database with the correct arguments, and
// returns whatever error occurs.
func TestSelector(t *testing.T) {
	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()
	sut := NewSelector(db)
	username := "bob123"
	boardID := "123"
	query := `SELECT isAdmin FROM app.user_board` +
		` WHERE username = \$1 AND boardID = \$2`
	sqlErr := errors.New("postgres query error")

	for _, c := range []struct {
		name        string
		wantIsAdmin bool
		wantErr     error
		setUpMock   func(mock *sqlmock.ExpectedQuery)
	}{
		{
			name:        "SqlErr",
			wantIsAdmin: false,
			wantErr:     sqlErr,
			setUpMock: func(mock *sqlmock.ExpectedQuery) {
				mock.WillReturnError(sqlErr)
			},
		},
		{
			name:        "IsNotAdmin",
			wantIsAdmin: false,
			wantErr:     nil,
			setUpMock: func(mock *sqlmock.ExpectedQuery) {
				mock.WillReturnRows(
					sqlmock.NewRows([]string{"isAdmin"}).AddRow(false),
				)
			},
		},
		{
			name:        "IsAdmin",
			wantIsAdmin: true,
			wantErr:     nil,
			setUpMock: func(mock *sqlmock.ExpectedQuery) {
				mock.WillReturnRows(
					sqlmock.NewRows([]string{"isAdmin"}).AddRow(true),
				).WillReturnError(nil)
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			c.setUpMock(mock.ExpectQuery(query).WithArgs(username, boardID))

			isAdmin, err := sut.Select(username, boardID)

			if err = assert.Equal(c.wantErr, err); err != nil {
				t.Error(err)
			}
			if err := assert.Equal(c.wantIsAdmin, isAdmin); err != nil {
				t.Error(err)
			}
		})
	}
}
