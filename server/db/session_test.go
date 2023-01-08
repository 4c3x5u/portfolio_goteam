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
		db, mock, teardown := setupTest(t)
		defer teardown()
		mock.
			ExpectExec(query).
			WithArgs(id, username, expiry.String()).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sut := NewCreatorSession(db)

		err := sut.Create(NewSession(id, username, expiry))

		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
	})

	t.Run("CreateErr", func(t *testing.T) {
		wantErr := errors.New("db: fatal error")
		db, mock, teardown := setupTest(t)
		defer teardown()
		mock.
			ExpectExec(query).
			WithArgs(id, username, expiry.String()).
			WillReturnError(wantErr)
		sut := NewCreatorSession(db)

		err := sut.Create(NewSession(id, username, expiry))

		if err = assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})
}

func TestUpserterSession(t *testing.T) {
	var (
		session = NewSession("3aad526f-afea-4d07-986a-72fcf245bd18", "bob21", time.Now().Add(1*time.Hour))
		query   = `INSERT INTO sessions\(id, username, expiry\) VALUES \(\$1, \$2, \$3\) ON CONFLICT \(username\) DO UPDATE SET expiry = \$3`
	)

	t.Run("UpsertOK", func(t *testing.T) {
		db, mock, teardown := setupTest(t)
		defer teardown()
		mock.
			ExpectExec(query).
			WithArgs(session.ID, session.Username, session.Expiry.String()).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sut := NewUpserterSession(db)

		err := sut.Upsert(session)

		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
	})

	t.Run("UpsertErr", func(t *testing.T) {
		wantErr := errors.New("db: fatal error")
		db, mock, teardown := setupTest(t)
		defer teardown()
		mock.
			ExpectExec(query).
			WithArgs(session.ID, session.Username, session.Expiry.String()).
			WillReturnError(wantErr)
		sut := NewUpserterSession(db)

		err := sut.Upsert(session)

		if err = assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})
}
