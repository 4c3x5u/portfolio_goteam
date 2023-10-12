//go:build itest

package itest

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	subtaskAPI "github.com/kxplxn/goteam/server/api/subtask"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/auth"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestSubtaskHandler tests the http.Handler for the subtask API route and
// asserts that it behaves correctly during various execution paths.
func TestSubtaskHandler(t *testing.T) {
	sut := api.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodPatch: subtaskAPI.NewPATCHHandler(
				subtaskAPI.NewIDValidator(),
				pkgLog.New(),
			),
		},
	)

	t.Run("IDEmpty", func(t *testing.T) {
		wantStatusCode := http.StatusBadRequest
		r, err := http.NewRequest(http.MethodPatch, "?id=", nil)
		if err != nil {
			t.Fatal(err)
		}
		addBearerAuth(jwtBob123)(r)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)
		res := w.Result()

		if err = assert.Equal(wantStatusCode, res.StatusCode); err != nil {
			t.Error(err)
		}

		assert.OnResErr("Subtask ID cannot be empty.")(t, res, "")
	})

	t.Run("IDNotInt", func(t *testing.T) {
		wantStatusCode := http.StatusBadRequest
		r, err := http.NewRequest(http.MethodPatch, "?id=A", nil)
		if err != nil {
			t.Fatal(err)
		}
		addBearerAuth(jwtBob123)(r)
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, r)
		res := w.Result()

		if err = assert.Equal(wantStatusCode, res.StatusCode); err != nil {
			t.Error(err)
		}

		assert.OnResErr("Subtask ID must be an integer.")(t, res, "")
	})
}
