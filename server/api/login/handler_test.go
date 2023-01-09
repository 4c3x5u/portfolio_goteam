package login

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"server/assert"
	"server/cookie"
	"server/db"
)

func TestHandler(t *testing.T) {
	var (
		userReader      = &db.FakeUserReader{}
		hashComparer    = &fakeHashComparer{}
		cookieGenerator = &cookie.FakeAuthGenerator{}
		cookieExpiry    = time.Now().Add(1 * time.Hour).Truncate(1 * time.Second).UTC()
	)
	sut := NewHandler(userReader, hashComparer, cookieGenerator)

	for _, c := range []struct {
		name                  string
		httpMethod            string
		reqBody               ReqBody
		userReaderOutRes      db.User
		userReaderOutErr      error
		hashComparerOutRes    bool
		hashComparerOutErr    error
		cookieGeneratorOutRes *http.Cookie
		tokenGeneratorOutErr  error
		wantStatusCode        int
	}{
		{
			name:                  "MethodNotAllowed",
			httpMethod:            http.MethodGet,
			reqBody:               ReqBody{},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      nil,
			hashComparerOutRes:    false,
			hashComparerOutErr:    nil,
			cookieGeneratorOutRes: nil,
			tokenGeneratorOutErr:  nil,
			wantStatusCode:        http.StatusMethodNotAllowed,
		},
		{
			name:                  "NoUsername",
			httpMethod:            http.MethodPost,
			reqBody:               ReqBody{},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      nil,
			hashComparerOutRes:    false,
			hashComparerOutErr:    nil,
			cookieGeneratorOutRes: nil,
			tokenGeneratorOutErr:  nil,
			wantStatusCode:        http.StatusBadRequest,
		},
		{
			name:                  "UsernameEmpty",
			httpMethod:            http.MethodPost,
			reqBody:               ReqBody{Username: ""},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      nil,
			hashComparerOutRes:    false,
			hashComparerOutErr:    nil,
			cookieGeneratorOutRes: nil,
			tokenGeneratorOutErr:  nil,
			wantStatusCode:        http.StatusBadRequest,
		},
		{
			name:                  "UserNotFound",
			httpMethod:            http.MethodPost,
			reqBody:               ReqBody{Username: "bob21"},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      sql.ErrNoRows,
			hashComparerOutRes:    false,
			hashComparerOutErr:    nil,
			cookieGeneratorOutRes: nil,
			tokenGeneratorOutErr:  nil,
			wantStatusCode:        http.StatusBadRequest,
		},
		{
			name:                  "UserReaderError",
			httpMethod:            http.MethodPost,
			reqBody:               ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      errors.New("user reader error"),
			hashComparerOutRes:    false,
			hashComparerOutErr:    nil,
			cookieGeneratorOutRes: nil,
			tokenGeneratorOutErr:  nil,
			wantStatusCode:        http.StatusInternalServerError,
		},
		{
			name:                  "NoPassword",
			httpMethod:            http.MethodPost,
			reqBody:               ReqBody{Username: "bob21"},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      nil,
			hashComparerOutRes:    false,
			hashComparerOutErr:    nil,
			cookieGeneratorOutRes: nil,
			tokenGeneratorOutErr:  nil,
			wantStatusCode:        http.StatusBadRequest,
		},
		{
			name:                  "PasswordEmpty",
			httpMethod:            http.MethodPost,
			reqBody:               ReqBody{Username: "bob21", Password: ""},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      nil,
			hashComparerOutRes:    false,
			hashComparerOutErr:    nil,
			cookieGeneratorOutRes: nil,
			tokenGeneratorOutErr:  nil,
			wantStatusCode:        http.StatusBadRequest,
		},
		{
			name:                  "WrongPassword",
			httpMethod:            http.MethodPost,
			reqBody:               ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      nil,
			hashComparerOutRes:    false,
			hashComparerOutErr:    nil,
			cookieGeneratorOutRes: nil,
			tokenGeneratorOutErr:  nil,
			wantStatusCode:        http.StatusBadRequest,
		},
		{
			name:                  "HashComparerError",
			httpMethod:            http.MethodPost,
			reqBody:               ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      nil,
			hashComparerOutRes:    true,
			hashComparerOutErr:    errors.New("hash comparer error"),
			cookieGeneratorOutRes: nil,
			tokenGeneratorOutErr:  nil,
			wantStatusCode:        http.StatusInternalServerError,
		},
		{
			name:                  "TokenGeneratorError",
			httpMethod:            http.MethodPost,
			reqBody:               ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      nil,
			hashComparerOutRes:    true,
			hashComparerOutErr:    nil,
			cookieGeneratorOutRes: nil,
			tokenGeneratorOutErr:  errors.New("token generator error"),
			wantStatusCode:        http.StatusInternalServerError,
		},
		{
			name:               "Success",
			httpMethod:         http.MethodPost,
			reqBody:            ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:   db.User{},
			userReaderOutErr:   nil,
			hashComparerOutRes: true,
			hashComparerOutErr: nil,
			cookieGeneratorOutRes: &http.Cookie{
				Name:    "authToken",
				Value:   "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
				Expires: cookieExpiry,
			},
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusOK,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userReader.OutRes = c.userReaderOutRes
			userReader.OutErr = c.userReaderOutErr
			hashComparer.outRes = c.hashComparerOutRes
			hashComparer.outErr = c.hashComparerOutErr
			cookieGenerator.OutRes = c.cookieGeneratorOutRes
			cookieGenerator.OutErr = c.tokenGeneratorOutErr

			reqBodyJSON, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(c.httpMethod, "/login", bytes.NewReader(reqBodyJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, req)

			if err = assert.Equal(c.wantStatusCode, w.Result().StatusCode); err != nil {
				t.Error(err)
			}

			// If 200 expected - auth token must be set.
			if c.wantStatusCode != http.StatusOK {
				return
			}
			authTokenFound := false
			for _, ck := range w.Result().Cookies() {
				if ck.Name == "authToken" {
					authTokenFound = true
					if err = assert.Equal(c.cookieGeneratorOutRes.Value, ck.Value); err != nil {
						t.Error(err)
					}
					if err = assert.Equal(cookieExpiry, ck.Expires); err != nil {
						t.Error(err)
					}
				}
			}
			if err = assert.True(authTokenFound); err != nil {
				t.Error(err)
			}
		})
	}
}
