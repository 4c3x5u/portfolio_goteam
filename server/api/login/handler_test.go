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

	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/auth"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"

	"golang.org/x/crypto/bcrypt"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly.
func TestHandler(t *testing.T) {
	var (
		validator          = &fakeReqValidator{}
		userSelector       = &userTable.FakeSelector{}
		passwordComparer   = &fakeHashComparer{}
		authTokenGenerator = &auth.FakeTokenGenerator{}
		log                = &pkgLog.FakeErrorer{}
	)
	sut := NewHandler(
		validator, userSelector, passwordComparer, authTokenGenerator, log,
	)

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
				res := w.Result()

				if err = assert.Equal(
					http.StatusMethodNotAllowed, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				if err := assert.Equal(
					http.MethodPost,
					res.Header.Get("Access-Control-Allow-Methods"),
				); err != nil {
					t.Error(err)
				}
			})
		}
	})

	// Used on cases where no case-specific assertions are required.
	emptyAssertFunc := func(*testing.T, *http.Response, string) {}

	for _, c := range []struct {
		name              string
		reqIsValid        bool
		userRecord        userTable.Record
		userSelectorErr   error
		hashComparerErr   error
		authToken         string
		tokenGeneratorErr error
		wantStatusCode    int
		assertFunc        func(*testing.T, *http.Response, string)
	}{
		{
			name:              "InvalidRequest",
			reqIsValid:        false,
			userRecord:        userTable.Record{},
			userSelectorErr:   nil,
			hashComparerErr:   nil,
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusBadRequest,
			assertFunc:        emptyAssertFunc,
		},
		{
			name:              "UserNotFound",
			reqIsValid:        true,
			userRecord:        userTable.Record{},
			userSelectorErr:   sql.ErrNoRows,
			hashComparerErr:   nil,
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusBadRequest,
			assertFunc:        emptyAssertFunc,
		},
		{
			name:              "UserSelectorError",
			reqIsValid:        true,
			userRecord:        userTable.Record{},
			userSelectorErr:   errors.New("user selector error"),
			hashComparerErr:   nil,
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr("user selector error"),
		},
		{
			name:       "WrongPassword",
			reqIsValid: true,
			userRecord: userTable.Record{
				Username: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			userSelectorErr:   nil,
			hashComparerErr:   bcrypt.ErrMismatchedHashAndPassword,
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusBadRequest,
			assertFunc:        emptyAssertFunc,
		},
		{
			name:       "HashComparerError",
			reqIsValid: true,
			userRecord: userTable.Record{
				Username: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			userSelectorErr:   nil,
			hashComparerErr:   errors.New("hash comparer error"),
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr("hash comparer error"),
		},
		{
			name:       "TokenGeneratorError",
			reqIsValid: true,
			userRecord: userTable.Record{
				Username: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			userSelectorErr:   nil,
			hashComparerErr:   nil,
			authToken:         "",
			tokenGeneratorErr: errors.New("token generator error"),
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr("token generator error"),
		},
		{
			name:       "Success",
			reqIsValid: true,
			userRecord: userTable.Record{
				Username: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			userSelectorErr:   nil,
			hashComparerErr:   nil,
			authToken:         "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusOK,
			assertFunc: func(t *testing.T, r *http.Response, _ string) {
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
						if err := assert.True(ck.Secure); err != nil {
							t.Error(err)
						}
						if err := assert.Equal(
							http.SameSiteNoneMode, ck.SameSite,
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
			validator.isValid = c.reqIsValid
			userSelector.Rec = c.userRecord
			userSelector.Err = c.userSelectorErr
			passwordComparer.err = c.hashComparerErr
			authTokenGenerator.AuthToken = c.authToken
			authTokenGenerator.Err = c.tokenGeneratorErr

			// Prepare request and response recorder.
			reqBody, err := json.Marshal(ReqBody{})
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
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
