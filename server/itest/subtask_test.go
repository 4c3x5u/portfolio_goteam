//go:build itest

package itest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server/api"
	subtaskAPI "github.com/kxplxn/goteam/server/api/subtask"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/auth"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	subtaskTable "github.com/kxplxn/goteam/server/dbaccess/subtask"
	taskTable "github.com/kxplxn/goteam/server/dbaccess/task"
	userboardTable "github.com/kxplxn/goteam/server/dbaccess/userboard"
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
				taskTable.NewSelector(db),
				columnTable.NewSelector(db),
				userboardTable.NewSelector(db),
				subtaskTable.NewUpdater(db),
				pkgLog.New(),
			),
		},
	)

	for _, c := range []struct {
		name           string
		id             string
		wantStatusCode int
		assertFunc     func(*testing.T, *http.Response, string)
	}{
		{
			name:           "IDEmpty",
			id:             "",
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     assert.OnResErr("Subtask ID cannot be empty."),
		},
		{
			name:           "IDNotInt",
			id:             "A",
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     assert.OnResErr("Subtask ID must be an integer."),
		},
		{
			name:           "SubtaskNotFound",
			id:             "1001",
			wantStatusCode: http.StatusNotFound,
			assertFunc:     assert.OnResErr("Subtask not found."),
		},
		{
			name:           "NoAccess",
			id:             "6",
			wantStatusCode: http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"You do not have access to this board.",
			),
		},
		{
			name:           "NotAdmin",
			id:             "7",
			wantStatusCode: http.StatusForbidden,
			assertFunc: assert.OnResErr(
				"Only board admins can edit subtasks.",
			),
		},
		{
			name:           "Success",
			id:             "8",
			wantStatusCode: http.StatusOK,
			assertFunc: func(t *testing.T, _ *http.Response, _ string) {
				var isDone bool
				if err := db.QueryRow(
					"SELECT isDone FROM app.subtask WHERE id = 8",
				).Scan(&isDone); err != nil {
					t.Fatal(err)
				}
				if err := assert.True(isDone); err != nil {
					t.Error(err)
				}
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			reqBody, err := json.Marshal(map[string]any{"done": true})
			if err != nil {
				t.Fatal(err)
			}
			r, err := http.NewRequest(
				http.MethodPatch, "?id="+c.id, bytes.NewReader(reqBody),
			)
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

			c.assertFunc(t, res, "")
		})
	}
}
