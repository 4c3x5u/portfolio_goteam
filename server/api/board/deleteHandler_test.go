package board

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/db"
)

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that
// it behaves correctly.
func TestDELETEHandler(t *testing.T) {
	userBoardSelector := &db.FakeRelSelector{}
	url := "/board?id=123"
	sut := NewDELETEHandler(userBoardSelector)

	for _, c := range []struct {
		name                        string
		userBoardSelectorOutIsAdmin bool
		userBoardSelectorOutErr     error
		wantStatusCode              int
	}{
		{
			name:                        "NoRows",
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     sql.ErrNoRows,
			wantStatusCode:              http.StatusNotFound,
		},
		{
			name:                        "ConnDone",
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     sql.ErrConnDone,
			wantStatusCode:              http.StatusInternalServerError,
		},
		{
			name:                        "NotAdmin",
			userBoardSelectorOutIsAdmin: false,
			userBoardSelectorOutErr:     nil,
			wantStatusCode:              http.StatusUnauthorized,
		},
		{
			name:                        "Success",
			userBoardSelectorOutIsAdmin: true,
			userBoardSelectorOutErr:     nil,
			// todo: set 200 when deletehandler work is done
			wantStatusCode: http.StatusNotImplemented,
		},
	} {
		userBoardSelector.OutIsAdmin = c.userBoardSelectorOutIsAdmin
		userBoardSelector.OutErr = c.userBoardSelectorOutErr

		req, err := http.NewRequest(http.MethodPost, url, nil)
		if err != nil {
			t.Fatal(err)
		}

		w := httptest.NewRecorder()

		sut.Handle(w, req, "bob123")

		if err := assert.Equal(c.wantStatusCode, w.Result().StatusCode); err != nil {
			t.Error(err)
		}
	}
}
