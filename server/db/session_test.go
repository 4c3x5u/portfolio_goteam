package db

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"server/assert"
	"testing"
	"time"
)

func TestCreatorSession(t *testing.T) {
	id, username, expiry := "ba929ec9-91a4-4fca-b6ec-edf10d6c827e", "bob21", time.Now()
	query := `INSERT INTO sessions\(id, username, expiry\) VALUES \(\$1, \$2, \$3\)`

	t.Run("CreateOK", func(t *testing.T) {
		db, mock, def := setup(t)
		defer def(db)
		mock.
			ExpectExec(query).
			WithArgs(id, username, expiry.String()).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sut := NewCreatorSession(db)

		err := sut.Create(id, username, expiry)

		assert.Nil(t, err)
	})

	t.Run("CreateErr", func(t *testing.T) {
		wantErr := errors.New("db: fatal error")
		db, mock, def := setup(t)
		defer def(db)
		mock.
			ExpectExec(query).
			WithArgs(id, username, expiry.String()).
			WillReturnError(wantErr)
		sut := NewCreatorSession(db)

		err := sut.Create(id, username, expiry)

		assert.Equal(t, wantErr.Error(), err.Error())
	})
}
