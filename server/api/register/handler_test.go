package register

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
)

func TestHandler(t *testing.T) {
	// handler setup
	var (
		validatorReq   = &fakeValidatorReq{}
		existorUser    = &fakeExistorUser{}
		hasherPwd      = &fakeHasherPwd{}
		creatorUser    = &fakeCreatorUser{}
		creatorSession = &fakeCreatorSession{}
	)

	for _, c := range []struct {
		name                 string
		req                  *Req
		outErrValidatorReq   *Errs
		outResExistorUser    bool
		outErrExistorUser    error
		outResHasherPwd      []byte
		outErrHasherPwd      error
		outErrCreatorUser    error
		outErrCreatorSession error
		wantStatusCode       int
		wantFieldErrs        *Errs
	}{
		{
			name:                 "ErrValidator",
			req:                  &Req{Username: "bobobobobobobobob", Password: "myNOdigitPASSWORD!"},
			outErrValidatorReq:   &Errs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
			outResExistorUser:    false,
			outErrExistorUser:    nil,
			outResHasherPwd:      nil,
			outErrHasherPwd:      nil,
			outErrCreatorUser:    nil,
			outErrCreatorSession: nil,
			wantStatusCode:       http.StatusBadRequest,
			wantFieldErrs:        &Errs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
		},
		{
			name:                 "ResExistorTrue",
			req:                  &Req{Username: "bob21", Password: "Myp4ssword!"},
			outErrValidatorReq:   nil,
			outResExistorUser:    true,
			outErrExistorUser:    nil,
			outResHasherPwd:      nil,
			outErrHasherPwd:      nil,
			outErrCreatorUser:    nil,
			outErrCreatorSession: nil,
			wantStatusCode:       http.StatusBadRequest,
			wantFieldErrs:        &Errs{Username: []string{errFieldUsernameTaken}},
		},
		{
			name:                 "ErrExistor",
			req:                  &Req{Username: "bob2121", Password: "Myp4ssword!"},
			outErrValidatorReq:   nil,
			outResExistorUser:    false,
			outErrExistorUser:    errors.New("existor fatal error"),
			outResHasherPwd:      nil,
			outErrHasherPwd:      nil,
			outErrCreatorUser:    nil,
			outErrCreatorSession: nil,
			wantStatusCode:       http.StatusInternalServerError,
			wantFieldErrs:        nil,
		},
		{
			name:                 "ErrHasher",
			req:                  &Req{Username: "bob2121", Password: "Myp4ssword!"},
			outErrValidatorReq:   nil,
			outResExistorUser:    false,
			outErrExistorUser:    nil,
			outResHasherPwd:      nil,
			outErrHasherPwd:      errors.New("hasher fatal error"),
			outErrCreatorUser:    nil,
			outErrCreatorSession: nil,
			wantStatusCode:       http.StatusInternalServerError,
			wantFieldErrs:        nil,
		},
		{
			name:                 "ErrCreatorUser",
			req:                  &Req{Username: "bob2121", Password: "Myp4ssword!"},
			outErrValidatorReq:   nil,
			outResExistorUser:    false,
			outErrExistorUser:    nil,
			outResHasherPwd:      nil,
			outErrHasherPwd:      nil,
			outErrCreatorUser:    errors.New("creator fatal error"),
			outErrCreatorSession: nil,
			wantStatusCode:       http.StatusInternalServerError,
			wantFieldErrs:        nil,
		},
		{
			name:                 "ErrCreatorSession",
			req:                  &Req{Username: "bob2121", Password: "Myp4ssword!"},
			outErrValidatorReq:   nil,
			outResExistorUser:    false,
			outErrExistorUser:    nil,
			outResHasherPwd:      nil,
			outErrHasherPwd:      nil,
			outErrCreatorUser:    nil,
			outErrCreatorSession: errors.New("session creator error"),
			wantStatusCode:       http.StatusUnauthorized,
			wantFieldErrs:        &Errs{Session: errSession},
		},
		{
			name:                 "ResHandlerOK",
			req:                  &Req{Username: "bob2121", Password: "Myp4ssword!"},
			outErrValidatorReq:   nil,
			outResExistorUser:    false,
			outErrExistorUser:    nil,
			outResHasherPwd:      nil,
			outErrHasherPwd:      nil,
			outErrCreatorUser:    nil,
			outErrCreatorSession: nil,
			wantStatusCode:       http.StatusOK,
			wantFieldErrs:        nil,
		},
		// TODO: Expand – stages? Curried function that takes in *testing.T and
		//       whatever else arg needed to make its assertions. Simpler.
		// TODO: Abstract a Logger to make assertions on logged messages?
	} {
		t.Run(c.name, func(t *testing.T) {
			// set pre-determinate return values for Handler dependencies
			validatorReq.outErrs = c.outErrValidatorReq
			existorUser.outExists = c.outResExistorUser
			existorUser.outErr = c.outErrExistorUser
			hasherPwd.outHash = c.outResHasherPwd
			hasherPwd.outErr = c.outErrHasherPwd
			creatorUser.outErr = c.outErrCreatorUser
			creatorSession.outErr = c.outErrCreatorSession

			// create handler
			sut := NewHandler(validatorReq, existorUser, hasherPwd, creatorUser, creatorSession)

			// parse response body - done only to assert tha the creator and
			// the valiadtor receives the correct input based on the request
			// passed in
			reqBody, err := json.Marshal(c.req)
			if err != nil {
				t.Fatal(err)
			}

			// create request (arrange)
			req, err := http.NewRequest("POST", "/register", bytes.NewReader(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// send request (act)
			sut.ServeHTTP(w, req)

			// Input-based assertions to be run up onto the point where handler
			// stops execution. Conditionals serve to determine which
			// dependencies should have received their function arguments.
			assert.Equal(t, c.req.Username, validatorReq.inReq.Username)
			assert.Equal(t, c.req.Password, validatorReq.inReq.Password)
			if c.outErrValidatorReq == nil {
				// validatorReq.Validate doesn't error – existorUser.Exists is called
				assert.Equal(t, c.req.Username, existorUser.inUsername)
				if c.outErrExistorUser == nil && c.outResExistorUser == false {
					// existorUser.Exists return true and doesn't error - hasherPwd.Hash is called
					assert.Equal(t, c.req.Password, hasherPwd.inPlaintext)
					if c.outErrHasherPwd == nil {
						// hasherPwd.Hash doesn't error – creatorUser.Create is called
						assert.Equal(t, c.req.Username, creatorUser.inUsername)
						assert.Equal(t, string(c.outResHasherPwd), string(creatorUser.inPassword))
						if c.outErrCreatorUser == nil {
							// creatorUser.Create doesn't error – creatorSession.Create is called
							assert.Equal(t, c.req.Username, creatorSession.inUsername)
						}
					}
				}
			}

			// assert on status code
			res := w.Result()
			assert.Equal(t, c.wantStatusCode, res.StatusCode)

			// Assert on response body – however, there are some cases such as
			// internal server errors where an empty res body is returned and
			// these assertions are not run.
			if c.outErrExistorUser == nil && c.outErrHasherPwd == nil && c.outErrCreatorUser == nil {
				resBody := &Res{}
				if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}

				if c.wantFieldErrs != nil {
					// field errors - assert on them
					assert.EqualArr(t, c.wantFieldErrs.Username, resBody.Errs.Username)
					assert.EqualArr(t, c.wantFieldErrs.Password, resBody.Errs.Password)
					assert.Equal(t, c.wantFieldErrs.Session, resBody.Errs.Session)
				} else {
					// no field errors - assert on session token
					foundSessionToken := false
					for _, cookie := range res.Cookies() {
						if cookie.Name == "sessionToken" {
							foundSessionToken = true
						}
					}
					assert.Equal(t, true, foundSessionToken)
				}
			}
		})
	}
}
