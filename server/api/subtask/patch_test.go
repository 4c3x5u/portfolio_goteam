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

	t.Run("IDEmpty", func(t *testing.T) {
		idValidator.Err = api.ErrStrEmpty

		wantStatusCode := http.StatusBadRequest
		wantErrMsg := "Subtask ID cannot be empty."

		r, err := http.NewRequest("", "?id=", nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()

		sut.Handle(w, r, "")
		res := w.Result()

		if err = assert.Equal(wantStatusCode, res.StatusCode); err != nil {
			t.Error(err)
		}

		assert.OnResErr(wantErrMsg)(t, res, "")
	})

	t.Run("IDNotInt", func(t *testing.T) {
		idValidator.Err = api.ErrStrNotInt

		wantStatusCode := http.StatusBadRequest
		wantErrMsg := "Subtask ID must be an integer."

		r, err := http.NewRequest("", "?id=", nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()

		sut.Handle(w, r, "")
		res := w.Result()

		if err = assert.Equal(wantStatusCode, res.StatusCode); err != nil {
			t.Error(err)
		}

		assert.OnResErr(wantErrMsg)(t, res, "")
	})

	t.Run("IDUnexpectedErr", func(t *testing.T) {
		idValidator.Err = api.ErrStrTooLong

		wantStatusCode := http.StatusInternalServerError
		wantLoggedErr := api.ErrStrTooLong.Error()

		r, err := http.NewRequest("", "?id=", nil)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()

		sut.Handle(w, r, "")
		res := w.Result()

		if err = assert.Equal(wantStatusCode, res.StatusCode); err != nil {
			t.Error(err)
		}

		assert.OnLoggedErr(wantLoggedErr)(t, res, log.InMessage)
	})
}
