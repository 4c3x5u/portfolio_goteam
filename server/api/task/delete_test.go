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

	t.Run("IDEmpty", func(t *testing.T) {
		idValidator.Err = api.ErrStrEmpty

		wantStatusCode := http.StatusBadRequest
		wantErrMsg := "Task ID cannot be empty."

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
		wantErrMsg := "Task ID must be an integer."

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
}
