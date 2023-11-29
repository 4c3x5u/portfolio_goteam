//go:build utest

package board

import (
	"database/sql"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

	t.Run("OK", func(t *testing.T) {
		wantRecs := []Record{
			{ID: 1, Name: "Board 1"},
			{ID: 2, Name: "Board 2"},
		}
		rows := sqlmock.NewRows([]string{"id", "name"})
		for _, rec := range wantRecs {
			rows.AddRow(rec.ID, rec.Name)
		}

		mock.ExpectQuery(sqlSelect).WithArgs(teamID).WillReturnRows(rows)

		recs, err := sut.Select(teamID)
		assert.Nil(t.Fatal, err)

		for i, wantRec := range wantRecs {
			assert.Equal(t.Error, recs[i].ID, wantRec.ID)
			assert.Equal(t.Error, recs[i].Name, wantRec.Name)
		}
	})
}
