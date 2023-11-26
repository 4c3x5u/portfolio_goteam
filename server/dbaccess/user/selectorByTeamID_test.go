//go:build utest

package user

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
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

	t.Run("OK", func(t *testing.T) {
		wantRecs := []Record{
			{Username: "foo", IsAdmin: true},
			{Username: "bar", IsAdmin: false},
			{Username: "baz", IsAdmin: false},
		}
		rows := sqlmock.NewRows([]string{"username", "isAdmin"})
		for _, user := range wantRecs {
			rows.AddRow(user.Username, user.IsAdmin)
		}
		mock.ExpectQuery(sqlSelect).WithArgs(teamID).WillReturnRows(rows)

		recs, err := sut.Select(teamID)
		if err = assert.Nil(err); err != nil {
			t.Error(err)
		}

		for i, user := range wantRecs {
			if err = assert.Equal(user.Username, recs[i].Username); err != nil {
				t.Error(err)
			}
			if err = assert.Equal(user.IsAdmin, recs[i].IsAdmin); err != nil {
				t.Error(err)
			}
		}
	})
}
