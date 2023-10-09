//go:build utest

package task

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"server/assert"
	"server/dbaccess"
	"testing"
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
	if err := assert.SameError(sql.ErrNoRows, err); err != nil {
		t.Error(err)
	}

	task, err = sut.Select(id)
	if err = assert.Nil(err); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(wantID, task.ID); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(wantColumnID, task.ColumnID); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(wantDescription, *task.Description); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(wantOrder, task.Order); err != nil {
		t.Error(err)
	}
}
