//go:build utest

package subtask

import (
	"database/sql"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/dbaccess"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestUpdater tests the Update method of Updater to assert that it sends the
// correct queries to the database with the correct arguments, and returns
// whatever error occurs.
func TestUpdater(t *testing.T) {
	const (
		sqlUpdateSubtask = `UPDATE app.subtask SET isDone = \$1 WHERE id = \$2`
		id               = "1"
	)
	rec := NewUpRecord(true)

	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()
	sut := NewUpdater(db)

	t.Run("Error", func(t *testing.T) {
		wantErr := sql.ErrNoRows

		mock.ExpectExec(sqlUpdateSubtask).
			WithArgs(rec.isDone, id).
			WillReturnError(wantErr)

		err := sut.Update(id, rec)

		assert.ErrIs(t.Error, err, wantErr)
	})

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(sqlUpdateSubtask).
			WithArgs(rec.isDone, id).
			WillReturnResult(sqlmock.NewResult(0, 0))

		err := sut.Update(id, rec)

		assert.Nil(t.Error, err)
	})
}
