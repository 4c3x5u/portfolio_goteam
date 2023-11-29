//go:build utest

package register

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/kxplxn/goteam/server/api"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/auth"
	teamTable "github.com/kxplxn/goteam/server/dbaccess/team"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/server/log"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it
// behaves correctly.
func TestHandler(t *testing.T) {
	var (
		userValidator       = &fakeUserValidator{}
		inviteCodeValidator = &api.FakeStringValidator{}
		teamSelector        = &teamTable.FakeSelector{}
		userSelector        = &userTable.FakeSelector{}
		hasher              = &fakeHasher{}
		userInserter        = &userTable.FakeInserter{}
		tokenGenerator      = &auth.FakeTokenGenerator{}
		log                 = &pkgLog.FakeErrorer{}
	)
	sut := NewPOSTHandler(
		userValidator,
		inviteCodeValidator,
		teamSelector,
		userSelector,
		hasher,
		userInserter,
		tokenGenerator,
		log,
	)

	// Used in status 400 cases to assert on validation errors.
	assertOnValidationErrs := func(
		wantValidationErrs ValidationErrors,
	) func(*testing.T, *http.Response, string) {
		return func(t *testing.T, r *http.Response, _ string) {
			resBody := &POSTResp{}
			if err := json.NewDecoder(r.Body).Decode(&resBody); err != nil {
				t.Fatal(err)
			}

			if err := assert.EqualArr(
				wantValidationErrs.Username,
				resBody.ValidationErrs.Username,
			); err != nil {
				t.Error(err)
			}

			if err := assert.EqualArr(
				wantValidationErrs.Password,
				resBody.ValidationErrs.Password,
			); err != nil {
				t.Error(err)
			}
		}
	}

	assertOnResErr := func(
		wantErrMsg string,
	) func(*testing.T, *http.Response, string) {
		return func(t *testing.T, res *http.Response, _ string) {
			var resBody POSTResp
			if err := json.NewDecoder(
				res.Body,
			).Decode(&resBody); err != nil {
				t.Fatal(err)
			}
			if err := assert.Equal(
				wantErrMsg, resBody.Err,
			); err != nil {
				t.Error(err)
			}
		}
	}

	validReqBody := POSTReq{Username: "bob123", Password: "Myp4ssword!"}
	for _, c := range []struct {
		name              string
		reqBody           POSTReq
		validationErrs    ValidationErrors
		inviteCode        string
		inviteCodeErr     error
		teamSelectorErr   error
		userRecord        userTable.Record
		userSelectorErr   error
		hashedPwd         []byte
		hasherErr         error
		userInserterErr   error
		authToken         string
		tokenGeneratorErr error
		wantStatusCode    int
		assertFunc        func(*testing.T, *http.Response, string)
	}{
		{
			name:    "BasicValidatorErrs",
			reqBody: POSTReq{},
			validationErrs: ValidationErrors{
				Username: []string{usnTooLong}, Password: []string{pwdNoDigit},
			},
			inviteCode:        "",
			inviteCodeErr:     nil,
			teamSelectorErr:   nil,
			userRecord:        userTable.Record{},
			userSelectorErr:   nil,
			hashedPwd:         nil,
			hasherErr:         nil,
			userInserterErr:   nil,
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusBadRequest,
			assertFunc: assertOnValidationErrs(
				ValidationErrors{
					Username: []string{usnTooLong},
					Password: []string{pwdNoDigit},
				},
			),
		},
		{
			name:              "InvalidInviteCode",
			reqBody:           POSTReq{},
			validationErrs:    ValidationErrors{},
			inviteCode:        "someinvitecode",
			inviteCodeErr:     errors.New("an error"),
			teamSelectorErr:   nil,
			userRecord:        userTable.Record{},
			userSelectorErr:   nil,
			hashedPwd:         nil,
			hasherErr:         nil,
			userInserterErr:   nil,
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusBadRequest,
			assertFunc:        assertOnResErr("Invalid invite code."),
		},
		{
			name:              "TeamNotFound",
			reqBody:           POSTReq{},
			validationErrs:    ValidationErrors{},
			inviteCode:        "someinvitecode",
			inviteCodeErr:     nil,
			teamSelectorErr:   sql.ErrNoRows,
			userRecord:        userTable.Record{},
			userSelectorErr:   nil,
			hashedPwd:         nil,
			hasherErr:         nil,
			userInserterErr:   nil,
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusNotFound,
			assertFunc:        assertOnResErr("Team not found."),
		},
		{
			name:              "TeamSelectorErr",
			reqBody:           POSTReq{},
			validationErrs:    ValidationErrors{},
			inviteCode:        "someinvitecode",
			inviteCodeErr:     nil,
			teamSelectorErr:   sql.ErrConnDone,
			userRecord:        userTable.Record{},
			userSelectorErr:   nil,
			hashedPwd:         nil,
			hasherErr:         nil,
			userInserterErr:   nil,
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr(sql.ErrConnDone.Error()),
		},
		{
			name:              "UsernameTaken",
			reqBody:           POSTReq{},
			validationErrs:    ValidationErrors{},
			inviteCode:        "",
			inviteCodeErr:     nil,
			teamSelectorErr:   nil,
			userRecord:        userTable.Record{},
			userSelectorErr:   nil,
			hashedPwd:         nil,
			hasherErr:         nil,
			userInserterErr:   nil,
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusBadRequest,
			assertFunc: assertOnValidationErrs(
				ValidationErrors{
					Username: []string{"Username is already taken."},
				},
			),
		},
		{
			name:              "UserSelectorError",
			reqBody:           validReqBody,
			validationErrs:    ValidationErrors{},
			inviteCode:        "",
			inviteCodeErr:     nil,
			teamSelectorErr:   nil,
			userRecord:        userTable.Record{},
			userSelectorErr:   errors.New("user selector error"),
			hashedPwd:         nil,
			hasherErr:         nil,
			userInserterErr:   nil,
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr("user selector error"),
		},
		{
			name:              "CreateTeamError",
			reqBody:           validReqBody,
			validationErrs:    ValidationErrors{},
			inviteCode:        "",
			inviteCodeErr:     nil,
			teamSelectorErr:   nil,
			userRecord:        userTable.Record{},
			userSelectorErr:   sql.ErrNoRows,
			hashedPwd:         nil,
			hasherErr:         errors.New("hasher error"),
			userInserterErr:   nil,
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr("hasher error"),
		},
		{
			name:              "HasherError",
			reqBody:           validReqBody,
			validationErrs:    ValidationErrors{},
			inviteCode:        "",
			inviteCodeErr:     nil,
			teamSelectorErr:   nil,
			userRecord:        userTable.Record{},
			userSelectorErr:   sql.ErrNoRows,
			hashedPwd:         nil,
			hasherErr:         errors.New("hasher error"),
			userInserterErr:   nil,
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr("hasher error"),
		},
		{
			name:              "UserInserterError",
			reqBody:           validReqBody,
			validationErrs:    ValidationErrors{},
			inviteCode:        "",
			inviteCodeErr:     nil,
			teamSelectorErr:   nil,
			userRecord:        userTable.Record{},
			userSelectorErr:   sql.ErrNoRows,
			hashedPwd:         nil,
			hasherErr:         nil,
			userInserterErr:   errors.New("inserter error"),
			authToken:         "",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc:        assert.OnLoggedErr("inserter error"),
		},
		{
			name:              "TokenGeneratorError",
			reqBody:           validReqBody,
			validationErrs:    ValidationErrors{},
			inviteCode:        "",
			inviteCodeErr:     nil,
			teamSelectorErr:   nil,
			userRecord:        userTable.Record{},
			userSelectorErr:   sql.ErrNoRows,
			hashedPwd:         nil,
			hasherErr:         nil,
			userInserterErr:   nil,
			authToken:         "",
			tokenGeneratorErr: errors.New("token generator error"),
			wantStatusCode:    http.StatusInternalServerError,
			assertFunc: assertOnResErr(
				"You have been registered successfully but something went " +
					"wrong. Please log in using the credentials you " +
					"registered with.",
			),
		},
		{
			name:              "Success",
			reqBody:           validReqBody,
			inviteCode:        "",
			validationErrs:    ValidationErrors{},
			userRecord:        userTable.Record{},
			userSelectorErr:   sql.ErrNoRows,
			hashedPwd:         nil,
			hasherErr:         nil,
			userInserterErr:   nil,
			authToken:         "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			tokenGeneratorErr: nil,
			wantStatusCode:    http.StatusOK,
			assertFunc: func(
				t *testing.T, r *http.Response, _ string,
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
			userValidator.validationErrs = c.validationErrs
			inviteCodeValidator.Err = c.inviteCodeErr
			teamSelector.Err = c.teamSelectorErr
			userSelector.Rec = c.userRecord
			userSelector.Err = c.userSelectorErr
			hasher.hash = c.hashedPwd
			hasher.err = c.hasherErr
			userInserter.Err = c.userInserterErr
			tokenGenerator.AuthToken = c.authToken
			tokenGenerator.Err = c.tokenGeneratorErr

			// Prepare request and response recorder.
			reqBody, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPost,
				"?inviteCode="+c.inviteCode,
				bytes.NewReader(reqBody),
			)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			// Handle request with sut and get the result.
			sut.Handle(w, req, "")
			res := w.Result()

			// Assert on the status code.
			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			// Run case-specific assertions
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
