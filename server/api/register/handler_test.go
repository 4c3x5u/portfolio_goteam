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

	"server/assert"
	"server/auth"
	"server/db"
	"server/log"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it
// behaves correctly.
func TestHandler(t *testing.T) {
	var (
		validator      = &fakeValidator{}
		userSelector   = &db.FakeUserSelector{}
		hasher         = &fakeHasher{}
		userInserter   = &db.FakeUserInserter{}
		tokenGenerator = &auth.FakeTokenGenerator{}
		logger         = &log.FakeLogger{}
	)
	sut := NewHandler(
		validator, userSelector, hasher, userInserter, tokenGenerator, logger,
	)

	t.Run("MethodNotAllowed", func(t *testing.T) {
		for _, httpMethod := range []string{
			http.MethodConnect, http.MethodDelete, http.MethodGet,
			http.MethodHead, http.MethodOptions, http.MethodPatch,
			http.MethodPut, http.MethodTrace,
		} {
			t.Run(httpMethod, func(t *testing.T) {
				req, err := http.NewRequest(httpMethod, "/register", nil)
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

	validReqBody := ReqBody{Username: "bob123", Password: "Myp4ssword!"}
	for _, c := range []struct {
		name                   string
		reqBody                ReqBody
		validatorOutErr        ValidationErrs
		userInserterOutRes     db.User
		userSelectorOutErr     error
		hasherOutRes           []byte
		hasherOutErr           error
		userInserterOutErr     error
		tokenGeneratorOutToken string
		tokenGeneratorOutErr   error
		wantStatusCode         int
		wantValidationErrs     ValidationErrs
	}{
		{
			name: "ValidatorErr",
			validatorOutErr: ValidationErrs{
				Username: []string{usnTooLong}, Password: []string{pwdNoDigit},
			},
			userInserterOutRes:     db.User{},
			userSelectorOutErr:     nil,
			hasherOutRes:           nil,
			hasherOutErr:           nil,
			userInserterOutErr:     nil,
			tokenGeneratorOutToken: "",
			tokenGeneratorOutErr:   nil,
			wantStatusCode:         http.StatusBadRequest,
			wantValidationErrs: ValidationErrs{
				Username: []string{usnTooLong}, Password: []string{pwdNoDigit},
			},
		},
		{
			name:                   "UsernameTaken",
			validatorOutErr:        ValidationErrs{},
			userInserterOutRes:     db.User{},
			userSelectorOutErr:     nil,
			hasherOutRes:           nil,
			hasherOutErr:           nil,
			userInserterOutErr:     nil,
			tokenGeneratorOutToken: "",
			tokenGeneratorOutErr:   nil,
			wantStatusCode:         http.StatusBadRequest,
			wantValidationErrs: ValidationErrs{
				Username: []string{errUsernameTaken},
			},
		},
		{
			name:                   "UserSelectorError",
			reqBody:                validReqBody,
			validatorOutErr:        ValidationErrs{},
			userInserterOutRes:     db.User{},
			userSelectorOutErr:     errors.New("user selector error"),
			hasherOutRes:           nil,
			hasherOutErr:           nil,
			userInserterOutErr:     nil,
			tokenGeneratorOutToken: "",
			tokenGeneratorOutErr:   nil,
			wantStatusCode:         http.StatusInternalServerError,
			wantValidationErrs:     ValidationErrs{},
		},
		{
			name:                   "HasherError",
			reqBody:                validReqBody,
			validatorOutErr:        ValidationErrs{},
			userInserterOutRes:     db.User{},
			userSelectorOutErr:     sql.ErrNoRows,
			hasherOutRes:           nil,
			hasherOutErr:           errors.New("hasher fatal error"),
			userInserterOutErr:     nil,
			tokenGeneratorOutToken: "",
			tokenGeneratorOutErr:   nil,
			wantStatusCode:         http.StatusInternalServerError,
			wantValidationErrs:     ValidationErrs{},
		},
		{
			name:                   "UserInserterError",
			reqBody:                validReqBody,
			validatorOutErr:        ValidationErrs{},
			userInserterOutRes:     db.User{},
			userSelectorOutErr:     sql.ErrNoRows,
			hasherOutRes:           nil,
			hasherOutErr:           nil,
			userInserterOutErr:     errors.New("inserter fatal error"),
			tokenGeneratorOutToken: "",
			tokenGeneratorOutErr:   nil,
			wantStatusCode:         http.StatusInternalServerError,
			wantValidationErrs:     ValidationErrs{},
		},
		{
			name:                   "TokenGeneratorError",
			reqBody:                validReqBody,
			validatorOutErr:        ValidationErrs{},
			userInserterOutRes:     db.User{},
			userSelectorOutErr:     sql.ErrNoRows,
			hasherOutRes:           nil,
			hasherOutErr:           nil,
			userInserterOutErr:     nil,
			tokenGeneratorOutToken: "",
			tokenGeneratorOutErr:   errors.New("token generator error"),
			wantStatusCode:         http.StatusUnauthorized,
			wantValidationErrs:     ValidationErrs{Auth: errAuth},
		},
		{
			name:                   "Success",
			reqBody:                validReqBody,
			validatorOutErr:        ValidationErrs{},
			userInserterOutRes:     db.User{},
			userSelectorOutErr:     sql.ErrNoRows,
			hasherOutRes:           nil,
			hasherOutErr:           nil,
			userInserterOutErr:     nil,
			tokenGeneratorOutToken: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			tokenGeneratorOutErr:   nil,
			wantStatusCode:         http.StatusOK,
			wantValidationErrs:     ValidationErrs{},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for sut's dependencies.
			validator.outErrs = c.validatorOutErr
			userSelector.OutRes = c.userInserterOutRes
			userSelector.OutErr = c.userSelectorOutErr
			hasher.outHash = c.hasherOutRes
			hasher.outErr = c.hasherOutErr
			userInserter.OutErr = c.userInserterOutErr
			tokenGenerator.OutRes = c.tokenGeneratorOutToken
			tokenGenerator.OutErr = c.tokenGeneratorOutErr

			// Prepare request and response recorder.
			reqBody, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPost, "", bytes.NewReader(reqBody),
			)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
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

			// If any validation errors were expected, assert on those.
			if c.wantValidationErrs.Any() {
				resBody := &ResBody{}
				if err = json.NewDecoder(res.Body).Decode(
					&resBody,
				); err != nil {
					t.Fatal(err)
				}

				if err = assert.EqualArr(
					c.wantValidationErrs.Username,
					resBody.ValidationErrs.Username,
				); err != nil {
					t.Error(err)
				}

				if err = assert.EqualArr(
					c.wantValidationErrs.Password,
					resBody.ValidationErrs.Password,
				); err != nil {
					t.Error(err)
				}

				if err = assert.Equal(
					c.wantValidationErrs.Auth, resBody.ValidationErrs.Auth,
				); err != nil {
					t.Error(err)
				}
			}

			switch c.wantStatusCode {
			case http.StatusOK:
				// 200 was expected - auth token must be set.
				authTokenFound := false
				for _, ck := range res.Cookies() {
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
				if !authTokenFound {
					t.Errorf("200 was expected but auth token was not set")
				}
			case http.StatusInternalServerError:
				// 500 was expected - an error must be logged.
				errFound := false
				for _, depErr := range []error{
					c.userSelectorOutErr,
					c.hasherOutErr,
					c.userInserterOutErr,
				} {
					if depErr != nil && depErr != sql.ErrNoRows {
						errFound = true
						if err = assert.Equal(
							log.LevelError, logger.InLevel,
						); err != nil {
							t.Error(err)
						}
						if err = assert.Equal(
							depErr.Error(), logger.InMessage,
						); err != nil {
							t.Error(err)
						}
					}
				}
				if !errFound {
					t.Errorf(
						"500 was expected but no errors were returned " +
							"from sut's dependencies",
					)
				}
			}
		})
	}
}
