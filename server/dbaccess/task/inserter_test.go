//go:build utest

package task

import (
	"database/sql"
	"errors"
	"testing"

	"server/assert"
	"server/dbaccess"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestInserter tests the Insert method of Inserter to assert that it sends the
// correct queries to the database with the correct arguments, and returns
// whatever error occurs.
func TestInserter(t *testing.T) {
	const (
		sqlSelectOrder = `SELECT "order" FROM app.task WHERE columnID = \$1 ` +
			`ORDER BY "order" DESC LIMIT 1`
		sqlInsertTask = `INSERT INTO app.task` +
			`\(columnID, title, description, \"order\"\)` +
			`VALUES \(\$1, \$2, \$3, \$4\)`
	)

	task := NewTask(
		2, "Task A", "Description A", []string{"Subtask A", "Subtask B"},
	)
	anErr := errors.New("an error occurred")

	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewInserter(db)

	for _, c := range []struct {
		name      string
		setUpMock func(sqlmock.Sqlmock)
		wantErr   error
	}{
		{
			name: "BeginErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin().WillReturnError(anErr)
			},
			wantErr: anErr,
		},
		{
			name: "SelectOrderErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectOrder).
					WithArgs(task.columnID).
					WillReturnError(anErr)
			},
			wantErr: anErr,
		},
		{
			name: "SelectOrderNoRows",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectOrder).
					WithArgs(task.columnID).
					WillReturnError(sql.ErrNoRows)
				mock.ExpectExec(sqlInsertTask).
					WithArgs(task.columnID, task.title, task.description, 1).
					WillReturnError(anErr)
			},
			wantErr: anErr,
		},
		{
			name: "InsertTaskErr",
			setUpMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectBegin()
				mock.ExpectQuery(sqlSelectOrder).
					WithArgs(task.columnID).
					WillReturnRows(sqlmock.NewRows([]string{"order"}).AddRow(5))
				mock.ExpectExec(sqlInsertTask).
					WithArgs(task.columnID, task.title, task.description, 6).
					WillReturnError(anErr)
			},
			wantErr: anErr,
		},
	} {
		c.setUpMock(mock)
		err := sut.Insert(task)
		if assertErr := assert.SameError(c.wantErr, err); assertErr != nil {
			t.Error(assertErr)
		}
	}
}
