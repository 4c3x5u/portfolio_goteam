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
		name                 string
		reqBody              ReqBody
		validatorOutErr      ValidationErrs
		userInserterOutRes   db.User
		userSelectorOutErr   error
		hasherOutRes         []byte
		hasherOutErr         error
		userInserterOutErr   error
		tokenGeneratorOutRes string
		tokenGeneratorOutErr error
		wantStatusCode       int
		wantValidationErrs   ValidationErrs
	}{
		{
			name: "ValidatorErr",
			validatorOutErr: ValidationErrs{
				Username: []string{usnTooLong}, Password: []string{pwdNoDigit},
			},
			userInserterOutRes:   db.User{},
			userSelectorOutErr:   nil,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userInserterOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
			wantValidationErrs: ValidationErrs{
				Username: []string{usnTooLong}, Password: []string{pwdNoDigit},
			},
		},
		{
			name:                 "UsernameTaken",
			validatorOutErr:      ValidationErrs{},
			userInserterOutRes:   db.User{},
			userSelectorOutErr:   nil,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userInserterOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
			wantValidationErrs: ValidationErrs{
				Username: []string{errUsernameTaken},
			},
		},
		{
			name:                 "UserSelectorError",
			reqBody:              validReqBody,
			validatorOutErr:      ValidationErrs{},
			userInserterOutRes:   db.User{},
			userSelectorOutErr:   errors.New("user selector error"),
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userInserterOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			wantValidationErrs:   ValidationErrs{},
		},
		{
			name:                 "HasherError",
			reqBody:              validReqBody,
			validatorOutErr:      ValidationErrs{},
			userInserterOutRes:   db.User{},
			userSelectorOutErr:   sql.ErrNoRows,
			hasherOutRes:         nil,
			hasherOutErr:         errors.New("hasher fatal error"),
			userInserterOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			wantValidationErrs:   ValidationErrs{},
		},
		{
			name:                 "UserInserterError",
			reqBody:              validReqBody,
			validatorOutErr:      ValidationErrs{},
			userInserterOutRes:   db.User{},
			userSelectorOutErr:   sql.ErrNoRows,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userInserterOutErr:   errors.New("inserter fatal error"),
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			wantValidationErrs:   ValidationErrs{},
		},
		{
			name:                 "TokenGeneratorError",
			reqBody:              validReqBody,
			validatorOutErr:      ValidationErrs{},
			userInserterOutRes:   db.User{},
			userSelectorOutErr:   sql.ErrNoRows,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userInserterOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: errors.New("token generator error"),
			wantStatusCode:       http.StatusUnauthorized,
			wantValidationErrs:   ValidationErrs{Auth: errAuth},
		},
		{
			name:                 "Success",
			reqBody:              validReqBody,
			validatorOutErr:      ValidationErrs{},
			userInserterOutRes:   db.User{},
			userSelectorOutErr:   sql.ErrNoRows,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userInserterOutErr:   nil,
			tokenGeneratorOutRes: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusOK,
			wantValidationErrs:   ValidationErrs{},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for Handler dependencies.
			validator.outErrs = c.validatorOutErr
			userSelector.OutRes = c.userInserterOutRes
			userSelector.OutErr = c.userSelectorOutErr
			hasher.outHash = c.hasherOutRes
			hasher.outErr = c.hasherOutErr
			userInserter.OutErr = c.userInserterOutErr
			tokenGenerator.OutRes = c.tokenGeneratorOutRes
			tokenGenerator.OutErr = c.tokenGeneratorOutErr

			// Parse request body.
			reqBody, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPost, "/register", bytes.NewReader(reqBody),
			)
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Send request (act).
			sut.ServeHTTP(w, req)

			// Assert on status code.
			res := w.Result()
			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			if c.wantStatusCode == http.StatusBadRequest {
				// 400 is expected - there must be validation errors in response
				// body.
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
			} else if c.wantStatusCode == http.StatusOK {
				// 200 is expected - auth token must be set.
				authTokenFound := false
				for _, ck := range res.Cookies() {
					if ck.Name == auth.CookieName {
						authTokenFound = true
						if err = assert.Equal(
							c.tokenGeneratorOutRes, ck.Value,
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
				if err = assert.Equal(true, authTokenFound); err != nil {
					t.Error(err)
				}
			}

			// DEPENDENCY-INPUT-BASED ASSERTIONS

			// If 405 isn't expected, validator must be called.
			if c.wantStatusCode == http.StatusMethodNotAllowed {
				return
			}
			if err = assert.Equal(
				c.reqBody.Username, validator.inReqBody.Username,
			); err != nil {
				t.Error(err)
			}
			if err = assert.Equal(
				c.reqBody.Password, validator.inReqBody.Password,
			); err != nil {
				t.Error(err)
			}

			// If no validator error is expected, user selector and db closer
			// must be called.
			if c.validatorOutErr.Any() {
				return
			}
			if err = assert.Equal(
				c.reqBody.Username, userSelector.InUserID,
			); err != nil {
				t.Error(err)
			}

			// If user is expected to not already exist, hasher must be called.
			if c.userSelectorOutErr != sql.ErrNoRows {
				return
			}
			if err = assert.Equal(
				c.reqBody.Password, hasher.inPlaintext,
			); err != nil {
				t.Error(err)
			}

			// If no hasher error is expected, user inserter must be called.
			if c.hasherOutErr != nil {
				return
			}
			if err = assert.Equal(
				c.reqBody.Username, userInserter.InUser.ID,
			); err != nil {
				t.Error(err)
			}
			if err = assert.Equal(
				string(c.hasherOutRes), string(userInserter.InUser.Password),
			); err != nil {
				t.Error(err)
			}

			// If no user inserter error is expected, token generator must be
			// called.
			if c.userInserterOutErr != nil {
				return
			}
			if err = assert.Equal(
				c.reqBody.Username, tokenGenerator.InSub,
			); err != nil {
				t.Error(err)
			}
		})
	}
}
