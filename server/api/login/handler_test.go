package login

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/db"
	"server/token"
)

func TestHandler(t *testing.T) {
	var (
		userReader     = &db.FakeUserReader{}
		hashComparer   = &fakeHashComparer{}
		tokenGenerator = &token.FakeGenerator{}
	)
	sut := NewHandler(userReader, hashComparer, tokenGenerator)

	for _, c := range []struct {
		name                 string
		httpMethod           string
		reqBody              ReqBody
		userReaderOutRes     db.User
		userReaderOutErr     error
		hashComparerOutRes   bool
		hashComparerOutErr   error
		tokenGeneratorOutRes string
		tokenGeneratorOutErr error
		wantStatusCode       int
	}{
		{
			name:                 "MethodNotAllowed",
			httpMethod:           http.MethodGet,
			reqBody:              ReqBody{},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutRes:   false,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusMethodNotAllowed,
		},
		{
			name:                 "NoUsername",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutRes:   false,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "UsernameEmpty",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: ""},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutRes:   false,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "UserNotFound",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob21"},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     sql.ErrNoRows,
			hashComparerOutRes:   false,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "UserReaderError",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     errors.New("user reader error"),
			hashComparerOutRes:   false,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
		},
		{
			name:                 "NoPassword",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob21"},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutRes:   false,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "PasswordEmpty",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob21", Password: ""},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutRes:   false,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "WrongPassword",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutRes:   false,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "HashComparerError",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutRes:   true,
			hashComparerOutErr:   errors.New("hash comparer error"),
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
		},
		{
			name:                 "TokenGeneratorError",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutRes:   true,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: errors.New("token generator error"),
			wantStatusCode:       http.StatusInternalServerError,
		},
		{
			name:                 "Success",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutRes:   true,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusOK,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userReader.OutRes = c.userReaderOutRes
			userReader.OutErr = c.userReaderOutErr
			hashComparer.outRes = c.hashComparerOutRes
			hashComparer.outErr = c.hashComparerOutErr
			tokenGenerator.OutRes = c.tokenGeneratorOutRes
			tokenGenerator.OutErr = c.tokenGeneratorOutErr

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
			if c.wantStatusCode == http.StatusOK {
				// no field errors - assert on authToken
				foundSessionToken := false
				for _, cookie := range w.Result().Cookies() {
					if cookie.Name == "authToken" {
						foundSessionToken = true
						if err = assert.Equal(c.tokenGeneratorOutRes, cookie.Value); err != nil {
							t.Error(err)
						}
					}
				}
				if err = assert.True(foundSessionToken); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
