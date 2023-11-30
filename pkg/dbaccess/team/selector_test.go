//go:build utest

package team

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/dbaccess"
)

// TestSelector tests Selector to assert that it queries the database correctly
// and handles the result and error appropriately.
func TestSelector(t *testing.T) {
	const (
		id            = "2"
		sqlSelectTeam = `SELECT id, inviteCode FROM app.team ` +
			`WHERE id = \$1`
	)

	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewSelector(db)

	t.Run("Error", func(t *testing.T) {
		wantErr := sql.ErrNoRows

		mock.ExpectQuery(sqlSelectTeam).
			WithArgs(id).
			WillReturnError(wantErr)

		_, err := sut.Select(id)

		assert.ErrIs(t.Error, err, wantErr)
	})

	t.Run("Success", func(t *testing.T) {
		wantInviteCode := "someinvitecode"
		mock.ExpectQuery(sqlSelectTeam).
			WithArgs(id).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "inviteCode"}).
					AddRow(id, wantInviteCode))

		team, err := sut.Select(id)

		assert.Nil(t.Error, err)
		assert.Equal(t.Error, strconv.Itoa(team.ID), id)
		assert.Equal(t.Error, team.InviteCode, wantInviteCode)
	})
}
