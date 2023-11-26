//go:build utest

package user

import (
	"errors"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
)

func TestSelectorByTeamID(t *testing.T) {
	teamID := "21"
	sqlSelect := `SELECT username, isAdmin FROM app.\"user\" WHERE teamID = \$1`

	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewSelectorByTeamID(db)

	t.Run("Error", func(t *testing.T) {
		wantErr := errors.New("error selecting user")

		mock.ExpectQuery(sqlSelect).WithArgs(teamID).WillReturnError(wantErr)

		_, err := sut.Select(teamID)

		if err := assert.Equal(wantErr, err); err != nil {
			t.Error(err)
		}
	})
}
