package db

import (
	"fmt"
	"testing"

	"server/assert"
)

func TestBoardUserSelector(t *testing.T) {
	db, mock, teardown := setUpDBTest(t)
	defer teardown()
	sut := NewUserBoardCounter(db)
	userID := "bob123"
	query := `SELECT COUNT\(\*\) FROM app.user_board ` +
		`WHERE userID = \$1 AND isAdmin = \$2`

	for _, wantCount := range []int{0, 1, 2, 3} {
		t.Run(fmt.Sprintf("count: %d", wantCount), func(t *testing.T) {
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
