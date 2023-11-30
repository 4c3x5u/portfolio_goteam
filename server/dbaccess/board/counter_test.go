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
		mock.ExpectQuery(cmdCount).
			WithArgs(teamID).
			WillReturnError(sql.ErrNoRows)

		_, err := sut.Count(teamID)

		assert.SameError(t.Error, err, sql.ErrNoRows)
	})

	t.Run("OK", func(t *testing.T) {
		teamID := "2"
		wantCount := 21
		mock.ExpectQuery(cmdCount).WithArgs(teamID).WillReturnRows(
			sqlmock.NewRows([]string{"count"}).AddRow(wantCount),
		)

		count, err := sut.Count(teamID)
		assert.Nil(t.Fatal, err)

		assert.Equal(t.Error, count, wantCount)
	})
}
