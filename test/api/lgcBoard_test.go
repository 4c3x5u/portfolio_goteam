//go:build itest

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	lgcBoardAPI "github.com/kxplxn/goteam/internal/api/board"
	"github.com/kxplxn/goteam/pkg/api"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/auth"
	lgcBoardTable "github.com/kxplxn/goteam/pkg/legacydb/board"
	lgcUserTable "github.com/kxplxn/goteam/pkg/legacydb/user"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
)

func TestLegacyBoardAPI(t *testing.T) {
	userSelector := lgcUserTable.NewSelector(db)
	boardInserter := lgcBoardTable.NewInserter(db)
	nameValidator := lgcBoardAPI.NewNameValidator()
	log := pkgLog.New()
	sut := api.NewHandler(
		auth.NewJWTValidator(jwtKey),
		map[string]api.MethodHandler{
			http.MethodPost: lgcBoardAPI.NewPOSTHandler(
				userSelector,
				nameValidator,
				lgcBoardTable.NewCounter(db),
				boardInserter,
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
			{name: "HeaderInvalid", authFunc: addCookieAuth("asdfasldfkjasd")},
		} {
			t.Run(c.name, func(t *testing.T) {
				t.Run("POST", func(t *testing.T) {
					req := httptest.NewRequest(http.MethodPost, "/board", nil)
					c.authFunc(req)
					w := httptest.NewRecorder()

					sut.ServeHTTP(w, req)
					res := w.Result()

					assert.Equal(t.Error,
						res.StatusCode, http.StatusUnauthorized,
					)
					assert.Equal(t.Error,
						res.Header.Values("WWW-Authenticate")[0], "Bearer",
					)
				})
			})
		}
	})

	t.Run("POST", func(t *testing.T) {
		for _, c := range []struct {
			name           string
			authFunc       func(*http.Request)
			boardName      string
			wantStatusCode int
			assertFunc     func(*testing.T, *http.Response, string)
		}{
			{
				name:           "NotAdmin",
				authFunc:       addCookieAuth(jwtTeam2Member),
				boardName:      "Team 2 Board 2",
				wantStatusCode: http.StatusForbidden,
				assertFunc: assert.OnResErr(
					"Only team admins can create boards.",
				),
			},
			{
				name:           "EmptyBoardName",
				authFunc:       addCookieAuth(jwtTeam1Admin),
				boardName:      "",
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("input is empty"),
			},
			{
				name:           "TooLongBoardName",
				authFunc:       addCookieAuth(jwtTeam1Admin),
				boardName:      "A Board Whose Name Is Just Too Long!",
				wantStatusCode: http.StatusBadRequest,
				assertFunc:     assert.OnResErr("input is too long"),
			},
			{
				name:           "TooManyBoards",
				authFunc:       addCookieAuth(jwtTeam1Admin),
				boardName:      "bob123's new board",
				wantStatusCode: http.StatusBadRequest,
				assertFunc: assert.OnResErr(
					"You have already created the maximum amount of boards " +
						"allowed per user. Please delete one of your boards " +
						"to create a new one.",
				),
			},
			{
				name:           "Success",
				authFunc:       addCookieAuth(jwtTeam4Admin),
				boardName:      "Team 4 Board 1",
				wantStatusCode: http.StatusOK,
				assertFunc: func(t *testing.T, _ *http.Response, _ string) {
					// Assert that bob124 is assigned to the board as admin.
					var count int
					err := db.QueryRow(
						"SELECT COUNT(*) boardID FROM app.board " +
							"WHERE teamID = 4",
					).Scan(&count)
					if err != nil {
						t.Error(err)
					}
					assert.Equal(t.Error, count, 1)

					var boardID int
					err = db.QueryRow(
						"SELECT id FROM app.board WHERE teamID = 4",
					).Scan(&boardID)
					if err != nil {
						t.Error(err)
					}

					// Assert that 4 columns are created for this board.
					for order := 1; order < 5; order++ {
						var columnCount int
						err = db.QueryRow(
							`SELECT COUNT(*) FROM app."column" `+
								`WHERE boardID = $1 AND "order" = $2`,
							boardID,
							order,
						).Scan(&columnCount)
						if err != nil {
							t.Error(err)
						}
						assert.Equal(t.Error, columnCount, 1)
					}
				},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				reqBody, err := json.Marshal(map[string]string{
					"name": c.boardName,
				})
				if err != nil {
					t.Fatal(err)
				}
				req := httptest.NewRequest(
					http.MethodPost, "/", bytes.NewReader(reqBody),
				)
				c.authFunc(req)
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)
				res := w.Result()

				assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

				// Run case-specific assertions.
				c.assertFunc(t, res, "")
			})
		}
	})
}
