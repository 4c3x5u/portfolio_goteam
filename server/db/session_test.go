package db

import (
	"errors"
	"testing"
	"time"

	"server/assert"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestCreatorSession(t *testing.T) {
	var (
		id       = "ba929ec9-91a4-4fca-b6ec-edf10d6c827e"
		username = "bob21"
		expiry   = time.Now().Add(1 * time.Hour)
		query    = `INSERT INTO sessions\(id, username, expiry\) VALUES \(\$1, \$2, \$3\)`
	)

	t.Run("CreateOK", func(t *testing.T) {
		db, mock, def := setupTest(t)
		defer def(db)
		mock.
			ExpectExec(query).
			WithArgs(id, username, expiry.String()).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sut := NewCreatorSession(db)

		err := sut.Create(NewSession(id, username, expiry))

		assert.Nil(t, err)
	})

	t.Run("CreateErr", func(t *testing.T) {
		wantErr := errors.New("db: fatal error")
		db, mock, def := setupTest(t)
		defer def(db)
		mock.
			ExpectExec(query).
			WithArgs(id, username, expiry.String()).
			WillReturnError(wantErr)
		sut := NewCreatorSession(db)

		err := sut.Create(NewSession(id, username, expiry))

		assert.Equal(t, wantErr.Error(), err.Error())
	})
}

func TestReaderSession(t *testing.T) {
	var (
		id       = "f06e2d5c-68bc-458f-8c7e-2a73b53543f5"
		username = "bob21"
		expiry   = time.Now().Add(30 * time.Minute)
		query    = `SELECT id, username, expiry FROM sessions WHERE username = \$1`
	)

	t.Run("ReadOK", func(t *testing.T) {
		db, mock, def := setupTest(t)
		defer def(db)
		mock.
			ExpectQuery(query).
			WithArgs(username).
			WillReturnRows(mock.
				NewRows([]string{"id", "username", "expiry"}).
				AddRow(id, username, expiry),
			)
		sut := NewReaderSession(db)

		session, err := sut.Read(username)

		assert.Nil(t, err)
		assert.Equal(t, id, session.ID)
		assert.Equal(t, username, session.Username)
		assert.Equal(t, expiry, session.Expiry)
	})

	t.Run("ReadErr", func(t *testing.T) {
		wantErr := errors.New("db: fatal error")
		db, mock, def := setupTest(t)
		defer def(db)
		mock.
			ExpectQuery(query).
			WithArgs(username).
			WillReturnError(wantErr)
		sut := NewReaderSession(db)

		_, err := sut.Read(username)

		assert.Equal(t, wantErr.Error(), err.Error())
	})
}
