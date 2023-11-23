//go:build utest

package board

import (
	"database/sql"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestCounter tests the Counter.Count to assert that it sends the correct
func TestCounter(t *testing.T) {
	cmdCount := `SELECT COUNT\(\*\) FROM app.board WHERE teamID = \$1`

	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewCounter(db)

	t.Run("Error", func(t *testing.T) {
		teamID := "2"
		mock.ExpectQuery(cmdCount).WithArgs(teamID).
			WillReturnError(sql.ErrNoRows)

		_, err := sut.Count(teamID)

		if err = assert.SameError(sql.ErrNoRows, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("OK", func(t *testing.T) {
		teamID := "2"
		wantCount := 21
		mock.ExpectQuery(cmdCount).WithArgs(teamID).WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(wantCount),
		)

		count, err := sut.Count(teamID)
		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(wantCount, count); err != nil {
			t.Error(err)
		}
	})
}
