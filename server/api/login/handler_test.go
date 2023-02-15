//go:build utest

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
	"server/log"

	"golang.org/x/crypto/bcrypt"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly.
func TestHandler(t *testing.T) {
	var (
		validator          = &fakeReqValidator{}
		dbUserSelector     = &db.FakeUserSelector{}
		passwordComparer   = &fakeHashComparer{}
		authTokenGenerator = &auth.FakeTokenGenerator{}
		logger             = &log.FakeLogger{}
	)
	sut := NewHandler(
		validator, dbUserSelector, passwordComparer, authTokenGenerator, logger,
	)

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

				if err = assert.Equal(
					http.StatusMethodNotAllowed, w.Result().StatusCode,
				); err != nil {
					t.Error(err)
				}

				if err := assert.Equal(
					http.MethodPost,
					w.Result().Header.Get("Access-Control-Allow-Methods"),
				); err != nil {
					t.Error(err)
				}
			})
		}
	})

	for _, c := range []struct {
		name                   string
		validatorOutOK         bool
		userSelectorOutUser    db.User
		userSelectorOutErr     error
		hashComparerOutErr     error
		tokenGeneratorOutToken string
		tokenGeneratorOutErr   error
		wantStatusCode         int
		wantErr                string
	}{
		{
			name:                   "InvalidRequest",
			validatorOutOK:         false,
			userSelectorOutUser:    db.User{},
			userSelectorOutErr:     nil,
			hashComparerOutErr:     nil,
			tokenGeneratorOutToken: "",
			tokenGeneratorOutErr:   nil,
			wantStatusCode:         http.StatusBadRequest,
		},
		{
			name:                   "UserNotFound",
			validatorOutOK:         true,
			userSelectorOutUser:    db.User{},
			userSelectorOutErr:     sql.ErrNoRows,
			hashComparerOutErr:     nil,
			tokenGeneratorOutToken: "",
			tokenGeneratorOutErr:   nil,
			wantStatusCode:         http.StatusBadRequest,
		},
		{
			name:                   "UserSelectorError",
			validatorOutOK:         true,
			userSelectorOutUser:    db.User{},
			userSelectorOutErr:     errors.New("user selector error"),
			hashComparerOutErr:     nil,
			tokenGeneratorOutToken: "",
			tokenGeneratorOutErr:   nil,
			wantStatusCode:         http.StatusInternalServerError,
		},
		{
			name:           "WrongPassword",
			validatorOutOK: true,
			userSelectorOutUser: db.User{
				ID: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			userSelectorOutErr:     nil,
			hashComparerOutErr:     bcrypt.ErrMismatchedHashAndPassword,
			tokenGeneratorOutToken: "",
			tokenGeneratorOutErr:   nil,
			wantStatusCode:         http.StatusBadRequest,
		},
		{
			name:           "HashComparerError",
			validatorOutOK: true,
			userSelectorOutUser: db.User{
				ID: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			userSelectorOutErr:     nil,
			hashComparerOutErr:     errors.New("hash comparer error"),
			tokenGeneratorOutToken: "",
			tokenGeneratorOutErr:   nil,
			wantStatusCode:         http.StatusInternalServerError,
		},
		{
			name:           "TokenGeneratorError",
			validatorOutOK: true,
			userSelectorOutUser: db.User{
				ID: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			userSelectorOutErr:     nil,
			hashComparerOutErr:     nil,
			tokenGeneratorOutToken: "",
			tokenGeneratorOutErr:   errors.New("token generator error"),
			wantStatusCode:         http.StatusInternalServerError,
		},
		{
			name:           "Success",
			validatorOutOK: true,
			userSelectorOutUser: db.User{
				ID: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			userSelectorOutErr:     nil,
			hashComparerOutErr:     nil,
			tokenGeneratorOutToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			tokenGeneratorOutErr:   nil,
			wantStatusCode:         http.StatusOK,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for sut's dependencies.
			validator.outOK = c.validatorOutOK
			dbUserSelector.OutRes = c.userSelectorOutUser
			dbUserSelector.OutErr = c.userSelectorOutErr
			passwordComparer.outErr = c.hashComparerOutErr
			authTokenGenerator.OutRes = c.tokenGeneratorOutToken
			authTokenGenerator.OutErr = c.tokenGeneratorOutErr

			// Prepare request and response recorder.
			reqBody := ReqBody{
				Username: "bob123", Password: "Myp4ssword!",
			}
			reqBodyJSON, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPost, "/login", bytes.NewReader(reqBodyJSON),
			)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.ServeHTTP(w, req)
			res := w.Result()

			if err = assert.Equal(c.wantStatusCode, res.StatusCode); err != nil {
				t.Error(err)
			}

			switch c.wantStatusCode {
			case http.StatusOK:
				// 200 was expected - auth token must be set.
				authTokenFound := false
				for _, ck := range res.Cookies() {
					if ck.Name == auth.CookieName {
						authTokenFound = true
						if err = assert.Equal(
							c.tokenGeneratorOutToken, ck.Value,
						); err != nil {
							t.Error(err)
						}
						if err = assert.True(
							ck.Expires.Unix() > time.Now().Unix(),
						); err != nil {
							t.Error(err)
						}
					}
				}
				if !authTokenFound {
					t.Errorf("200 was expected but auth token was not set")
				}
			case http.StatusInternalServerError:
				// 500 was expected - an error must be logged.
				errFound := false
				for _, depErr := range []error{
					c.userSelectorOutErr,
					c.hashComparerOutErr,
					c.tokenGeneratorOutErr,
				} {
					if depErr != nil {
						errFound = true
						if err = assert.Equal(
							log.LevelError, logger.InLevel,
						); err != nil {
							t.Error(err)
						}
						if err = assert.Equal(
							depErr.Error(), logger.InMessage,
						); err != nil {
							t.Error(err)
						}
					}
				}
				if !errFound {
					t.Errorf(
						"500 was expected but no errors were returned " +
							"from sut's dependencies",
					)
				}
			}
		})
	}
}
