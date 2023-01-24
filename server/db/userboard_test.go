package db

import (
	"errors"
	"fmt"
	"testing"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUserBoardSelector(t *testing.T) {
	db, mock, teardown := setUpDBTest(t)
	defer teardown()
	sut := NewUserBoardSelector(db)
	userID := "bob123"
	boardID := "123"
	query := `SELECT isAdmin FROM app.user_board` +
		` WHERE userID = \$1 AND boardID = \$2`
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
			c.setUpMock(mock.ExpectQuery(query).WithArgs(userID, boardID))

			isAdmin, err := sut.Select(userID, boardID)

			if err := assert.Equal(c.wantErr, err); err != nil {
				t.Error(err)
			}
			if err := assert.Equal(c.wantIsAdmin, isAdmin); err != nil {
				t.Error(err)
			}
		})
	}
}

// TestUserBoardCounter tests the Count method of UserBoardCounter to assert
// that it executes the correct SQL query with the correct arguments, and
// returns the count back.
func TestUserBoardCounter(t *testing.T) {
	db, mock, teardown := setUpDBTest(t)
	defer teardown()
	sut := NewUserBoardCounter(db)
	userID := "bob123"
	query := `SELECT COUNT\(\*\) FROM app.user_board ` +
		`WHERE userID = \$1 AND isAdmin = \$2`

	for _, wantCount := range []int{0, 1, 2, 3} {
		t.Run(fmt.Sprintf("Count%d", wantCount), func(t *testing.T) {
			mock.ExpectQuery(query).WithArgs(userID, true).WillReturnRows(
				mock.NewRows([]string{"count"}).AddRow(wantCount),
			)

			gotCount := sut.Count(userID)

			if err := assert.Equal(wantCount, gotCount); err != nil {
				t.Error(err)
			}
		})
	}
}
