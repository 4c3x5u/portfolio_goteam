//go:build utest

package team

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
)

// TestSelector tests Selector to assert that it queries the database and
// handles the result and error correctly.
func TestSelector(t *testing.T) {
	const (
		inviteCode    = "someinvitecode"
		sqlSelectTeam = `SELECT id FROM app.team WHERE inviteCode = \$1`
	)

	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewSelector(db)

	t.Run("Error", func(t *testing.T) {
		wantErr := sql.ErrNoRows

		mock.ExpectQuery(sqlSelectTeam).
			WithArgs(inviteCode).
			WillReturnError(wantErr)

		_, err := sut.Select(inviteCode)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		id := 3

		mock.ExpectQuery(sqlSelectTeam).
			WithArgs(inviteCode).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(id))

		team, err := sut.Select(inviteCode)

		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(id, team.ID); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(inviteCode, team.InviteCode); err != nil {
			t.Error(err)
		}
	})
}
