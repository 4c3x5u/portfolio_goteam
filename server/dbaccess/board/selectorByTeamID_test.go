//go:build utest

package board

import (
	"database/sql"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
)

func TestSelectorByTeamID(t *testing.T) {
	teamID := "21"
	sqlSelect := `SELECT id, name FROM app.board WHERE teamID = \$1`

	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewSelectorByTeamID(db)

	t.Run("Error", func(t *testing.T) {
		wantErr := sql.ErrNoRows
		mock.ExpectQuery(sqlSelect).WithArgs(teamID).WillReturnError(wantErr)

		_, err := sut.Select(teamID)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})
}
