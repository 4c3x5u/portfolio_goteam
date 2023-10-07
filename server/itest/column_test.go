//go:build itest

package itest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"server/api/board"
	columnTable "server/dbaccess/column"
	userboardTable "server/dbaccess/userboard"
	"testing"

	columnAPI "server/api/column"
	"server/assert"
	"server/auth"
	pkgLog "server/log"
)

func TestColumn(t *testing.T) {
	// Create board API handler.
	log := pkgLog.New()
	sut := columnAPI.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(jwtKey),
		columnAPI.NewIDValidator(),
		columnTable.NewSelector(db),
		userboardTable.NewSelector(db),
		columnTable.NewUpdater(db),
		log,
	)

	// Used in status 400 error cases to assert on the error message.
	assertOnErrMsg := func(
		wantErrMsg string,
	) func(*testing.T, *httptest.ResponseRecorder) {
		return func(t *testing.T, w *httptest.ResponseRecorder) {
			resBody := board.ResBody{}
			if err := json.NewDecoder(w.Result().Body).Decode(
				&resBody,
			); err != nil {
				t.Error(err)
			}
			if err := assert.Equal(
				wantErrMsg, resBody.Error,
			); err != nil {
				t.Error(err)
			}
		}
	}

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
			assertFunc func(*testing.T, *httptest.ResponseRecorder)
		}{
			{
				name:       "IDEmpty",
				id:         "",
				reqBody:    columnAPI.ReqBody{{ID: 0, Order: 0}},
				authFunc:   addBearerAuth(jwtBob123),
				statusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg("Column ID cannot be empty."),
			},
			{
				name:       "IDNotInt",
				id:         "A",
				reqBody:    columnAPI.ReqBody{{ID: 0, Order: 0}},
				authFunc:   addBearerAuth(jwtBob123),
				statusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg("Column ID must be an integer."),
			},
			{
				name:       "ColumnNotFound",
				id:         "1001",
				reqBody:    columnAPI.ReqBody{{ID: 0, Order: 0}},
				authFunc:   addBearerAuth(jwtBob123),
				statusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg("Column not found."),
			},
			{
				name:       "NoAccess",
				id:         "5",
				reqBody:    columnAPI.ReqBody{{ID: 0, Order: 0}},
				authFunc:   addBearerAuth(jwtBob124),
				statusCode: http.StatusUnauthorized,
				assertFunc: assertOnErrMsg(
					"You do not have access to this board.",
				),
			},
			{
				name:       "NotAdmin",
				id:         "6",
				reqBody:    columnAPI.ReqBody{{ID: 0, Order: 0}},
				authFunc:   addBearerAuth(jwtBob123),
				statusCode: http.StatusUnauthorized,
				assertFunc: assertOnErrMsg("Only board admins can move tasks."),
			},
			{
				name:       "TaskNotFound",
				id:         "5",
				reqBody:    columnAPI.ReqBody{{ID: 0, Order: 0}},
				authFunc:   addBearerAuth(jwtBob123),
				statusCode: http.StatusNotFound,
				assertFunc: assertOnErrMsg("Task not found."),
			},
			{
				name:       "Success",
				id:         "7",
				reqBody:    columnAPI.ReqBody{{ID: 5, Order: 2}},
				authFunc:   addBearerAuth(jwtBob123),
				statusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *httptest.ResponseRecorder) {
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

				c.assertFunc(t, w)
			})
		}
	})
}
