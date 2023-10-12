//go:build utest

package subtask

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	subtaskTable "github.com/kxplxn/goteam/server/dbaccess/subtask"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestPATCHHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	idValidator := &api.FakeStringValidator{}
	subtaskSelector := &subtaskTable.FakeSelector{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(idValidator, subtaskSelector, log)

	for _, c := range []struct {
		name               string
		idValidatorErr     error
		subtaskSelectorErr error
		wantStatusCode     int
		assertFunc         func(*testing.T, *http.Response, string)
	}{
		{
			name:               "IDEmpty",
			idValidatorErr:     api.ErrStrEmpty,
			subtaskSelectorErr: nil,
			wantStatusCode:     http.StatusBadRequest,
			assertFunc:         assert.OnResErr("Subtask ID cannot be empty."),
		},
		{
			name:               "IDNotInt",
			idValidatorErr:     api.ErrStrNotInt,
			subtaskSelectorErr: nil,
			wantStatusCode:     http.StatusBadRequest,
			assertFunc:         assert.OnResErr("Subtask ID must be an integer."),
		},
		{
			name:               "IDUnexpectedErr",
			idValidatorErr:     api.ErrStrTooLong,
			subtaskSelectorErr: nil,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(api.ErrStrTooLong.Error()),
		},
		{
			name:               "SubtaskSelectorErr",
			idValidatorErr:     nil,
			subtaskSelectorErr: sql.ErrConnDone,
			wantStatusCode:     http.StatusInternalServerError,
			assertFunc:         assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:               "SubtaskNotFound",
			idValidatorErr:     nil,
			subtaskSelectorErr: sql.ErrNoRows,
			wantStatusCode:     http.StatusNotFound,
			assertFunc:         assert.OnResErr("Subtask not found."),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			idValidator.Err = c.idValidatorErr
			subtaskSelector.Err = c.subtaskSelectorErr

			r, err := http.NewRequest("", "?id=", nil)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			sut.Handle(w, r, "")
			res := w.Result()

			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			c.assertFunc(t, res, log.InMessage)
		})
	}
}
