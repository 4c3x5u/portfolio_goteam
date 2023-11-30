//go:build utest

package login

import (
	"database/sql"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/auth"
	userTable "github.com/kxplxn/goteam/pkg/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/pkg/log"

	"golang.org/x/crypto/bcrypt"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly.
func TestPOSTHandler(t *testing.T) {
	var (
		validator          = &fakeReqValidator{}
		userSelector       = &userTable.FakeSelector{}
		passwordComparer   = &fakeHashComparer{}
		authTokenGenerator = &auth.FakeTokenGenerator{}
		log                = &pkgLog.FakeErrorer{}
	)
	sut := NewPOSTHandler(
		validator, userSelector, passwordComparer, authTokenGenerator, log,
	)

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
						assert.Equal(t.Error,
							ck.Value, "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
						)
						assert.True(t.Error,
							ck.Expires.Unix() > time.Now().Unix(),
						)
						assert.True(t.Error, ck.Secure)
						assert.Equal(t.Error,
							ck.SameSite, http.SameSiteNoneMode,
						)
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

			req := httptest.NewRequest("", "/", strings.NewReader("{}"))
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.Handle(w, req, "")
			res := w.Result()

			// Assert on the status code.
			assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

			// Run case-specific assertions.
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
