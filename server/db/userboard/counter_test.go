package userboard

import (
	"database/sql"
	"fmt"
	"testing"

	"server/assert"
	"server/db"
)

// TestCounter tests the Count method of Counter to assert that it executes the
// correct SQL query with the correct arguments, and returns the count back
// alongside any sql error occurred.
func TestCounter(t *testing.T) {
	db, mock, teardown := db.SetUpDBTest(t)
	defer teardown()
	sut := NewCounter(db)
	username := "bob123"
	query := `SELECT COUNT\(\*\) FROM app.user_board ` +
		`WHERE username = \$1 AND isAdmin = TRUE`

	for wantCount, wantErr := range map[int]error{
		0: sql.ErrConnDone,
		3: sql.ErrNoRows,
		6: nil,
	} {
		t.Run(fmt.Sprintf("Count%d", wantCount), func(t *testing.T) {
			mock.ExpectQuery(query).WithArgs(username).WillReturnRows(
				mock.NewRows([]string{"count"}).AddRow(wantCount),
			).WillReturnError(wantErr)

			count, errCount := sut.Count(username)

			if err := assert.Equal(wantErr, errCount); err != nil {
				t.Error(err)
			}

			if wantErr == nil {
				if err := assert.Equal(wantCount, count); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
