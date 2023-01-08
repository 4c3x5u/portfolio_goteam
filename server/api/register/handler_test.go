package register

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
	"server/db"
	"server/token"
)

func TestHandler(t *testing.T) {
	// handler setup
	var (
		validator      = &fakeValidator{}
		userReader     = &db.FakeUserReader{}
		hasher         = &fakeHasher{}
		userCreator    = &db.FakeUserCreator{}
		tokenGenerator = &token.FakeGenerator{}
	)
	sut := NewHandler(validator, userReader, hasher, userCreator, tokenGenerator)

	for _, c := range []struct {
		name                 string
		httpMethod           string
		reqBody              ReqBody
		validatorOutErr      ValidationErrs
		userReaderOutRes     db.User
		userReaderOutErr     error
		hasherOutRes         []byte
		hasherOutErr         error
		userCreatorOutErr    error
		tokenGeneratorOutRes string
		tokenGeneratorOutErr error
		wantStatusCode       int
		wantFieldErrs        ValidationErrs
	}{
		{
			name:                 "HttpMethodError",
			httpMethod:           http.MethodGet,
			reqBody:              ReqBody{Username: "bob2121", Password: "Myp4ssword!"},
			validatorOutErr:      ValidationErrs{},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userCreatorOutErr:    nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusMethodNotAllowed,
			wantFieldErrs:        ValidationErrs{},
		},
		{
			name:                 "ValidatorError",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bobobobobobobobob", Password: "myNOdigitPASSWORD!"},
			validatorOutErr:      ValidationErrs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userCreatorOutErr:    nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
			wantFieldErrs:        ValidationErrs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
		},
		{
			name:                 "UsernameTaken",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			validatorOutErr:      ValidationErrs{},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     nil,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userCreatorOutErr:    nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
			wantFieldErrs:        ValidationErrs{Username: []string{strErrUsernameTaken}},
		},
		{
			name:                 "UserReaderError",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob2121", Password: "Myp4ssword!"},
			validatorOutErr:      ValidationErrs{},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     errors.New("user reader error"),
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userCreatorOutErr:    nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			wantFieldErrs:        ValidationErrs{},
		},
		{
			name:                 "HasherError",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob2121", Password: "Myp4ssword!"},
			validatorOutErr:      ValidationErrs{},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     sql.ErrNoRows,
			hasherOutRes:         nil,
			hasherOutErr:         errors.New("hasher fatal error"),
			userCreatorOutErr:    nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			wantFieldErrs:        ValidationErrs{},
		},
		{
			name:                 "UserCreatorError",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob2121", Password: "Myp4ssword!"},
			validatorOutErr:      ValidationErrs{},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     sql.ErrNoRows,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userCreatorOutErr:    errors.New("creator fatal error"),
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			wantFieldErrs:        ValidationErrs{},
		},
		{
			name:                 "TokenGeneratorError",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob2121", Password: "Myp4ssword!"},
			validatorOutErr:      ValidationErrs{},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     sql.ErrNoRows,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userCreatorOutErr:    nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: errors.New("token generator error"),
			wantStatusCode:       http.StatusUnauthorized,
			wantFieldErrs:        ValidationErrs{Auth: errAuth},
		},
		{
			name:                 "Success",
			httpMethod:           http.MethodPost,
			reqBody:              ReqBody{Username: "bob2121", Password: "Myp4ssword!"},
			validatorOutErr:      ValidationErrs{},
			userReaderOutRes:     db.User{},
			userReaderOutErr:     sql.ErrNoRows,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userCreatorOutErr:    nil,
			tokenGeneratorOutRes: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusOK,
			wantFieldErrs:        ValidationErrs{},
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
			tokenGenerator.OutRes = c.tokenGeneratorOutRes
			tokenGenerator.OutErr = c.tokenGeneratorOutErr

			// Parse request body.
			reqBody, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(c.httpMethod, "/register", bytes.NewReader(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// Send request (act).
			sut.ServeHTTP(w, req)

			// Input-based assertions to be run up onto the point where handler
			// stops execution. Conditionals serve to determine which
			// dependencies should have received their function arguments.
			if c.httpMethod == http.MethodPost {
				if err = assert.Equal(c.reqBody.Username, validator.inReqBody.Username); err != nil {
					t.Error(err)
				}
				if err = assert.Equal(c.reqBody.Password, validator.inReqBody.Password); err != nil {
					t.Error(err)
				}

				if !c.validatorOutErr.Any() {
					// validator.Validate doesn't error – userReader.Exists is called.
					if err = assert.Equal(c.reqBody.Username, userReader.InArg); err != nil {
						t.Error(err)
					}
					if c.userReaderOutErr == sql.ErrNoRows {
						// userReader.Exists returns sql.ErrNoRows - hasher.Hash is called.
						if err = assert.Equal(c.reqBody.Password, hasher.inPlaintext); err != nil {
							t.Error(err)
						}

						if c.hasherOutErr == nil {
							// hasher.Hash doesn't error – userCreator.Create is called.
							if err = assert.Equal(c.reqBody.Username, userCreator.InArg.Username); err != nil {
								t.Error(err)
							}
							if err = assert.Equal(string(c.hasherOutRes), string(userCreator.InArg.Password)); err != nil {
								t.Error(err)
							}
							if c.userCreatorOutErr == nil {
								// userCreator.Create doesn't error – tokenGenerator.Create is called.
								if err = assert.Equal(c.reqBody.Username, tokenGenerator.InSub); err != nil {
									t.Error(err)
								}
							}
						}
					}
				}
			}

			// Assert on status code.
			res := w.Result()
			if err = assert.Equal(c.wantStatusCode, res.StatusCode); err != nil {
				t.Error(err)
			}

			// Assert on response body – however, there are some cases such as
			// internal server errors where an empty res body is returned and
			// these assertions are not run.
			if c.httpMethod != http.MethodPost ||
				c.userReaderOutErr != nil ||
				c.hasherOutErr != nil ||
				c.userCreatorOutErr != nil ||
				c.wantStatusCode == http.StatusOK {
				return
			}

			resBody := &ResBody{}
			if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
				t.Fatal(err)
			}

			if c.wantFieldErrs.Any() {
				// field errors - assert on them
				if err = assert.EqualArr(c.wantFieldErrs.Username, resBody.ValidationErrs.Username); err != nil {
					t.Error(err)
				}
				if err = assert.EqualArr(c.wantFieldErrs.Password, resBody.ValidationErrs.Password); err != nil {
					t.Error(err)
				}
				if err = assert.Equal(c.wantFieldErrs.Auth, resBody.ValidationErrs.Auth); err != nil {
					t.Error(err)
				}
			} else {
				// no field errors - assert on auth token
				tokenFound := false
				for _, cookie := range res.Cookies() {
					if cookie.Name == "authToken" {
						tokenFound = true
						if err = assert.Equal(c.tokenGeneratorOutRes, cookie.Value); err != nil {
							t.Error(err)
						}
					}
				}
				if err = assert.Equal(true, tokenFound); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
