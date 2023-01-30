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
		dbCloser           = &db.FakeCloser{}
		logger             = &log.FakeLogger{}
	)
	sut := NewHandler(
		validator, dbUserSelector, passwordComparer, authTokenGenerator,
		dbCloser, logger,
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
			validator.outOK = c.validatorOutOK
			dbUserSelector.OutRes = c.userSelectorOutUser
			dbUserSelector.OutErr = c.userSelectorOutErr
			passwordComparer.outErr = c.hashComparerOutErr
			authTokenGenerator.OutRes = c.tokenGeneratorOutToken
			authTokenGenerator.OutErr = c.tokenGeneratorOutErr

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

			sut.ServeHTTP(w, req)

			if err = assert.Equal(
				c.wantStatusCode, w.Result().StatusCode,
			); err != nil {
				t.Error(err)
			}

			// If 200 was expected, auth token must be set.
			if c.wantStatusCode == http.StatusOK {
				authTokenFound := false
				for _, ck := range w.Result().Cookies() {
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
				if err = assert.True(authTokenFound); err != nil {
					t.Error(err)
				}
			}

			// If 500 was expected, logger must be called.
			if c.wantStatusCode == http.StatusInternalServerError {
				errFound := false
				for _, err := range []error{
					c.userSelectorOutErr,
					c.hashComparerOutErr,
					c.tokenGeneratorOutErr,
				} {
					if err != nil {
						errFound = true

						if err := assert.Equal(
							log.LevelError, logger.InLevel,
						); err != nil {
							t.Error(err)
						}

						if err := assert.Equal(
							err.Error(), logger.InMessage,
						); err != nil {
							t.Error(err)
						}
					}
				}
				if !errFound {
					t.Errorf(
						"c.wantStatusCode was %d but no errors were logged.",
						http.StatusInternalServerError,
					)
				}
				return
			}

			// DEPENDENCY-INPUT-BASED ASSERTIONS

			// If request and HTTP method was valid, user selector and db closer
			// must be called.
			if !c.validatorOutOK ||
				c.wantStatusCode == http.StatusMethodNotAllowed {
				return
			}
			if err = assert.Equal(
				reqBody.Username, dbUserSelector.InUserID,
			); err != nil {
				t.Error(err)
			}
			if err = assert.True(dbCloser.IsCalled); err != nil {
				t.Error(err)
			}

			// If no user selector error is expected, hash comparer must be
			// called.
			if dbUserSelector.OutErr != nil {
				return
			}
			if err = assert.Equal(
				string(c.userSelectorOutUser.Password),
				string(passwordComparer.inHash),
			); err != nil {
				t.Error(err)
			}
			if err = assert.Equal(
				reqBody.Password, passwordComparer.inPlaintext,
			); err != nil {
				t.Error(err)
			}

			// If no hash comparer error is expected, token generator must be
			// called.
			if passwordComparer.outErr != nil {
				return
			}
			if err = assert.Equal(
				reqBody.Username, authTokenGenerator.InSub,
			); err != nil {
				t.Error(err)
			}
		})
	}
}
