//go:build utest

package team

import (
	"database/sql"
	"strconv"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
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

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		wantInviteCode := "someinvitecode"
		mock.ExpectQuery(sqlSelectTeam).
			WithArgs(id).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "inviteCode"}).
					AddRow(id, wantInviteCode))

		team, err := sut.Select(id)

		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(id, strconv.Itoa(team.ID)); err != nil {
			t.Error(err)
		}
		if err = assert.Equal(wantInviteCode, team.InviteCode); err != nil {
			t.Error(err)
		}
	})
}

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

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
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
