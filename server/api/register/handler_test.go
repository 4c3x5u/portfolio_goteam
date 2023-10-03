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
	userTable "server/db/user"
	pkgLog "server/log"
)

// TestHandler tests the ServeHTTP method of Handler to assert that it
// behaves correctly.
func TestHandler(t *testing.T) {
	var (
		validator      = &fakeValidator{}
		userSelector   = &userTable.FakeSelector{}
		hasher         = &fakeHasher{}
		userInserter   = &userTable.FakeInserter{}
		tokenGenerator = &auth.FakeTokenGenerator{}
		log            = &pkgLog.FakeErrorer{}
	)
	sut := NewHandler(
		validator, userSelector, hasher, userInserter, tokenGenerator, log,
	)

	// Used in status 400 cases to assert on validation errors.
	assertOnValidationErrs := func(
		wantValidationErrs ValidationErrs,
	) func(*testing.T, *pkgLog.FakeErrorer, *http.Response) {
		return func(t *testing.T, _ *pkgLog.FakeErrorer, r *http.Response) {
			resBody := &ResBody{}
			if err := json.NewDecoder(r.Body).Decode(&resBody); err != nil {
				t.Fatal(err)
			}

			if err := assert.EqualArr(
				wantValidationErrs.Username,
				resBody.Errs.Username,
			); err != nil {
				t.Error(err)
			}

			if err := assert.EqualArr(
				wantValidationErrs.Password,
				resBody.Errs.Password,
			); err != nil {
				t.Error(err)
			}
		}
	}

	// Used in status 500 error cases to assert on the logged error message.
	assertOnLoggedErr := func(
		wantErrMsg string,
	) func(*testing.T, *pkgLog.FakeErrorer, *http.Response) {
		return func(t *testing.T, l *pkgLog.FakeErrorer, _ *http.Response) {
			if err := assert.Equal(wantErrMsg, l.InMessage); err != nil {
				t.Error(err)
			}
		}
	}

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
		userInserterOutRes   userTable.Record
		userSelectorOutErr   error
		hasherOutRes         []byte
		hasherOutErr         error
		userInserterOutErr   error
		tokenGeneratorOutRes string
		tokenGeneratorOutErr error
		wantStatusCode       int
		assertFunc           func(
			*testing.T, *pkgLog.FakeErrorer, *http.Response,
		)
	}{
		{
			name: "BasicValidatorErrs",
			validatorOutErr: ValidationErrs{
				Username: []string{usnTooLong}, Password: []string{pwdNoDigit},
			},
			userInserterOutRes:   userTable.Record{},
			userSelectorOutErr:   nil,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userInserterOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc: assertOnValidationErrs(
				ValidationErrs{
					Username: []string{usnTooLong},
					Password: []string{pwdNoDigit},
				},
			),
		},
		{
			name:                 "UsernameTaken",
			validatorOutErr:      ValidationErrs{},
			userInserterOutRes:   userTable.Record{},
			userSelectorOutErr:   nil,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userInserterOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusBadRequest,
			assertFunc: assertOnValidationErrs(
				ValidationErrs{
					Username: []string{"Username is already taken."},
				},
			),
		},
		{
			name:                 "UserSelectorError",
			reqBody:              validReqBody,
			validatorOutErr:      ValidationErrs{},
			userInserterOutRes:   userTable.Record{},
			userSelectorOutErr:   errors.New("user selector error"),
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userInserterOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assertOnLoggedErr("user selector error"),
		},
		{
			name:                 "HasherError",
			reqBody:              validReqBody,
			validatorOutErr:      ValidationErrs{},
			userInserterOutRes:   userTable.Record{},
			userSelectorOutErr:   sql.ErrNoRows,
			hasherOutRes:         nil,
			hasherOutErr:         errors.New("hasher error"),
			userInserterOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assertOnLoggedErr("hasher error"),
		},
		{
			name:                 "UserInserterError",
			reqBody:              validReqBody,
			validatorOutErr:      ValidationErrs{},
			userInserterOutRes:   userTable.Record{},
			userSelectorOutErr:   sql.ErrNoRows,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userInserterOutErr:   errors.New("inserter error"),
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc:           assertOnLoggedErr("inserter error"),
		},
		{
			name:                 "TokenGeneratorError",
			reqBody:              validReqBody,
			validatorOutErr:      ValidationErrs{},
			userInserterOutRes:   userTable.Record{},
			userSelectorOutErr:   sql.ErrNoRows,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userInserterOutErr:   nil,
			tokenGeneratorOutRes: "",
			tokenGeneratorOutErr: errors.New("token generator error"),
			wantStatusCode:       http.StatusInternalServerError,
			assertFunc: func(
				t *testing.T, _ *pkgLog.FakeErrorer, r *http.Response,
			) {
				resBody := &ResBody{}
				if err := json.NewDecoder(r.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}
				if err := assert.Equal(
					"You have been registered successfully but something "+
						"went wrong. Please log in using the credentials you "+
						"registered with.",
					resBody.Msg,
				); err != nil {
					t.Error(err)
				}
			},
		},
		{
			name:                 "Success",
			reqBody:              validReqBody,
			validatorOutErr:      ValidationErrs{},
			userInserterOutRes:   userTable.Record{},
			userSelectorOutErr:   sql.ErrNoRows,
			hasherOutRes:         nil,
			hasherOutErr:         nil,
			userInserterOutErr:   nil,
			tokenGeneratorOutRes: "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9...",
			tokenGeneratorOutErr: nil,
			wantStatusCode:       http.StatusOK,
			assertFunc: func(
				t *testing.T, _ *pkgLog.FakeErrorer, r *http.Response,
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
			validator.outErrs = c.validatorOutErr
			userSelector.OutRes = c.userInserterOutRes
			userSelector.OutErr = c.userSelectorOutErr
			hasher.outHash = c.hasherOutRes
			hasher.outErr = c.hasherOutErr
			userInserter.OutErr = c.userInserterOutErr
			tokenGenerator.OutRes = c.tokenGeneratorOutRes
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

			// Run case-specific assertions
			c.assertFunc(t, log, res)
		})
	}
}
