//go:build utest

package subtask

import (
	"database/sql"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestSelector tests the Select method of Selector to assert that it sends the
// correct query to the database with the correct arguments, and returns all
// data or whatever error occurs.
func TestSelector(t *testing.T) {
	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewSelector(db)

	const (
		sqlSelectTask = `SELECT taskID FROM app.subtask WHERE id = \$1`
		id            = "1"
	)

	t.Run("Error", func(t *testing.T) {
		wantErr := sql.ErrNoRows

		mock.ExpectQuery(sqlSelectTask).WithArgs(id).WillReturnError(wantErr)

		_, err := sut.Select(id)
		if err := assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		wantColumnID := 2

		mock.ExpectQuery(sqlSelectTask).
			WithArgs(id).
			WillReturnRows(
				sqlmock.NewRows([]string{"columnID"}).AddRow(wantColumnID),
			)

		subtask, err := sut.Select(id)
		assert.Nil(t.Error, err)
		assert.Equal(t.Error, subtask.TaskID, wantColumnID)
	})
}
