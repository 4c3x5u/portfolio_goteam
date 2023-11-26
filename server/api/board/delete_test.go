//go:build utest

package board

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDELETEHandler(t *testing.T) {
	validator := &api.FakeStringValidator{}
	userSelector := &userTable.FakeSelector{}
	boardSelector := &boardTable.FakeSelector{}
	userBoardDeleter := &dbaccess.FakeDeleter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDELETEHandler(
		userSelector, validator, boardSelector, userBoardDeleter, log,
	)

	// Used on cases where no case-specific assertions are required.
	emptyAssertFunc := func(*testing.T, *http.Response, string) {}

	for _, c := range []struct {
		name           string
		user           userTable.Record
		selectUserErr  error
		validatorErr   error
		board          boardTable.Record
		selectBoardErr error
		deleteBoardErr error
		wantStatusCode int
		assertFunc     func(*testing.T, *http.Response, string)
	}{
		{
			name:           "UserNotFound",
			user:           userTable.Record{},
			selectUserErr:  sql.ErrNoRows,
			validatorErr:   nil,
			board:          boardTable.Record{},
			selectBoardErr: nil,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "NotAdmin",
			user:           userTable.Record{IsAdmin: false},
			selectUserErr:  nil,
			validatorErr:   nil,
			board:          boardTable.Record{},
			selectBoardErr: nil,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "SelectUserErr",
			user:           userTable.Record{},
			selectUserErr:  sql.ErrConnDone,
			validatorErr:   nil,
			board:          boardTable.Record{},
			selectBoardErr: nil,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:           "ValidatorErr",
			user:           userTable.Record{IsAdmin: true},
			selectUserErr:  nil,
			validatorErr:   errors.New("some validator err"),
			board:          boardTable.Record{},
			selectBoardErr: nil,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "BoardNotFound",
			user:           userTable.Record{IsAdmin: true},
			selectUserErr:  nil,
			validatorErr:   nil,
			board:          boardTable.Record{},
			selectBoardErr: sql.ErrNoRows,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusNotFound,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "SelectBoardErr",
			user:           userTable.Record{IsAdmin: true},
			selectUserErr:  nil,
			validatorErr:   nil,
			board:          boardTable.Record{},
			selectBoardErr: sql.ErrConnDone,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusInternalServerError,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "BoardWrongTeam",
			user:           userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:  nil,
			validatorErr:   nil,
			board:          boardTable.Record{TeamID: 2},
			selectBoardErr: nil,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "DeleteErr",
			user:           userTable.Record{IsAdmin: true, TeamID: 1},
			validatorErr:   nil,
			selectUserErr:  nil,
			board:          boardTable.Record{TeamID: 1},
			selectBoardErr: nil,
			deleteBoardErr: errors.New("delete board error"),
			wantStatusCode: http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				"delete board error",
			),
		},
		{
			name:           "Success",
			user:           userTable.Record{IsAdmin: true, TeamID: 1},
			selectUserErr:  nil,
			validatorErr:   nil,
			board:          boardTable.Record{TeamID: 1},
			selectBoardErr: nil,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusOK,
			assertFunc:     emptyAssertFunc,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for sut's dependencies.
			userSelector.Rec = c.user
			userSelector.Err = c.selectUserErr
			validator.Err = c.validatorErr
			boardSelector.Board = c.board
			boardSelector.Err = c.selectBoardErr
			userBoardDeleter.Err = c.deleteBoardErr

			// Prepare request and response recorder.
			req, err := http.NewRequest(http.MethodPost, "", nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.Handle(w, req, "")
			res := w.Result()

			// Assert on the status code.
			if err := assert.Equal(c.wantStatusCode, res.StatusCode); err != nil {
				t.Error(err)
			}

			// Run case-specific assertions.
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
