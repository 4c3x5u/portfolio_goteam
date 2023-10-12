//go:build itest

package itest

import (
	"bytes"
	"encoding/json"
	"github.com/kxplxn/goteam/server/api"
	"net/http"
	"net/http/httptest"
	"testing"

	columnAPI "github.com/kxplxn/goteam/server/api/column"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/auth"
	columnTable "github.com/kxplxn/goteam/server/dbaccess/column"
	userboardTable "github.com/kxplxn/goteam/server/dbaccess/userboard"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestColumnHandler tests the http.Handler for the column API route and asserts
// that it behaves correctly during various execution paths.
func TestColumnHandler(t *testing.T) {
	// Create board API handler.
	log := pkgLog.New()
	sut := api.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodPatch: columnAPI.NewPATCHHandler(
				columnAPI.NewIDValidator(),
				columnTable.NewSelector(db),
				userboardTable.NewSelector(db),
				columnTable.NewUpdater(db),
				log,
			),
		},
	)

	t.Run("Auth", func(t *testing.T) {
		for _, c := range []struct {
			name     string
			authFunc func(*http.Request)
		}{
			// Auth Cases
			{name: "HeaderEmpty", authFunc: func(*http.Request) {}},
			{name: "HeaderInvalid", authFunc: addBearerAuth("asdfasldfkjasd")},
		} {
			t.Run(c.name, func(t *testing.T) {
				t.Run(http.MethodPatch, func(t *testing.T) {
					req, err := http.NewRequest(http.MethodPatch, "", nil)
					if err != nil {
						t.Fatal(err)
					}
					c.authFunc(req)
					w := httptest.NewRecorder()

					sut.ServeHTTP(w, req)
					res := w.Result()

					if err = assert.Equal(
						http.StatusUnauthorized, res.StatusCode,
					); err != nil {
						t.Error(err)
					}

					if err = assert.Equal(
						"Bearer", res.Header.Values("WWW-Authenticate")[0],
					); err != nil {
						t.Error(err)
					}
				})
			})
		}
	})
	t.Run(http.MethodPatch, func(t *testing.T) {
		for _, c := range []struct {
			name       string
			id         string
			reqBody    columnAPI.ReqBody
			authFunc   func(*http.Request)
			statusCode int
			assertFunc func(*testing.T, *http.Response, string)
		}{
			{
				name:       "IDEmpty",
				id:         "",
				reqBody:    columnAPI.ReqBody{{ID: 0, Order: 0}},
				authFunc:   addBearerAuth(jwtBob123),
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Column ID cannot be empty."),
			},
			{
				name:       "IDNotInt",
				id:         "A",
				reqBody:    columnAPI.ReqBody{{ID: 0, Order: 0}},
				authFunc:   addBearerAuth(jwtBob123),
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Column ID must be an integer."),
			},
			{
				name:       "ColumnNotFound",
				id:         "1001",
				reqBody:    columnAPI.ReqBody{{ID: 0, Order: 0}},
				authFunc:   addBearerAuth(jwtBob123),
				statusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr("Column not found."),
			},
			{
				name:       "NoAccess",
				id:         "5",
				reqBody:    columnAPI.ReqBody{{ID: 0, Order: 0}},
				authFunc:   addBearerAuth(jwtBob124),
				statusCode: http.StatusUnauthorized,
				assertFunc: assert.OnResErr(
					"You do not have access to this board.",
				),
			},
			{
				name:       "NotAdmin",
				id:         "6",
				reqBody:    columnAPI.ReqBody{{ID: 0, Order: 0}},
				authFunc:   addBearerAuth(jwtBob123),
				statusCode: http.StatusUnauthorized,
				assertFunc: assert.OnResErr("Only board admins can move tasks."),
			},
			{
				name:       "TaskNotFound",
				id:         "5",
				reqBody:    columnAPI.ReqBody{{ID: 0, Order: 0}},
				authFunc:   addBearerAuth(jwtBob123),
				statusCode: http.StatusNotFound,
				assertFunc: assert.OnResErr("Task not found."),
			},
			{
				name:       "Success",
				id:         "7",
				reqBody:    columnAPI.ReqBody{{ID: 5, Order: 2}},
				authFunc:   addBearerAuth(jwtBob123),
				statusCode: http.StatusOK,
				assertFunc: func(
					t *testing.T, _ *http.Response, _ string,
				) {
					var columnID, order int
					if err := db.QueryRow(
						`SELECT columnID, "order" FROM app.task WHERE id = $1`,
						5,
					).Scan(&columnID, &order); err != nil {
						t.Fatal(err)
					}
					if err := assert.Equal(7, columnID); err != nil {
						t.Error(err)
					}
					if err := assert.Equal(2, order); err != nil {
						t.Error(err)
					}
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				tasks, err := json.Marshal(c.reqBody)
				if err != nil {
					t.Fatal(err)
				}
				req, err := http.NewRequest(
					http.MethodPatch, "?id="+c.id, bytes.NewReader(tasks),
				)
				if err != nil {
					t.Fatal(err)
				}
				c.authFunc(req)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				if err = assert.Equal(
					c.statusCode, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				c.assertFunc(t, res, "")
			})
		}
	})
}
