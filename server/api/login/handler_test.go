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

	// Used in status 500 cases to assert on the logged error message.
	assertOnLoggedErr := func(
		wantErrMsg string,
	) func(*testing.T, *log.FakeLogger, *http.Response) {
		return func(t *testing.T, l *log.FakeLogger, _ *http.Response) {
			if err := assert.Equal(log.LevelError, l.InLevel); err != nil {
				t.Error(err)
			}
			if err := assert.Equal(wantErrMsg, l.InMessage); err != nil {
				t.Error(err)
			}
		}
	}

	t.Run("MethodNotAllowed", func(t *testing.T) {
		for _, httpMethod := range []string{
			http.MethodConnect, http.MethodDelete, http.MethodGet,
			http.MethodHead, http.MethodOptions, http.MethodPatch,
			http.MethodPut, http.MethodTrace,
		} {
			t.Run(httpMethod, func(t *testing.T) {
				req, err := http.NewRequest(httpMethod, "", nil)
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
		assertFunc             func(*testing.T, *log.FakeLogger, *http.Response)
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
			assertFunc:             func(*testing.T, *log.FakeLogger, *http.Response) {},
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
			assertFunc:             func(*testing.T, *log.FakeLogger, *http.Response) {},
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
			assertFunc:             assertOnLoggedErr("user selector error"),
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
			assertFunc:             func(*testing.T, *log.FakeLogger, *http.Response) {},
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
			assertFunc:             assertOnLoggedErr("hash comparer error"),
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
			assertFunc:             assertOnLoggedErr("token generator error"),
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
			assertFunc: func(
				t *testing.T, _ *log.FakeLogger, r *http.Response,
			) {
				authTokenFound := false
				for _, ck := range r.Cookies() {
					if ck.Name == "auth-token" {
						authTokenFound = true
						if err := assert.Equal(
							"eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...", ck.Value,
						); err != nil {
							t.Error(err)
						}
						if err := assert.True(
							ck.Expires.Unix() > time.Now().Unix(),
						); err != nil {
							t.Error(err)
						}
					}
				}
				if !authTokenFound {
					t.Errorf("auth token was not found")
				}
			},
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
			reqBody, err := json.Marshal(ReqBody{})
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodPost, "", bytes.NewReader(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.ServeHTTP(w, req)
			res := w.Result()

			// Assert on the status code.
			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			// Run case-specific assertions.
			c.assertFunc(t, logger, res)
		})
	}
}
