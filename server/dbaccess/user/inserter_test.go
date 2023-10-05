//go:build utest

package user

import (
	"errors"
	"server/dbaccess"
	"testing"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

// TestInserter tests the Insert method of Inserter to assert that it
// sends the correct query to the database with the correct arguments, and
// returns whatever error occurs.
func TestInserter(t *testing.T) {
	const (
		username = "bob123"
		pwdHash  = "asd..fasdf.asdfa/sdf.asdfa.sdfa"
		query    = `INSERT INTO app.\"user\"\(username, password\) ` +
			`VALUES \(\$1, \$2\)`
	)

	t.Run("Error", func(t *testing.T) {
		wantErr := errors.New("db: fatal error")
		db, mock, teardown := dbaccess.SetUpDBTest(t)
		defer teardown()
		mock.ExpectExec(query).
			WithArgs(username, pwdHash).
			WillReturnError(wantErr)
		sut := NewInserter(db)

		err := sut.Insert(NewRecord(username, []byte(pwdHash)))

		if err = assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		db, mock, teardown := dbaccess.SetUpDBTest(t)
		defer teardown()
		mock.ExpectExec(query).
			WithArgs(username, pwdHash).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sut := NewInserter(db)

		err := sut.Insert(NewRecord(username, []byte(pwdHash)))

		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
	})
}
