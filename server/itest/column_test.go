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

	const bob123AuthToken = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJ" +
		"ib2IxMjMifQ.Y8_6K50EHUEJlJf4X21fNCFhYWhVIqN3Tw1niz8XwZc"

	// used in various test cases to authenticate the request sent
	addBearerAuth := func(token string) func(*http.Request) {
		return func(req *http.Request) {
			req.Header.Add("Authorization", "Bearer "+token)
		}
	}

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
			authFunc   func(*http.Request)
			statusCode int
			assertFunc func(*testing.T, *httptest.ResponseRecorder)
		}{
			{
				name:       "IDEmpty",
				id:         "",
				authFunc:   addBearerAuth(bob123AuthToken),
				statusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg("Column ID cannot be empty."),
			},
			{
				name:       "IDNotInt",
				id:         "A",
				authFunc:   addBearerAuth(bob123AuthToken),
				statusCode: http.StatusBadRequest,
				assertFunc: assertOnErrMsg("Column ID must be an integer."),
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				tasks, err := json.Marshal([]map[string]int{
					{"id": 0, "order": 0},
				})
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
