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
	userboardTable "github.com/kxplxn/goteam/server/dbaccess/userboard"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDELETEHandler(t *testing.T) {
	validator := &api.FakeStringValidator{}
	userBoardSelector := &userboardTable.FakeSelector{}
	userBoardDeleter := &dbaccess.FakeDeleter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDELETEHandler(
		validator, userBoardSelector, userBoardDeleter, log,
	)

	// Used on cases where no case-specific assertions are required.
	emptyAssertFunc := func(*testing.T, *http.Response, string) {}

	for _, c := range []struct {
		name                 string
		validatorErr         error
		userIsAdmin          bool
		userBoardSelectorErr error
		boardDeleterErr      error
		wantStatusCode       int
		assertFunc           func(*testing.T, *http.Response, string)
	}{
		{
			name:                 "ValidatorErr",
			validatorErr:         errors.New("some validator err"),
			userIsAdmin:          true,
			userBoardSelectorErr: nil,
			boardDeleterErr:      nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc:           emptyAssertFunc,
		},
		{
			name:                 "ConnDone",
			validatorErr:         nil,
			userIsAdmin:          false,
			userBoardSelectorErr: sql.ErrConnDone,
			boardDeleterErr:      nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				sql.ErrConnDone.Error(),
			),
		},
		{
			name:                 "UserBoardNotFound",
			validatorErr:         nil,
			userIsAdmin:          false,
			userBoardSelectorErr: sql.ErrNoRows,
			boardDeleterErr:      nil,
			wantStatusCode:       http.StatusForbidden,
			assertFunc:           emptyAssertFunc,
		},
		{
			name:                 "NotAdmin",
			validatorErr:         nil,
			userIsAdmin:          false,
			userBoardSelectorErr: nil,
			boardDeleterErr:      nil,
			wantStatusCode:       http.StatusForbidden,
			assertFunc:           emptyAssertFunc,
		},
		{
			name:                 "DeleteErr",
			validatorErr:         nil,
			userIsAdmin:          true,
			userBoardSelectorErr: nil,
			boardDeleterErr:      errors.New("delete board error"),
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: assert.OnLoggedErr(
				"delete board error",
			),
		},
		{
			name:                 "Success",
			validatorErr:         nil,
			userIsAdmin:          true,
			userBoardSelectorErr: nil,
			boardDeleterErr:      nil,
			wantStatusCode:       http.StatusOK,
			assertFunc:           emptyAssertFunc,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for sut's dependencies.
			validator.Err = c.validatorErr
			userBoardSelector.IsAdmin = c.userIsAdmin
			userBoardSelector.Err = c.userBoardSelectorErr
			userBoardDeleter.Err = c.boardDeleterErr

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
