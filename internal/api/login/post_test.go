//go:build utest

package login

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
	userTable "github.com/kxplxn/goteam/pkg/db/user"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"

	"golang.org/x/crypto/bcrypt"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly.
func TestPOSTHandler(t *testing.T) {
	var (
		validator        = &fakeReqValidator{}
		userGetter       = &userTable.FakeGetter{}
		passwordComparer = &fakeHashComparer{}
		encodeAuthToken  = &token.FakeEncodeAuth{}
		log              = &pkgLog.FakeErrorer{}
	)
	sut := NewPostHandler(
		validator, userGetter, passwordComparer, encodeAuthToken.Func, log,
	)

	// Used on cases where no case-specific assertions are required.
	assertNone := func(*testing.T, *http.Response, string) {}

	for _, c := range []struct {
		name             string
		reqIsValid       bool
		userRecord       userTable.User
		errGetUser       error
		errCompareHash   error
		authToken        string
		errGenerateToken error
		wantStatus       int
		assertFunc       func(*testing.T, *http.Response, string)
	}{
		{
			name:             "InvalidRequest",
			reqIsValid:       false,
			userRecord:       userTable.User{},
			errGetUser:       nil,
			errCompareHash:   nil,
			authToken:        "",
			errGenerateToken: nil,
			wantStatus:       http.StatusBadRequest,
			assertFunc:       assertNone,
		},
		{
			name:             "UserNotFound",
			reqIsValid:       true,
			userRecord:       userTable.User{},
			errGetUser:       db.ErrNoItem,
			errCompareHash:   nil,
			authToken:        "",
			errGenerateToken: nil,
			wantStatus:       http.StatusBadRequest,
			assertFunc:       assertNone,
		},
		{
			name:             "UserSelectorError",
			reqIsValid:       true,
			userRecord:       userTable.User{},
			errGetUser:       errors.New("user selector error"),
			errCompareHash:   nil,
			authToken:        "",
			errGenerateToken: nil,
			wantStatus:       http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr("user selector error"),
		},
		{
			name:       "WrongPassword",
			reqIsValid: true,
			userRecord: userTable.User{
				Username: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			errGetUser:       nil,
			errCompareHash:   bcrypt.ErrMismatchedHashAndPassword,
			authToken:        "",
			errGenerateToken: nil,
			wantStatus:       http.StatusBadRequest,
			assertFunc:       assertNone,
		},
		{
			name:       "HashComparerError",
			reqIsValid: true,
			userRecord: userTable.User{
				Username: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			errGetUser:       nil,
			errCompareHash:   errors.New("hash comparer error"),
			authToken:        "",
			errGenerateToken: nil,
			wantStatus:       http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr("hash comparer error"),
		},
		{
			name:       "TokenGeneratorError",
			reqIsValid: true,
			userRecord: userTable.User{
				Username: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			errGetUser:       nil,
			errCompareHash:   nil,
			authToken:        "",
			errGenerateToken: errors.New("token generator error"),
			wantStatus:       http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr("token generator error"),
		},
		{
			name:       "Success",
			reqIsValid: true,
			userRecord: userTable.User{
				Username: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			errGetUser:       nil,
			errCompareHash:   nil,
			authToken:        "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			errGenerateToken: nil,
			wantStatus:       http.StatusOK,
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
			userGetter.User = c.userRecord
			userGetter.Err = c.errGetUser
			passwordComparer.err = c.errCompareHash
			encodeAuthToken.Encoded = c.authToken
			encodeAuthToken.Err = c.errGenerateToken

			req := httptest.NewRequest("", "/", strings.NewReader("{}"))
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.Handle(w, req, "")
			res := w.Result()

			// Assert on the status code.
			assert.Equal(t.Error, res.StatusCode, c.wantStatus)

			// Run case-specific assertions.
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
