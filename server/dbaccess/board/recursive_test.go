//go:build utest

package board

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
)

// TestRecursiveSelector tests the Select method of RecursiveSelector to assert
// that it sends the correct queries to the database with the correct arguments,
// and returns whatever error occurs alongside the correct data.
func TestRecursiveSelector(t *testing.T) {
	boardID := "21"
	sqlSelectBoard := `SELECT id, name, teamID FROM app.board WHERE id = \$1`
	sqlSelectColumn := `SELECT id, \"order\" FROM app.column ` +
		`WHERE boardID = \$1`
	sqlSelectTask := `SELECT id, title, description, \"order\" FROM app.task ` +
		`WHERE columnID = \$1`

	db, mock, teardown := dbaccess.SetUpDBTest(t)
	defer teardown()

	sut := NewRecursiveSelector(db)

	t.Run("SelectBoardErr", func(t *testing.T) {
		wantErr := errors.New("error selecting board")
		mock.ExpectQuery(sqlSelectBoard).
			WithArgs(boardID).
			WillReturnError(wantErr)

		_, err := sut.Select(boardID)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("SelectColumnErr", func(t *testing.T) {
		wantErr := errors.New("error selecting column")
		mock.ExpectQuery(sqlSelectBoard).
			WithArgs(boardID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "teamID"}).
					AddRow(1, "board 1", 1),
			)
		mock.ExpectQuery(sqlSelectColumn).
			WithArgs(1).
			WillReturnError(wantErr)

		_, err := sut.Select(boardID)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})

	t.Run("SelectTaskErr", func(t *testing.T) {
		wantErr := errors.New("error selecting task")
		mock.ExpectQuery(sqlSelectBoard).
			WithArgs(boardID).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "name", "teamID"}).
					AddRow(1, "board 1", 1),
			)
		mock.ExpectQuery(sqlSelectColumn).
			WithArgs(1).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "order"}).
					AddRow(2, 1).AddRow(3, 2).AddRow(4, 3).AddRow(5, 4),
			)
		mock.ExpectQuery(sqlSelectTask).
			WithArgs(2).
			WillReturnError(wantErr)

		_, err := sut.Select(boardID)

		if err = assert.SameError(wantErr, err); err != nil {
			t.Error(err)
		}
	})
}
