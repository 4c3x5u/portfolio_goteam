//go:build utest

package subtask

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestPATCHHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	idValidator := &api.FakeStringValidator{}
	log := &pkgLog.FakeErrorer{}
	sut := NewPATCHHandler(idValidator, log)

	for _, c := range []struct {
		name           string
		idValidatorErr error
		wantStatusCode int
		assertFunc     func(*testing.T, *http.Response, string)
	}{
		{
			name:           "IDEmpty",
			idValidatorErr: api.ErrStrEmpty,
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     assert.OnResErr("Subtask ID cannot be empty."),
		},
		{
			name:           "IDNotInt",
			idValidatorErr: api.ErrStrNotInt,
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     assert.OnResErr("Subtask ID must be an integer."),
		},
		{
			name:           "IDUnexpectedErr",
			idValidatorErr: api.ErrStrTooLong,
			wantStatusCode: http.StatusInternalServerError,
			assertFunc:     assert.OnLoggedErr(api.ErrStrTooLong.Error()),
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			idValidator.Err = c.idValidatorErr

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
