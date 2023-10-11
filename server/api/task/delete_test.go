//go:build utest

package task

import (
	"net/http"
	"net/http/httptest"
	"server/api"
	"server/assert"
	pkgLog "server/log"
	"testing"
)

// TestDELETEHandler tests the Handle method of DELETEHandler to assert that it
// behaves correctly in all possible scenarios.
func TestDELETEHandler(t *testing.T) {
	idValidator := &api.FakeStringValidator{}
	log := &pkgLog.FakeErrorer{}
	sut := NewDELETEHandler(idValidator, log)

	for _, c := range []struct {
		name           string
		idValidatorErr error
		wantStatusCode int
		wantErrMsg     string
	}{
		{
			name:           "IDEmpty",
			idValidatorErr: api.ErrStrEmpty,
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     "Task ID cannot be empty.",
		},
		{
			name:           "IDNotInt",
			idValidatorErr: api.ErrStrNotInt,
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     "Task ID must be an integer.",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			idValidator.Err = c.idValidatorErr

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

			assert.OnResErr(c.wantErrMsg)(t, res, "")
		})
	}
}
