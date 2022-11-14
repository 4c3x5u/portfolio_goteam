package register

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	"github.com/kxplxn/goteam/server-v2/assert"
)

func TestCreatorDBUser(t *testing.T) {
	t.Run("UsernameIsTaken", func(t *testing.T) {
		db, mock, err := sqlmock.New()
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()

		const (
			username            = "bobby"
			wantUsernameIsTaken = true
		)

		//goland:noinspection SqlResolve
		mock.
			ExpectQuery(`SELECT username FROM users WHERE username = \$1`).
			WithArgs(username).
			WillReturnRows(sqlmock.NewRows([]string{"username"}).AddRow(username))

		sut := NewCreatorDBUser(db)

		errUsername, err := sut.CreateUser(username, "")
		if err != nil {
			t.Fatal(err)
		}
		assert.Equal(t, wantUsernameIsTaken, errUsername)
	})
}
