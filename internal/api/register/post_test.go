//go:build utest

package register

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db"
	userTable "github.com/kxplxn/goteam/pkg/db/user"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it
// behaves correctly.
func TestHandler(t *testing.T) {
	var (
		userValidator = &fakeReqValidator{}
		hasher        = &fakeHasher{}
		decodeInvite  = &token.FakeDecode[token.Invite]{}
		userInserter  = &db.FakeInserter[userTable.User]{}
		encodeAuth    = &token.FakeEncode[token.Auth]{}
		log           = &pkgLog.FakeErrorer{}
	)
	sut := NewPostHandler(
		userValidator,
		decodeInvite.Func,
		hasher,
		userInserter,
		encodeAuth.Func,
		log,
	)

	// Used in status 400 cases to assert on validation errors.
	assertOnErrsValidate := func(
		wantValidationErrs ValidationErrs,
	) func(*testing.T, *http.Response, string) {
		return func(t *testing.T, r *http.Response, _ string) {
			resBody := &PostResp{}
			if err := json.NewDecoder(r.Body).Decode(&resBody); err != nil {
				t.Fatal(err)
			}

			assert.AllEqual(t.Error,
				resBody.ValidationErrs.Username, wantValidationErrs.Username,
			)

			assert.AllEqual(t.Error,
				resBody.ValidationErrs.Password, wantValidationErrs.Password,
			)
		}
	}

	assertOnResErr := func(
		wantErrMsg string,
	) func(*testing.T, *http.Response, string) {
		return func(t *testing.T, res *http.Response, _ string) {
			var resBody PostResp
			if err := json.NewDecoder(
				res.Body,
			).Decode(&resBody); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t.Error, resBody.Err, wantErrMsg)
		}
	}

	validReqBody := PostReq{Username: "bob123", Password: "Myp4ssword!"}
	for _, c := range []struct {
		name            string
		req             PostReq
		errValidate     ValidationErrs
		inviteToken     string
		inviteDecoded   token.Invite
		errDecodeInvite error
		pwdHash         []byte
		errHash         error
		errInsertUser   error
		authToken       string
		errEncodeAuth   error
		wantStatus      int
		assertFunc      func(*testing.T, *http.Response, string)
	}{
		{
			name: "ErrsValidate",
			req:  PostReq{},
			errValidate: ValidationErrs{
				Username: []string{idTooLong}, Password: []string{pwdNoDigit},
			},
			inviteToken:     "",
			inviteDecoded:   token.Invite{},
			errDecodeInvite: nil,
			pwdHash:         nil,
			errHash:         nil,
			errInsertUser:   nil,
			authToken:       "",
			errEncodeAuth:   nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc: assertOnErrsValidate(
				ValidationErrs{
					Username: []string{idTooLong},
					Password: []string{pwdNoDigit},
				},
			),
		},
		{
			name:            "ErrDecodeInvite",
			req:             PostReq{},
			errValidate:     ValidationErrs{},
			inviteToken:     "someinvitetoken",
			inviteDecoded:   token.Invite{},
			errDecodeInvite: errors.New("an error"),
			pwdHash:         nil,
			errHash:         nil,
			errInsertUser:   nil,
			authToken:       "",
			errEncodeAuth:   nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc:      assertOnResErr("Invalid invite token."),
		},
		{
			name:            "ErrUsnTaken",
			req:             PostReq{},
			errValidate:     ValidationErrs{},
			inviteToken:     "",
			inviteDecoded:   token.Invite{},
			errDecodeInvite: nil,
			pwdHash:         nil,
			errHash:         nil,
			errInsertUser:   db.ErrDupKey,
			authToken:       "",
			errEncodeAuth:   nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc: assertOnErrsValidate(
				ValidationErrs{
					Username: []string{"Username is already taken."},
				},
			),
		},
		{
			name:          "ErrHash",
			req:           validReqBody,
			errValidate:   ValidationErrs{},
			inviteToken:   "",
			inviteDecoded: token.Invite{},
			pwdHash:       nil,
			errHash:       errors.New("hasher error"),
			errInsertUser: nil,
			authToken:     "",
			errEncodeAuth: nil,
			wantStatus:    http.StatusInternalServerError,
			assertFunc:    assert.OnLoggedErr("hasher error"),
		},
		{
			name:          "ErrUsnTaken",
			req:           PostReq{},
			errValidate:   ValidationErrs{},
			inviteToken:   "",
			inviteDecoded: token.Invite{},
			pwdHash:       nil,
			errHash:       nil,
			errInsertUser: db.ErrDupKey,
			authToken:     "",
			errEncodeAuth: nil,
			wantStatus:    http.StatusBadRequest,
			assertFunc: assertOnErrsValidate(
				ValidationErrs{
					Username: []string{"Username is already taken."},
				},
			),
		},
		{
			name:          "ErrPutUser",
			req:           validReqBody,
			errValidate:   ValidationErrs{},
			inviteToken:   "",
			inviteDecoded: token.Invite{},
			errInsertUser: errors.New("failed to put user"),
			pwdHash:       nil,
			errHash:       nil,
			authToken:     "",
			errEncodeAuth: nil,
			wantStatus:    http.StatusInternalServerError,
			assertFunc:    assert.OnLoggedErr("failed to put user"),
		},
		{
			name:          "ErrEncodeAuth",
			req:           validReqBody,
			errValidate:   ValidationErrs{},
			inviteToken:   "",
			inviteDecoded: token.Invite{},
			pwdHash:       nil,
			errHash:       nil,
			errInsertUser: nil,
			authToken:     "",
			errEncodeAuth: errors.New("error encoding auth token"),
			wantStatus:    http.StatusInternalServerError,
			assertFunc: assertOnResErr(
				"You have been registered successfully but something went " +
					"wrong. Please log in using the credentials you " +
					"registered with.",
			),
		},
		{
			name: "Success",
			req:  validReqBody,
			inviteToken: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFtSUQiOi" +
				"J0ZWFtaWQifQ.1h_fmLJ1ip-Z6kJq9JXYDgGuWDPOcOf8abwCgKtHHcY",
			errValidate:   ValidationErrs{},
			errInsertUser: nil,
			pwdHash:       nil,
			errHash:       nil,
			authToken:     "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			errEncodeAuth: nil,
			wantStatus:    http.StatusOK,
			assertFunc: func(
				t *testing.T, r *http.Response, _ string,
			) {
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
			userValidator.validationErrs = c.errValidate
			decodeInvite.Decoded = c.inviteDecoded
			decodeInvite.Err = c.errDecodeInvite
			hasher.hash = c.pwdHash
			hasher.err = c.errHash
			userInserter.Err = c.errInsertUser
			encodeAuth.Encoded = c.authToken
			encodeAuth.Err = c.errEncodeAuth

			// Prepare request and response recorder.
			reqBody, err := json.Marshal(c.req)
			if err != nil {
				t.Fatal(err)
			}
			req := httptest.NewRequest(
				http.MethodPost,
				"/",
				bytes.NewReader(reqBody),
			)
			if c.inviteToken != "" {
				req.AddCookie(&http.Cookie{
					Name:     token.InviteName,
					Value:    c.inviteToken,
					SameSite: http.SameSiteNoneMode,
					Secure:   true,
				})
			}
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.Handle(w, req, "")
			res := w.Result()

			// Assert on the status code.
			assert.Equal(t.Error, res.StatusCode, c.wantStatus)

			// Run case-specific assertions
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
