package board

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/auth"
	"server/db"
)

func TestHandler(t *testing.T) {
	tokenValidator := &auth.FakeTokenValidator{}
	userBoardCounter := &db.FakeCounter{}
	sut := NewHandler(tokenValidator, userBoardCounter)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		for _, httpMethod := range []string{
			http.MethodConnect, http.MethodDelete, http.MethodGet,
			http.MethodHead, http.MethodOptions, http.MethodPatch,
			http.MethodPut, http.MethodTrace,
		} {
			t.Run(httpMethod, func(t *testing.T) {
				req, err := http.NewRequest(httpMethod, "/board", nil)
				if err != nil {
					t.Fatal(err)
				}
				w := httptest.NewRecorder()

				sut.ServeHTTP(w, req)

				if err = assert.Equal(http.StatusMethodNotAllowed, w.Result().StatusCode); err != nil {
					t.Error(err)
				}
			})
		}
	})

	authCookie := &http.Cookie{Name: auth.CookieName}
	for _, c := range []struct {
		name                   string
		cookie                 *http.Cookie
		tokenValidatorOutErr   error
		userBoardCounterOutRes int
		wantStatusCode         int
		wantErr                string
	}{
		{
			name:                   "NoAuthCookie",
			cookie:                 nil,
			tokenValidatorOutErr:   nil,
			userBoardCounterOutRes: 3,
			wantStatusCode:         http.StatusUnauthorized,
			wantErr:                "",
		},
		{
			name:                   "InvalidAuthCookie",
			cookie:                 authCookie,
			tokenValidatorOutErr:   errors.New("token validator error"),
			userBoardCounterOutRes: 3,
			wantStatusCode:         http.StatusUnauthorized,
			wantErr:                "",
		},
		{
			name:                   "MaxBoardsCreated",
			cookie:                 authCookie,
			tokenValidatorOutErr:   nil,
			userBoardCounterOutRes: 3,
			wantStatusCode:         http.StatusBadRequest,
			wantErr:                errMaxBoards,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			tokenValidator.OutErr = c.tokenValidatorOutErr
			userBoardCounter.OutRes = c.userBoardCounterOutRes
			req, err := http.NewRequest(http.MethodPost, "/board", nil)
			if err != nil {
				t.Fatal(err)
			}
			if c.cookie != nil {
				req.AddCookie(c.cookie)
			}
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, req)

			if err = assert.Equal(c.wantStatusCode, w.Result().StatusCode); err != nil {
				t.Error(err)
			}
			if c.wantErr != "" {
				resBody := ResBody{}
				if err := json.NewDecoder(w.Result().Body).Decode(&resBody); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(c.wantErr, resBody.Error); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
