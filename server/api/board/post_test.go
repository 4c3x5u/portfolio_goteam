//go:build utest

package board

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/dbaccess"
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestPOSTHandler tests the Handle method of POSTHandler to assert that it
// behaves correctly in all possible scenarios.
func TestPOSTHandler(t *testing.T) {
	userSelector := &userTable.FakeSelector{}
	validator := &api.FakeStringValidator{}
	boardCounter := &dbaccess.FakeCounter{}
	boardInserter := &boardTable.FakeInserter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPOSTHandler(
		userSelector, validator, boardCounter, boardInserter, log,
	)

	t.Run(http.MethodPost, func(t *testing.T) {
		for _, c := range []struct {
			name             string
			user             userTable.Record
			userSelectorErr  error
			validatorErr     error
			boardCount       int
			boardCounterErr  error
			boardInserterErr error
			wantStatusCode   int
			assertFunc       func(*testing.T, *http.Response, string)
		}{
			{
				name:             "UserNotRecognised",
				user:             userTable.Record{},
				userSelectorErr:  sql.ErrNoRows,
				validatorErr:     nil,
				boardCount:       0,
				boardCounterErr:  nil,
				boardInserterErr: nil,
				wantStatusCode:   http.StatusUnauthorized,
				assertFunc: assert.OnResErr(
					"Username is not recognised.",
				),
			},
			{
				name:             "UserSelectorErr",
				user:             userTable.Record{IsAdmin: false},
				userSelectorErr:  sql.ErrConnDone,
				validatorErr:     nil,
				boardCount:       0,
				boardCounterErr:  nil,
				boardInserterErr: nil,
				wantStatusCode:   http.StatusInternalServerError,
				assertFunc: assert.OnLoggedErr(
					sql.ErrConnDone.Error(),
				),
			},
			{
				name:             "UserNotAdmin",
				user:             userTable.Record{IsAdmin: false},
				userSelectorErr:  nil,
				validatorErr:     nil,
				boardCount:       0,
				boardCounterErr:  sql.ErrConnDone,
				boardInserterErr: nil,
				wantStatusCode:   http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can create boards.",
				),
			},
			{
				name:             "InvalidRequest",
				user:             userTable.Record{IsAdmin: true},
				userSelectorErr:  nil,
				validatorErr:     errors.New("Board name cannot be empty."),
				boardCount:       0,
				boardCounterErr:  nil,
				boardInserterErr: nil,
				wantStatusCode:   http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"Board name cannot be empty.",
				),
			},
			{
				name:             "BoardCounterErr",
				user:             userTable.Record{IsAdmin: true},
				userSelectorErr:  nil,
				validatorErr:     nil,
				boardCount:       0,
				boardCounterErr:  sql.ErrConnDone,
				boardInserterErr: nil,
				wantStatusCode:   http.StatusInternalServerError,
				assertFunc: assert.OnLoggedErr(
					sql.ErrConnDone.Error(),
				),
			},
			{
				name:             "MaxBoardsCreated",
				user:             userTable.Record{IsAdmin: true},
				userSelectorErr:  nil,
				validatorErr:     nil,
				boardCount:       3,
				boardCounterErr:  nil,
				boardInserterErr: nil,
				wantStatusCode:   http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"You have already created the maximum amount of boards " +
						"allowed per user. Please delete one of your boards " +
						"to create a new one.",
				),
			},
			{
				name:             "BoardInserterErr",
				user:             userTable.Record{IsAdmin: true},
				userSelectorErr:  nil,
				validatorErr:     nil,
				boardCount:       0,
				boardCounterErr:  sql.ErrNoRows,
				boardInserterErr: errors.New("create board error"),
				wantStatusCode:   http.StatusInternalServerError,
				assertFunc: assert.OnLoggedErr(
					"create board error",
				),
			},
			{
				name:             "Success",
				user:             userTable.Record{IsAdmin: true},
				userSelectorErr:  nil,
				validatorErr:     nil,
				boardCount:       0,
				boardCounterErr:  sql.ErrNoRows,
				boardInserterErr: nil,
				wantStatusCode:   http.StatusOK,
				assertFunc: func(*testing.T, *http.Response, string) {
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				userSelector.Rec = c.user
				userSelector.Err = c.userSelectorErr
				validator.Err = c.validatorErr
				boardCounter.BoardCount = c.boardCount
				boardCounter.Err = c.boardCounterErr
				boardInserter.Err = c.boardInserterErr

				req := httptest.NewRequest("", "/", strings.NewReader("{}"))
				w := httptest.NewRecorder()

				sut.Handle(w, req, "")
				res := w.Result()

				// Assert on the status code.
				assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

				// Run case-specific assertions.
				c.assertFunc(t, res, log.InMessage)
			})
		}
	})
}
