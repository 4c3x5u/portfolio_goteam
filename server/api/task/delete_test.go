//go:build utest

package task

import (
	"database/sql"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/api"
	"server/assert"
	taskTable "server/dbaccess/task"
	pkgLog "server/log"
)

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDELETEHandler(t *testing.T) {
	idValidator := &api.FakeStringValidator{}
	taskSelector := &taskTable.FakeSelector{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDELETEHandler(
		idValidator, taskSelector, log,
	)

	for _, c := range []struct {
		name                 string
		idValidatorErr       error
		taskSelectorErr      error
		columnSelectorErr    error
		userBoardSelectorErr error
		wantStatusCode       int
		assertFunc           func(*testing.T, *http.Response, string)
	}{
		{
			name:              "IDEmpty",
			idValidatorErr:    api.ErrStrEmpty,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			wantStatusCode:    http.StatusBadRequest,
			assertFunc:        assert.OnResErr("Task ID cannot be empty."),
		},
		{
			name:              "IDNotInt",
			idValidatorErr:    api.ErrStrNotInt,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			wantStatusCode:    http.StatusBadRequest,
			assertFunc:        assert.OnResErr("Task ID must be an integer."),
		},
		{
			name:              "IDUnexpectedErr",
			idValidatorErr:    api.ErrStrTooLong,
			taskSelectorErr:   nil,
			columnSelectorErr: nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(api.ErrStrTooLong.Error()),
		},
		{
			name:              "TaskSelectorErr",
			idValidatorErr:    nil,
			taskSelectorErr:   sql.ErrConnDone,
			columnSelectorErr: nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:              "TaskNotFound",
			idValidatorErr:    nil,
			taskSelectorErr:   sql.ErrNoRows,
			columnSelectorErr: nil,
			wantStatusCode:    http.StatusNotFound,
			assertFunc:        assert.OnResErr("Task not found."),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			idValidator.Err = c.idValidatorErr
			taskSelector.Err = c.taskSelectorErr

			r, err := http.NewRequest("", "", nil)
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
