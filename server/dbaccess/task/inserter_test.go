//go:build utest

package task

import (
	"errors"
	"server/assert"
	"server/dbaccess"
	"testing"
)

// TestInserter tests the Insert method of Inserter to assert that it sends the
// correct queries to the database with the correct arguments, and returns
// whatever error occurs.
func TestInserter(t *testing.T) {
	task := NewTask()
	wantErr := errors.New("an error occurred")

	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewInserter(db)

	mock.ExpectBegin().WillReturnError(wantErr)

	err := sut.Insert(task)

	if assertErr := assert.SameError(err, wantErr); assertErr != nil {
		t.Error(assertErr)
	}
}
