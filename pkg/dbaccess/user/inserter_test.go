//go:build utest

package user

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/dbaccess"
)

// TestInserter tests the Insert method of Inserter to assert that it
// sends the correct query to the database with the correct arguments, and
// returns whatever error occurs.
func TestInserter(t *testing.T) {
	const (
		username      = "bob123"
		pwdHash       = "asd..fasdf.asdfa/sdf.asdfa.sdfa"
		cmdInsertUser = `INSERT INTO app.\"user\"\(username, password, ` +
			`teamID, isAdmin\) VALUES \(\$1, \$2\, \$3\, \$4\)`
		cmdInsertTeam = `INSERT INTO app.team\(inviteCode\) VALUES \(\$1\) ` +
			`RETURNING id`
	)

	t.Run("ErrBegin", func(t *testing.T) {
		rec := NewRecord("", []byte{}, -1, false)
		errBegin := errors.New("error beginning transaction")
		db, mock, teardown := dbaccess.SetUpDBTest(t)
		defer teardown()
		mock.ExpectBegin().WillReturnError(errBegin)
		sut := NewInserter(db)

		err := sut.Insert(rec)

		if !errors.Is(err, errBegin) {
			t.Errorf("wrong error - got: %s, want: %s", err, errBegin)
		}
	})

	t.Run("ErrInsertTeam", func(t *testing.T) {
		rec := NewRecord("", []byte{}, -1, true)
		errInsertTeam := errors.New("error inserting team")
		db, mock, teardown := dbaccess.SetUpDBTest(t)
		defer teardown()
		mock.ExpectBegin()
		mock.ExpectQuery(cmdInsertTeam).WillReturnError(errInsertTeam)
		mock.ExpectRollback()
		sut := NewInserter(db)

		err := sut.Insert(rec)

		if !errors.Is(err, errInsertTeam) {
			t.Errorf("wrong error - got: %s, want: %s", err, errInsertTeam)
		}
	})

	t.Run("ErrInsertUser", func(t *testing.T) {
		rec := NewRecord(username, []byte(pwdHash), 21, false)
		wantErr := errors.New("error inserting user")
		db, mock, teardown := dbaccess.SetUpDBTest(t)
		defer teardown()
		mock.ExpectBegin()
		mock.ExpectExec(cmdInsertUser).
			WithArgs(username, pwdHash, 21, false).
			WillReturnError(wantErr)
		mock.ExpectRollback()
		sut := NewInserter(db)

		err := sut.Insert(rec)

		assert.Equal(t.Error, err, wantErr)
	})

	t.Run("OKWithTeam", func(t *testing.T) {
		rec := NewRecord(username, []byte(pwdHash), -1, true)
		db, mock, teardown := dbaccess.SetUpDBTest(t)
		defer teardown()
		mock.ExpectBegin()
		mock.ExpectQuery(cmdInsertTeam).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(21))
		mock.ExpectExec(cmdInsertUser).
			WithArgs(username, pwdHash, 21, true).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		sut := NewInserter(db)

		err := sut.Insert(rec)

		assert.Nil(t.Error, err)
	})

	t.Run("OK", func(t *testing.T) {
		rec := NewRecord(username, []byte(pwdHash), 32, false)
		db, mock, teardown := dbaccess.SetUpDBTest(t)
		defer teardown()
		mock.ExpectBegin()
		mock.ExpectExec(cmdInsertUser).
			WithArgs(username, pwdHash, rec.TeamID, false).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()
		sut := NewInserter(db)

		err := sut.Insert(rec)

		assert.Nil(t.Error, err)
	})
}
