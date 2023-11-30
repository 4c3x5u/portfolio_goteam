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

	"github.com/kxplxn/goteam/internal/api"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/auth"
	teamTable "github.com/kxplxn/goteam/pkg/dbaccess/team"
	userTable "github.com/kxplxn/goteam/pkg/dbaccess/user"
	pkgLog "github.com/kxplxn/goteam/pkg/log"
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

			assert.AllEqual(t.Error,
				wantValidationErrs.Username,
				resBody.ValidationErrs.Username,
			)

			assert.AllEqual(t.Error,
				wantValidationErrs.Password,
				resBody.ValidationErrs.Password,
			)
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
			assert.Equal(t.Error, resBody.Err, wantErrMsg)
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
			assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)

			// Run case-specific assertions
			c.assertFunc(t, res, log.InMessage)
		})
	}
}
