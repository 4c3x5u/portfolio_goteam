package column

import (
	"database/sql"
	"server/assert"
	"server/dbaccess"
	"testing"
)

// TestUpdater tests the Update method of Updater to assert that it
// sends the correct query to the database with the correct arguments, and
// returns whatever error occurs.
func TestUpdater(t *testing.T) {
	t.Run("Error", func(t *testing.T) {
		var (
			columnID = "1"
			tasks    = []Task{{ID: 1, Order: 1}, {ID: 2, Order: 2}, {ID: 3, Order: 3}}
			wantErr  = sql.ErrNoRows
		)

		db, mock, teardown := dbaccess.SetUpDBTest(t)
		defer teardown()

		sut := NewUpdater(db)

		mock.ExpectExec(
			`UPDATE app.task SET columnID = \$1 AND order = \$2 WHERE id = \$3`,
		).WithArgs(columnID, tasks[0].Order, tasks[0].ID).
			WillReturnError(wantErr)

		err := sut.Update(columnID, tasks)

		if assertErr := assert.SameError(wantErr, err); assertErr != nil {
			t.Error(assertErr)
		}
	})
}
