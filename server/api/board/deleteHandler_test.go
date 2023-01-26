package board

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/db"
	"server/log"
)

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDELETEHandler(t *testing.T) {
	userBoardSelector := &db.FakeRelSelector{}
	userBoardDeleter := &db.FakeDeleter{}
	logger := &log.FakeLogger{}
	// sub is inconsequential since errors returned from dependencies are
	// directly manipulated
	sub := ""
	url := "/board?id=123"
	sut := NewDeleteHandler(userBoardSelector, userBoardDeleter, logger)

	for _, c := range []struct {
		name                        string
		userBoardSelectorOutIsAdmin bool
		userBoardSelectorOutErr     error
		boardDeleterOutErr          error
		wantStatusCode              int
	}{
		{
			name:                        "NoRows",
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     sql.ErrNoRows,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusNotFound,
		},
		{
			name:                        "ConnDone",
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     sql.ErrConnDone,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusInternalServerError,
		},
		{
			name:                        "NotAdmin",
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     nil,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusUnauthorized,
		},
		{
			name:                        "DeleteErr",
			userBoardSelectorOutIsAdmin: true,
			userBoardSelectorOutErr:     nil,
			boardDeleterOutErr:          errors.New("delete board error"),
			wantStatusCode:              http.StatusInternalServerError,
		},
		{
			name:                        "Success",
			userBoardSelectorOutIsAdmin: true,
			userBoardSelectorOutErr:     nil,
			boardDeleterOutErr:          nil,
			wantStatusCode:              http.StatusOK,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userBoardSelector.OutIsAdmin = c.userBoardSelectorOutIsAdmin
			userBoardSelector.OutErr = c.userBoardSelectorOutErr
			userBoardDeleter.OutErr = c.boardDeleterOutErr

			req, err := http.NewRequest(http.MethodPost, url, nil)
			if err != nil {
				t.Fatal(err)
			}

			w := httptest.NewRecorder()

			sut.Handle(w, req, sub)

			if err := assert.Equal(c.wantStatusCode, w.Result().StatusCode); err != nil {
				t.Error(err)
			}

			if c.wantStatusCode == http.StatusInternalServerError {
				errFound := false
				for _, err := range []error{
					c.userBoardSelectorOutErr,
					c.boardDeleterOutErr,
				} {
					if err != nil {
						errFound = true
						if err := assert.Equal(
							log.LevelError, logger.InLevel,
						); err != nil {
							t.Error(err)
						}
						if err := assert.Equal(err.Error(), logger.InMessage); err != nil {
							t.Error(err)
						}
					}
				}
				if !errFound {
					t.Errorf(
						"c.wantStatusCode was %d but no errors were logged.",
						http.StatusInternalServerError,
					)
				}
			}
		})
	}
}
