//go:build utest

package registerapi

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/cookie"
	"github.com/kxplxn/goteam/pkg/db"
	"github.com/kxplxn/goteam/pkg/db/usertbl"
	"github.com/kxplxn/goteam/pkg/log"
)

func TestHandler(t *testing.T) {
	var (
		userValidator = &fakeReqValidator{}
		hasher        = &fakeHasher{}
		inviteDecoder = &cookie.FakeStringDecoder[cookie.Invite]{}
		userInserter  = &db.FakeInserter[usertbl.User]{}
		authEncoder   = &cookie.FakeEncoder[cookie.Auth]{}
		log           = &log.FakeErrorer{}
	)
	sut := NewPostHandler(
		userValidator, inviteDecoder, hasher, userInserter, authEncoder, log,
	)

	// Used in status 400 cases to assert on validation errors.
	assertOnErrsValidate := func(
		wantValidationErrs ValidationErrs,
	) func(*testing.T, *http.Response, []any) {
		return func(t *testing.T, resp *http.Response, _ []any) {
			respBody := &PostResp{}
			if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
				t.Fatal(err)
			}

			assert.AllEqual(t.Error,
				respBody.ValidationErrs.Username, wantValidationErrs.Username,
			)

			assert.AllEqual(t.Error,
				respBody.ValidationErrs.Password, wantValidationErrs.Password,
			)
		}
	}

	validRBody := `{"username": "bob123", "password": "Myp4ssword!"}`
	for _, c := range []struct {
		name            string
		req             string
		errValidate     ValidationErrs
		tkInvite        string
		inviteDecoded   cookie.Invite
		errDecodeInvite error
		pwdHash         []byte
		errHash         error
		errInsertUser   error
		authToken       http.Cookie
		errEncodeAuth   error
		wantStatus      int
		assertFunc      func(*testing.T, *http.Response, []any)
	}{
		{
			name: "ErrsValidate",
			req:  "{}",
			errValidate: ValidationErrs{
				Username: []string{idTooLong}, Password: []string{pwdNoDigit},
			},
			tkInvite:        "",
			inviteDecoded:   cookie.Invite{},
			errDecodeInvite: nil,
			pwdHash:         nil,
			errHash:         nil,
			errInsertUser:   nil,
			authToken:       http.Cookie{},
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
			req:             "{}",
			errValidate:     ValidationErrs{},
			tkInvite:        "someinvitetoken",
			inviteDecoded:   cookie.Invite{},
			errDecodeInvite: errors.New("an error"),
			pwdHash:         nil,
			errHash:         nil,
			errInsertUser:   nil,
			authToken:       http.Cookie{},
			errEncodeAuth:   nil,
			wantStatus:      http.StatusBadRequest,
			assertFunc:      assert.OnRespErr("Invalid invite token."),
		},
		{
			name:            "ErrUsnTaken",
			req:             "{}",
			errValidate:     ValidationErrs{},
			tkInvite:        "",
			inviteDecoded:   cookie.Invite{},
			errDecodeInvite: nil,
			pwdHash:         nil,
			errHash:         nil,
			errInsertUser:   db.ErrDupKey,
			authToken:       http.Cookie{},
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
			req:           validRBody,
			errValidate:   ValidationErrs{},
			tkInvite:      "{}",
			inviteDecoded: cookie.Invite{},
			pwdHash:       nil,
			errHash:       errors.New("hasher error"),
			errInsertUser: nil,
			authToken:     http.Cookie{},
			errEncodeAuth: nil,
			wantStatus:    http.StatusInternalServerError,
			assertFunc:    assert.OnLoggedErr("hasher error"),
		},
		{
			name:          "ErrUsnTaken",
			req:           "{}",
			errValidate:   ValidationErrs{},
			tkInvite:      "",
			inviteDecoded: cookie.Invite{},
			pwdHash:       nil,
			errHash:       nil,
			errInsertUser: db.ErrDupKey,
			authToken:     http.Cookie{},
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
			req:           validRBody,
			errValidate:   ValidationErrs{},
			tkInvite:      "",
			inviteDecoded: cookie.Invite{},
			errInsertUser: errors.New("failed to put user"),
			pwdHash:       nil,
			errHash:       nil,
			authToken:     http.Cookie{},
			errEncodeAuth: nil,
			wantStatus:    http.StatusInternalServerError,
			assertFunc:    assert.OnLoggedErr("failed to put user"),
		},
		{
			name:          "ErrEncodeAuth",
			req:           validRBody,
			errValidate:   ValidationErrs{},
			tkInvite:      "",
			inviteDecoded: cookie.Invite{},
			pwdHash:       nil,
			errHash:       nil,
			errInsertUser: nil,
			authToken:     http.Cookie{},
			errEncodeAuth: errors.New("error encoding auth token"),
			wantStatus:    http.StatusInternalServerError,
			assertFunc: assert.OnRespErr(
				"You have been registered successfully but something went " +
					"wrong. Please log in using the credentials you " +
					"registered with.",
			),
		},
		{
			name: "Success",
			req:  validRBody,
			tkInvite: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ0ZWFtSUQiOi" +
				"J0ZWFtaWQifQ.1h_fmLJ1ip-Z6kJq9JXYDgGuWDPOcOf8abwCgKtHHcY",
			errValidate:   ValidationErrs{},
			errInsertUser: nil,
			pwdHash:       nil,
			errHash:       nil,
			authToken:     http.Cookie{Name: "foo", Value: "bar"},
			errEncodeAuth: nil,
			wantStatus:    http.StatusOK,
			assertFunc: func(t *testing.T, resp *http.Response, _ []any) {
				ck := resp.Cookies()[0]
				assert.Equal(t.Error, ck.Name, "foo")
				assert.Equal(t.Error, ck.Value, "bar")
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			userValidator.validationErrs = c.errValidate
			inviteDecoder.Res = c.inviteDecoded
			inviteDecoder.Err = c.errDecodeInvite
			hasher.hash = c.pwdHash
			hasher.err = c.errHash
			userInserter.Err = c.errInsertUser
			authEncoder.Res = c.authToken
			authEncoder.Err = c.errEncodeAuth
			w := httptest.NewRecorder()
			r := httptest.NewRequest(
				http.MethodPost,
				"/?inviteToken="+c.tkInvite,
				strings.NewReader(c.req),
			)

			sut.Handle(w, r, "")

			resp := w.Result()
			assert.Equal(t.Error, resp.StatusCode, c.wantStatus)
			c.assertFunc(t, resp, log.Args)
		})
	}
}
