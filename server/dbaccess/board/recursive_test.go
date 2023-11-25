//go:build utest

package board

import (
	"errors"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
)

// TestRecursiveSelector tests the Select method of RecursiveSelector to assert
// that it sends the correct queries to the database with the correct arguments,
// and returns whatever error occurs alongside the correct data.
func TestRecursiveSelector(t *testing.T) {
	boardID := "21"
	sqlSelectBoard := `SELECT id, name FROM app.board WHERE id = \$1`

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

}
