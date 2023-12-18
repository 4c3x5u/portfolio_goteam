//go:build utest

package board

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/legacydb"
	boardTable "github.com/kxplxn/goteam/pkg/legacydb/board"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDELETEHandler(t *testing.T) {
	validator := &api.FakeStringValidator{}
	decodeAuth := &token.FakeDecode[token.Auth]{}
	boardSelector := &boardTable.FakeSelector{}
	userBoardDeleter := &legacydb.FakeDeleter{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDELETEHandler(
		decodeAuth.Func,
		validator,
		boardSelector,
		userBoardDeleter,
		log,
	)

	// Used on cases where no case-specific assertions are required.
	emptyAssertFunc := func(*testing.T, *http.Response, string) {}

	for _, c := range []struct {
		name           string
		authToken      string
		errDecodeAuth  error
		authDecoded    token.Auth
		validatorErr   error
		board          boardTable.Record
		selectBoardErr error
		deleteBoardErr error
		wantStatusCode int
		assertFunc     func(*testing.T, *http.Response, string)
	}{
		{
			name:           "NoAuth",
			authToken:      "",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{},
			validatorErr:   nil,
			board:          boardTable.Record{},
			selectBoardErr: nil,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusUnauthorized,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "InvalidAuth",
			authToken:      "nonempty",
			errDecodeAuth:  token.ErrInvalid,
			authDecoded:    token.Auth{},
			validatorErr:   nil,
			board:          boardTable.Record{},
			selectBoardErr: nil,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusUnauthorized,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "NotAdmin",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: false},
			validatorErr:   nil,
			board:          boardTable.Record{},
			selectBoardErr: nil,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "ValidatorErr",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			validatorErr:   errors.New("some validator err"),
			board:          boardTable.Record{},
			selectBoardErr: nil,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "BoardNotFound",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			validatorErr:   nil,
			board:          boardTable.Record{},
			selectBoardErr: sql.ErrNoRows,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusNotFound,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "SelectBoardErr",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			validatorErr:   nil,
			board:          boardTable.Record{},
			selectBoardErr: sql.ErrConnDone,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusInternalServerError,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "BoardWrongTeam",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true},
			validatorErr:   nil,
			board:          boardTable.Record{TeamID: 2},
			selectBoardErr: nil,
			deleteBoardErr: nil,
			wantStatusCode: http.StatusForbidden,
			assertFunc:     emptyAssertFunc,
		},
		{
			name:           "DeleteErr",
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true, TeamID: "1"},
			validatorErr:   nil,
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
			authToken:      "nonempty",
			errDecodeAuth:  nil,
			authDecoded:    token.Auth{IsAdmin: true, TeamID: "1"},
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
			decodeAuth.Err = c.errDecodeAuth
			decodeAuth.Res = c.authDecoded
			validator.Err = c.validatorErr
			boardSelector.Board = c.board
			boardSelector.Err = c.selectBoardErr
			userBoardDeleter.Err = c.deleteBoardErr

			// Prepare request and response recorder.
			r := httptest.NewRequest(http.MethodPost, "/", nil)
			if c.authToken != "" {
				r.AddCookie(&http.Cookie{
					Name:  "auth-token",
					Value: c.authToken,
				})
			}
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.Handle(w, r, "")
			res := w.Result()

			// Assert on the status code.
			assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

			// Run case-specific assertions.
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
