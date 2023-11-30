//go:build utest

package task

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/dbaccess"
)

// TestSelector tests the Select method of Selector to assert that it sends the
// correct query to the database with the correct arguments, and returns all
// data or whatever error occurs.
func TestSelector(t *testing.T) {
	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewSelector(db)

	const (
		sqlSelectTask = `SELECT id, columnID, title, description, \"order\" ` +
			`FROM app.task WHERE id = \$1`
		id              = "1"
		wantID          = 1
		wantColumnID    = 2
		wantTitle       = "Some Task"
		wantDescription = "Do something."
		wantOrder       = 3
	)

	mock.ExpectQuery(sqlSelectTask).
		WithArgs(id).
		WillReturnError(sql.ErrNoRows)

	mock.ExpectQuery(sqlSelectTask).
		WithArgs(id).
		WillReturnRows(
			sqlmock.
				NewRows(
					[]string{"id", "columnID", "title", "description", "order"},
				).
				AddRow(
					wantID, wantColumnID, wantTitle, wantDescription, wantOrder,
				),
		)

	task, err := sut.Select(id)
	assert.ErrIs(t.Error, err, sql.ErrNoRows)

	task, err = sut.Select(id)
	assert.Nil(t.Fatal, err)

	assert.Equal(t.Error, task.ID, wantID)
	assert.Equal(t.Error, task.ColumnID, wantColumnID)
	assert.Equal(t.Error, *task.Description, wantDescription)
	assert.Equal(t.Error, task.Order, wantOrder)
}
