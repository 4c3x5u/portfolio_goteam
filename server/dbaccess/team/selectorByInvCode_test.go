//go:build utest

package team

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
)

// TestSelectorByInvCode tests SelectorByInvCode to assert that it queries the
// database correctly and handles the result and error appropriately.
func TestSelectorByInvCode(t *testing.T) {
	const (
		inviteCode    = "someinvitecode"
		sqlSelectTeam = `SELECT id, inviteCode FROM app.team ` +
			`WHERE inviteCode = \$1`
	)

	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewSelectorByInvCode(db)

	t.Run("Error", func(t *testing.T) {
		wantErr := sql.ErrNoRows

		mock.ExpectQuery(sqlSelectTeam).
			WithArgs(inviteCode).
			WillReturnError(wantErr)

		_, err := sut.Select(inviteCode)

		assert.SameError(t.Error, err, wantErr)
	})

	t.Run("Success", func(t *testing.T) {
		id := 3

		mock.ExpectQuery(sqlSelectTeam).
			WithArgs(inviteCode).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "inviteCode"}).
					AddRow(id, inviteCode),
			)

		team, err := sut.Select(inviteCode)
		assert.Nil(t.Fatal, err)

		assert.Equal(t.Error, team.ID, id)
		assert.Equal(t.Error, team.InviteCode, inviteCode)
	})
}
