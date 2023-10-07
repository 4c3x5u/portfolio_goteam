//go:build itest

package itest

import (
	"net/http"
	"net/http/httptest"
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
}
