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
	"server/auth"
	"server/db"

	"golang.org/x/crypto/bcrypt"
)

// TestHandler tests Handler to ensure that it behaves correctly on all
// possible scenarios.
func TestHandler(t *testing.T) {
	var (
		userReader     = &db.FakeUserReader{}
		hashComparer   = &fakeHashComparer{}
		tokenGenerator = &auth.FakeTokenGenerator{}
		dbCloser       = &db.FakeCloser{}
	)
	sut := NewHandler(userReader, hashComparer, tokenGenerator, dbCloser)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		for _, httpMethod := range []string{
			http.MethodConnect, http.MethodDelete, http.MethodGet,
			http.MethodHead, http.MethodOptions, http.MethodPatch,
			http.MethodPut, http.MethodTrace,
		} {
			t.Run(httpMethod, func(t *testing.T) {
				req, err := http.NewRequest(httpMethod, "/login", nil)
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

	for _, c := range []struct {
		name                 string
		reqBody              ReqBody
		userReaderOutRes     db.User
		userReaderOutErr     error
		hashComparerOutErr   error
		tokenGeneratorOutRes string
		tokenGeneratorOutErr error
		wantStatusCode       int
	}{
		{
			name:                 "NoUsername",
			reqBody:              ReqBody{},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "UsernameEmpty",
			reqBody:              ReqBody{Username: ""},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "UserNotFound",
			reqBody:              ReqBody{Username: "bob21"},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     sql.ErrNoRows,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "UserReaderError",
			reqBody:              ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     errors.New("user reader error"),
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
		},
		{
			name:                 "NoPassword",
			reqBody:              ReqBody{Username: "bob21"},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "PasswordEmpty",
			reqBody:              ReqBody{Username: "bob21", Password: ""},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "WrongPassword",
			reqBody:              ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:     db.User{Username: "bob21", Password: []byte("$2a$ASasdflak$kajdsfh")},
			userReaderOutErr:     nil,
			hashComparerOutErr:   bcrypt.ErrMismatchedHashAndPassword,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
		},
		{
			name:                 "HashComparerError",
			reqBody:              ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:     db.User{Username: "bob21", Password: []byte("$2a$ASasdflak$kajdsfh")},
			userReaderOutErr:     nil,
			hashComparerOutErr:   errors.New("hash comparer error"),
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
		},
		{
			name:                 "TokenGeneratorError",
			reqBody:              ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:     db.User{Username: "bob21", Password: []byte("$2a$ASasdflak$kajdsfh")},
			userReaderOutErr:     nil,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: errors.New("token generator error"),
			wantStatusCode:       http.StatusInternalServerError,
		},
		{
			name:                 "Success",
			reqBody:              ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			userReaderOutRes:     db.User{Username: "bob21", Password: []byte("$2a$ASasdflak$kajdsfh")},
			userReaderOutErr:     nil,
			hashComparerOutErr:   nil,
			tokenGeneratorOutRes: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusOK,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userReader.OutRes = c.userReaderOutRes
			userReader.OutErr = c.userReaderOutErr
			hashComparer.outErr = c.hashComparerOutErr
			tokenGenerator.OutRes = c.tokenGeneratorOutRes
			tokenGenerator.OutErr = c.tokenGeneratorOutErr

			reqBodyJSON, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewReader(reqBodyJSON))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, req)

			if err = assert.Equal(c.wantStatusCode, w.Result().StatusCode); err != nil {
				t.Error(err)
			}

			// If 200 is expected - auth token must be set.
			if c.wantStatusCode == http.StatusOK {
				authTokenFound := false
				for _, ck := range w.Result().Cookies() {
					if ck.Name == auth.CookieName {
						authTokenFound = true
						if err = assert.Equal(c.tokenGeneratorOutRes, ck.Value); err != nil {
							t.Error(err)
						}
						if err = assert.True(ck.Expires.Unix() > time.Now().Unix()); err != nil {
							t.Error(err)
						}
					}
				}
				if err = assert.True(authTokenFound); err != nil {
					t.Error(err)
				}
			}

			// DEPENDENCY-INPUT-BASED ASSERTIONS

			// If username and password weren't empty, user reader and db closer must be called.
			if c.reqBody.Username == "" ||
				c.reqBody.Password == "" ||
				c.wantStatusCode == http.StatusMethodNotAllowed {
				return
			}
			if err = assert.Equal(c.reqBody.Username, userReader.InArg); err != nil {
				t.Error(err)
			}
			if err = assert.True(dbCloser.IsCalled); err != nil {
				t.Error(err)
			}

			// If no user reader error is expected, hash comparer must be called.
			if userReader.OutErr != nil {
				return
			}
			if err = assert.Equal(
				string(c.userReaderOutRes.Password), string(hashComparer.inHash),
			); err != nil {
				t.Error(err)
			}
			if err = assert.Equal(c.reqBody.Password, hashComparer.inPlaintext); err != nil {
				t.Error(err)
			}

			// If no hash comparer error is expected, token generator must be called.
			if hashComparer.outErr != nil {
				return
			}
			if err = assert.Equal(c.reqBody.Username, tokenGenerator.InSub); err != nil {
				t.Error(err)
			}

		})
	}
}
