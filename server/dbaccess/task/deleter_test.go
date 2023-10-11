//go:build utest

package task

import (
	"database/sql"
	"testing"

	"server/assert"
	"server/dbaccess"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestDeleter tests the Delete method of Deleter to assert that it sends the
// correct query to the database with the correct arguments, and returns
// whatever error occurs.
func TestDeleter(t *testing.T) {
	sqlDeleteTask := `DELETE FROM app.task WHERE id = \$1`
	id := "3"
	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()
	sut := NewDeleter(db)

	t.Run("Err", func(t *testing.T) {
		mock.ExpectExec(sqlDeleteTask).
			WithArgs(id).
			WillReturnError(sql.ErrNoRows)

		err := sut.Delete(id)

		if err = assert.SameError(sql.ErrNoRows, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		mock.ExpectExec(sqlDeleteTask).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		err := sut.Delete(id)

		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
	})
}
