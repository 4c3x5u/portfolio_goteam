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
	subtaskTable "github.com/kxplxn/goteam/server/dbaccess/subtask"
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
				subtaskTable.NewSelector(db),
				pkgLog.New(),
			),
		},
	)

	for _, c := range []struct {
		name           string
		id             string
		wantStatusCode int
		wantErrMsg     string
	}{
		{
			name:           "IDEmpty",
			id:             "",
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     "Subtask ID cannot be empty.",
		},
		{
			name:           "IDNotInt",
			id:             "A",
			wantStatusCode: http.StatusBadRequest,
			wantErrMsg:     "Subtask ID must be an integer.",
		},
		{
			name:           "SubtaskNotFound",
			id:             "1001",
			wantStatusCode: http.StatusNotFound,
			wantErrMsg:     "Subtask not found.",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			r, err := http.NewRequest(http.MethodPatch, "?id="+c.id, nil)
			if err != nil {
				t.Fatal(err)
			}
			addBearerAuth(jwtBob123)(r)
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, r)
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
