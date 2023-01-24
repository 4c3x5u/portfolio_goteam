package db

import (
	"fmt"
	"testing"

	"server/assert"
)

func TestUserBoardSelector(t *testing.T) {
	db, mock, teardown := setUpDBTest(t)
	defer teardown()
	sut := NewUserBoardSelector(db)
	userID := "bob123"
	boardID := "123"
	query := `SELECT isAdmin FROM app.user_board` +
		` WHERE userID = \$1 AND boardID = \$2`

	for _, wantIsAdmin := range []bool{true, false} {
		t.Run(fmt.Sprintf("IsAdmin%v", wantIsAdmin), func(t *testing.T) {
			mock.ExpectQuery(query).WithArgs(userID, boardID).WillReturnRows(
				mock.NewRows([]string{"isAdmin"}).AddRow(wantIsAdmin),
			)

			isAdmin := sut.Select(userID, boardID)

			if err := assert.Equal(wantIsAdmin, isAdmin); err != nil {
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
