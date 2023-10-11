//go:build utest

package subtask

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/assert"
)

// TestPATCHHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly in all possible scenarios.
func TestPATCHHandler(t *testing.T) {
	sut := NewPATCHHandler()

	t.Run("IDEmpty", func(t *testing.T) {
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
}
