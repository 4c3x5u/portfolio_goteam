//go:build utest

package loginapi

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"golang.org/x/crypto/bcrypt"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/usertbl"
	"github.com/kxplxn/goteam/pkg/log"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it behaves
// correctly.
func TestPOSTHandler(t *testing.T) {
	var (
		validator        = &fakeReqValidator{}
		userRetriever    = &db.FakeRetriever[usertbl.User]{}
		passwordComparer = &fakeHashComparer{}
		authEncoder      = &cookie.FakeEncoder[cookie.Auth]{}
		log              = &log.FakeErrorer{}
	)
	sut := NewPostHandler(
		validator, userRetriever, passwordComparer, authEncoder, log,
	)

	for _, c := range []struct {
		name             string
		reqIsValid       bool
		user             usertbl.User
		errRetrieveUser  error
		errCompareHash   error
		authToken        http.Cookie
		errGenerateToken error
		wantStatus       int
		assertFunc       func(*testing.T, *http.Response, []any)
	}{
		{
			name:             "InvalidRequest",
			reqIsValid:       false,
			user:             usertbl.User{},
			errRetrieveUser:  nil,
			errCompareHash:   nil,
			authToken:        http.Cookie{},
			errGenerateToken: nil,
			wantStatus:       http.StatusBadRequest,
			assertFunc:       func(*testing.T, *http.Response, []any) {},
		},
		{
			name:             "UserNotFound",
			reqIsValid:       true,
			user:             usertbl.User{},
			errRetrieveUser:  db.ErrNoItem,
			errCompareHash:   nil,
			authToken:        http.Cookie{},
			errGenerateToken: nil,
			wantStatus:       http.StatusBadRequest,
			assertFunc:       func(*testing.T, *http.Response, []any) {},
		},
		{
			name:             "UserSelectorError",
			reqIsValid:       true,
			user:             usertbl.User{},
			errRetrieveUser:  errors.New("user selector error"),
			errCompareHash:   nil,
			authToken:        http.Cookie{},
			errGenerateToken: nil,
			wantStatus:       http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr("user selector error"),
		},
		{
			name:       "WrongPassword",
			reqIsValid: true,
			user: usertbl.User{
				Username: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			errRetrieveUser:  nil,
			errCompareHash:   bcrypt.ErrMismatchedHashAndPassword,
			authToken:        http.Cookie{},
			errGenerateToken: nil,
			wantStatus:       http.StatusBadRequest,
			assertFunc:       func(*testing.T, *http.Response, []any) {},
		},
		{
			name:       "HashComparerError",
			reqIsValid: true,
			user: usertbl.User{
				Username: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			errRetrieveUser:  nil,
			errCompareHash:   errors.New("hash comparer error"),
			authToken:        http.Cookie{},
			errGenerateToken: nil,
			wantStatus:       http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr("hash comparer error"),
		},
		{
			name:       "TokenGeneratorError",
			reqIsValid: true,
			user: usertbl.User{
				Username: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			errRetrieveUser:  nil,
			errCompareHash:   nil,
			authToken:        http.Cookie{},
			errGenerateToken: errors.New("token generator error"),
			wantStatus:       http.StatusInternalServerError,
			assertFunc:       assert.OnLoggedErr("token generator error"),
		},
		{
			name:       "Success",
			reqIsValid: true,
			user: usertbl.User{
				Username: "bob123", Password: []byte("$2a$ASasdflak$kajdsfh"),
			},
			errRetrieveUser:  nil,
			errCompareHash:   nil,
			authToken:        http.Cookie{Name: "foo", Value: "bar"},
			errGenerateToken: nil,
			wantStatus:       http.StatusOK,
			assertFunc: func(t *testing.T, r *http.Response, _ []any) {
				ck := r.Cookies()[0]
				assert.Equal(t.Error, ck.Name, "foo")
				assert.Equal(t.Error, ck.Value, "bar")
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for sut's dependencies.
			validator.isValid = c.reqIsValid
			userRetriever.Res = c.user
			userRetriever.Err = c.errRetrieveUser
			passwordComparer.err = c.errCompareHash
			authEncoder.Res = c.authToken
			authEncoder.Err = c.errGenerateToken

			req := httptest.NewRequest("", "/", strings.NewReader("{}"))
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.Handle(w, req, "")
			res := w.Result()

			// Assert on the status code.
			assert.Equal(t.Error, res.StatusCode, c.wantStatus)

			// Run case-specific assertions.
			c.assertFunc(t, res, log.Args)
		})
	}
}
