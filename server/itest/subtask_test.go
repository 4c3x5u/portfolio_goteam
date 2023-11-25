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
	boardTable "github.com/kxplxn/goteam/server/dbaccess/board"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	subtaskTable "github.com/kxplxn/goteam/server/dbaccess/subtask"
	taskTable "github.com/kxplxn/goteam/server/dbaccess/task"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
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
				boardTable.NewSelector(db),
				userTable.NewSelector(db),
				subtaskTable.NewUpdater(db),
				pkgLog.New(),
			),
		},
	)

	t.Run("PATCH", func(t *testing.T) {
		for _, c := range []struct {
			name           string
			id             string
			authFunc       func(*http.Request)
			wantStatusCode int
			assertFunc     func(*testing.T, *http.Response, string)
		}{
			{
				name:           "IDEmpty",
				id:             "",
				authFunc:       addBearerAuth(jwtTeam1Admin),
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Subtask ID cannot be empty."),
			},
			{
				name:           "IDNotInt",
				id:             "A",
				authFunc:       addBearerAuth(jwtTeam1Admin),
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("Subtask ID must be an integer."),
			},
			{
				name:           "SubtaskNotFound",
				id:             "1001",
				authFunc:       addBearerAuth(jwtTeam1Admin),
				wantStatusCode: http.StatusNotFound,
				assertFunc:     assert.OnResErr("Subtask not found."),
			},
			{
				name:           "BoardWrongTeam",
				id:             "5",
				authFunc:       addBearerAuth(jwtTeam2Admin),
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"You do not have access to this board.",
				),
			},
			{
				name:           "NotAdmin",
				id:             "5",
				authFunc:       addBearerAuth(jwtTeam1Member),
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can edit subtasks.",
				),
			},
			{
				name:           "Success",
				id:             "5",
				authFunc:       addBearerAuth(jwtTeam1Admin),
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ string) {
					var isDone bool
					if err := db.QueryRow(
						"SELECT isDone FROM app.subtask WHERE id = 5",
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
				c.authFunc(r)
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

	})
}
