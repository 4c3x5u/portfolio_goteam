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
	"server/cookie"
	"server/db"
)

// TestHandler tests Handler to ensure that it behaves correctly on all
// possible scenarios.
func TestHandler(t *testing.T) {
	var (
		validator       = &fakeValidator{}
		userReader      = &db.FakeUserReader{}
		hasher          = &fakeHasher{}
		userCreator     = &db.FakeUserCreator{}
		cookieGenerator = &cookie.FakeAuthGenerator{}
		cookieExpiry    = time.Now().Add(1 * time.Hour).Truncate(1 * time.Second).UTC()
		dbCloser        = &db.FakeCloser{}
	)
	sut := NewHandler(validator, userReader, hasher, userCreator, cookieGenerator, dbCloser)

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

				if err = assert.Equal(http.StatusMethodNotAllowed, w.Result().StatusCode); err != nil {
					t.Error(err)
				}
			})
		}
	})

	for _, c := range []struct {
		name                  string
		reqBody               ReqBody
		validatorOutErr       ValidationErrs
		userReaderOutRes      db.User
		userReaderOutErr      error
		hasherOutRes          []byte
		hasherOutErr          error
		userCreatorOutErr     error
		cookieGeneratorOutRes *http.Cookie
		cookieGeneratorOutErr error
		wantStatusCode        int
		wantValidationErrs    ValidationErrs
	}{
		{
			name:                  "ValidatorError",
			reqBody:               ReqBody{Username: "bobobobobobobobob", Password: "myNOdigitPASSWORD!"},
			validatorOutErr:       ValidationErrs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      nil,
			hasherOutRes:          nil,
			hasherOutErr:          nil,
			userCreatorOutErr:     nil,
			cookieGeneratorOutRes: nil,
			cookieGeneratorOutErr: nil,
			wantStatusCode:        http.StatusBadRequest,
			wantValidationErrs:    ValidationErrs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
		},
		{
			name:                  "UsernameTaken",
			reqBody:               ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			validatorOutErr:       ValidationErrs{},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      nil,
			hasherOutRes:          nil,
			hasherOutErr:          nil,
			userCreatorOutErr:     nil,
			cookieGeneratorOutRes: nil,
			cookieGeneratorOutErr: nil,
			wantStatusCode:        http.StatusBadRequest,
			wantValidationErrs:    ValidationErrs{Username: []string{errUsernameTaken}},
		},
		{
			name:                  "UserReaderError",
			reqBody:               ReqBody{Username: "bob2121", Password: "Myp4ssword!"},
			validatorOutErr:       ValidationErrs{},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      errors.New("user reader error"),
			hasherOutRes:          nil,
			hasherOutErr:          nil,
			userCreatorOutErr:     nil,
			cookieGeneratorOutRes: nil,
			cookieGeneratorOutErr: nil,
			wantStatusCode:        http.StatusInternalServerError,
			wantValidationErrs:    ValidationErrs{},
		},
		{
			name:                  "HasherError",
			reqBody:               ReqBody{Username: "bob2121", Password: "Myp4ssword!"},
			validatorOutErr:       ValidationErrs{},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      sql.ErrNoRows,
			hasherOutRes:          nil,
			hasherOutErr:          errors.New("hasher fatal error"),
			userCreatorOutErr:     nil,
			cookieGeneratorOutRes: nil,
			cookieGeneratorOutErr: nil,
			wantStatusCode:        http.StatusInternalServerError,
			wantValidationErrs:    ValidationErrs{},
		},
		{
			name:                  "UserCreatorError",
			reqBody:               ReqBody{Username: "bob2121", Password: "Myp4ssword!"},
			validatorOutErr:       ValidationErrs{},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      sql.ErrNoRows,
			hasherOutRes:          nil,
			hasherOutErr:          nil,
			userCreatorOutErr:     errors.New("creator fatal error"),
			cookieGeneratorOutRes: nil,
			cookieGeneratorOutErr: nil,
			wantStatusCode:        http.StatusInternalServerError,
			wantValidationErrs:    ValidationErrs{},
		},
		{
			name:                  "TokenGeneratorError",
			reqBody:               ReqBody{Username: "bob2121", Password: "Myp4ssword!"},
			validatorOutErr:       ValidationErrs{},
			userReaderOutRes:      db.User{},
			userReaderOutErr:      sql.ErrNoRows,
			hasherOutRes:          nil,
			hasherOutErr:          nil,
			userCreatorOutErr:     nil,
			cookieGeneratorOutRes: nil,
			cookieGeneratorOutErr: errors.New("token generator error"),
			wantStatusCode:        http.StatusUnauthorized,
			wantValidationErrs:    ValidationErrs{Auth: errAuth},
		},
		{
			name:              "Success",
			reqBody:           ReqBody{Username: "bob2121", Password: "Myp4ssword!"},
			validatorOutErr:   ValidationErrs{},
			userReaderOutRes:  db.User{},
			userReaderOutErr:  sql.ErrNoRows,
			hasherOutRes:      nil,
			hasherOutErr:      nil,
			userCreatorOutErr: nil,
			cookieGeneratorOutRes: &http.Cookie{
				Name:    "authToken",
				Value:   "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
				Expires: cookieExpiry,
			},
			cookieGeneratorOutErr: nil,
			wantStatusCode:        http.StatusOK,
			wantValidationErrs:    ValidationErrs{},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// Set pre-determinate return values for Handler dependencies.
			validator.outErrs = c.validatorOutErr
			userReader.OutRes = c.userReaderOutRes
			userReader.OutErr = c.userReaderOutErr
			hasher.outHash = c.hasherOutRes
			hasher.outErr = c.hasherOutErr
			userCreator.OutErr = c.userCreatorOutErr
			cookieGenerator.OutRes = c.cookieGeneratorOutRes
			cookieGenerator.OutErr = c.cookieGeneratorOutErr

			// Parse request body.
			reqBody, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(http.MethodPost, "/register", bytes.NewReader(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Send request (act).
			sut.ServeHTTP(w, req)

			// Assert on status code.
			res := w.Result()
			if err = assert.Equal(c.wantStatusCode, res.StatusCode); err != nil {
				t.Error(err)
			}

			if c.wantStatusCode == http.StatusBadRequest {
				// 400 is expected - there must be validation errors in request body.
				resBody := &ResBody{}
				if err = json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}

				if err = assert.EqualArr(c.wantValidationErrs.Username, resBody.ValidationErrs.Username); err != nil {
					t.Error(err)
				}
				if err = assert.EqualArr(c.wantValidationErrs.Password, resBody.ValidationErrs.Password); err != nil {
					t.Error(err)
				}
				if err = assert.Equal(c.wantValidationErrs.Auth, resBody.ValidationErrs.Auth); err != nil {
					t.Error(err)
				}
			} else if c.wantStatusCode == http.StatusOK {
				// 200 is expected - auth token must be set.
				authTokenFound := false
				for _, ck := range res.Cookies() {
					if ck.Name == "authToken" {
						authTokenFound = true
						if err = assert.Equal(c.cookieGeneratorOutRes.Value, ck.Value); err != nil {
							t.Error(err)
						}
						if err = assert.Equal(c.cookieGeneratorOutRes.Expires, ck.Expires); err != nil {
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
			if err = assert.Equal(c.reqBody.Username, validator.inReqBody.Username); err != nil {
				t.Error(err)
			}
			if err = assert.Equal(c.reqBody.Password, validator.inReqBody.Password); err != nil {
				t.Error(err)
			}

			// If no validator error is expected, userReader and dbCloser must be called.
			if c.validatorOutErr.Any() {
				return
			}
			if err = assert.Equal(c.reqBody.Username, userReader.InArg); err != nil {
				t.Error(err)
			}

			// If user is expected to not already exist, hasher must be called.
			if c.userReaderOutErr != sql.ErrNoRows {
				return
			}
			if err = assert.Equal(c.reqBody.Password, hasher.inPlaintext); err != nil {
				t.Error(err)
			}
			if err = assert.True(dbCloser.IsCalled); err != nil {
				t.Error(err)
			}

			// If no hasher error is expected, user creator must be called.
			if c.hasherOutErr != nil {
				return
			}
			if err = assert.Equal(c.reqBody.Username, userCreator.InArg.Username); err != nil {
				t.Error(err)
			}
			if err = assert.Equal(string(c.hasherOutRes), string(userCreator.InArg.Password)); err != nil {
				t.Error(err)
			}

			// If no user creator error is expected, cookie generator must be called.
			if c.userCreatorOutErr != nil {
				return
			}
			if err = assert.Equal(c.reqBody.Username, cookieGenerator.InSub); err != nil {
				t.Error(err)
			}
		})
	}
}
