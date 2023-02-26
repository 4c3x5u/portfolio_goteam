//go:build itest

package itest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/api"
	boardAPI "server/api/board"
	"server/assert"
	"server/auth"
	"server/db"
	"server/log"
)

func TestBoard(t *testing.T) {
	// Create board API handler.
	logger := log.NewAppLogger()
	sut := boardAPI.NewHandler(
		auth.NewBearerTokenReader(),
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodPost: boardAPI.NewPOSTHandler(
				boardAPI.NewPOSTValidator(),
				db.NewUserBoardCounter(dbConnPool),
				db.NewBoardInserter(dbConnPool),
				logger,
			),
			http.MethodDelete: boardAPI.NewDELETEHandler(
				boardAPI.NewDELETEValidator(),
				db.NewUserBoardSelector(dbConnPool),
				db.NewBoardDeleter(dbConnPool),
				logger,
			),
		},
	)

	t.Run("NoAuthHeader", func(t *testing.T) {
		reqBody, err := json.Marshal(boardAPI.POSTReqBody{
			Name: "New Board",
		})
		if err != nil {
			t.Fatal(err)
		}
		req, err := http.NewRequest(
			http.MethodPost, "", bytes.NewReader(reqBody),
		)
		if err != nil {
			t.Fatal(err)
		}
		w := httptest.NewRecorder()

		sut.ServeHTTP(w, req)

		res := w.Result()

		if err = assert.Equal(
			http.StatusUnauthorized, res.StatusCode,
		); err != nil {
			t.Error(err)
		}

		wantAuthHeaderName, wantAuthHeaderValue := auth.WWWAuthenticate()
		gotAuthHeaderValue := res.Header.Values(wantAuthHeaderName)[0]
		if err := assert.Equal(
			wantAuthHeaderValue,
			gotAuthHeaderValue,
		); err != nil {
			t.Error(err)
		}
	})
}
