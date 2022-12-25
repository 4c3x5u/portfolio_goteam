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
		validator      = &fakeValidatorReq{}
		existorUser    = &fakeExistorUser{}
		hasherPwd      = &fakeHasherPwd{}
		creatorUser    = &fakeCreatorUser{}
		creatorSession = &fakeCreatorSession{}
	)

	for _, c := range []struct {
		name                 string
		req                  *Req
		outErrValidatorReq   *Errs
		outResExistor        bool
		outErrExistor        error
		outResHasher         []byte
		outErrHasherPwd      error
		outErrCreatorUser    error
		wantStatusCode       int
		wantFieldErrs        *Errs
		outErrCreatorSession error
	}{
		{
			name:                 "ErrValidator",
			req:                  &Req{Username: "bobobobobobobobob", Password: "myNOdigitPASSWORD!"},
			outErrValidatorReq:   &Errs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
			outResExistor:        false,
			outErrExistor:        nil,
			outResHasher:         nil,
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
			outResExistor:        true,
			outErrExistor:        nil,
			outResHasher:         nil,
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
			outResExistor:        false,
			outErrExistor:        errors.New("existor fatal error"),
			outResHasher:         nil,
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
			outResExistor:        false,
			outErrExistor:        nil,
			outResHasher:         nil,
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
			outResExistor:        false,
			outErrExistor:        nil,
			outResHasher:         nil,
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
			outResExistor:        false,
			outErrExistor:        nil,
			outResHasher:         nil,
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
			outResExistor:        false,
			outErrExistor:        nil,
			outResHasher:         nil,
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
			sut := NewHandler(validator, existorUser, hasherPwd, creatorUser, creatorSession)

			// parse response body - done only to assert tha the creator and
			// the validator receives the correct input based on the request
			// passed in
			reqBody, err := json.Marshal(c.req)
			if err != nil {
				t.Fatal(err)
			}

			// set pre-determinate return values for Handler dependencies
			validator.outErrs = c.outErrValidatorReq
			existorUser.outExists = c.outResExistor
			existorUser.outErr = c.outErrExistor
			hasherPwd.outHash = c.outResHasher
			hasherPwd.outErr = c.outErrHasherPwd
			creatorUser.outErr = c.outErrCreatorUser
			creatorSession.outErr = c.outErrCreatorSession

			// create request (arrange)
			req, err := http.NewRequest("POST", "/register", bytes.NewReader(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// send request (act)
			sut.ServeHTTP(w, req)

			// There are multiple breakpoints to the assertions here. Each stage
			// of the Handler is either run or not run based on the error
			// returns on Handler's dependencies. The conditionals below serve
			// to only run the assertions up to the point where the Handler
			// exits execution.
			assert.Equal(t, c.req.Username, validator.inReq.Username)
			assert.Equal(t, c.req.Password, validator.inReq.Password)
			if c.outErrValidatorReq == nil {
				assert.Equal(t, c.req.Username, existorUser.inUsername)
				if existorUser.outErr == nil && existorUser.outExists == false {
					assert.Equal(t, c.req.Password, hasherPwd.inPlaintext)
					if c.outErrHasherPwd == nil {
						assert.Equal(t, c.req.Username, creatorUser.inUsername)
						assert.Equal(t, string(c.outResHasher), string(creatorUser.inPassword))
						if c.outErrCreatorUser == nil {
							assert.Equal(t, c.req.Username, creatorSession.inUsername)
						}
					}
				}
			}

			// assert on status code
			res := w.Result()
			assert.Equal(t, c.wantStatusCode, res.StatusCode)

			// assert on response body – however, there are some cases such as
			// internal server errors where an empty res body is returned and
			// these assertions are not viable
			if c.outErrExistor == nil && c.outErrHasherPwd == nil && c.outErrCreatorUser == nil {
				resBody := &Res{}
				if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}

				if c.wantFieldErrs != nil {
					assert.EqualArr(t, c.wantFieldErrs.Username, resBody.Errs.Username)
					assert.EqualArr(t, c.wantFieldErrs.Password, resBody.Errs.Password)
					assert.Equal(t, c.wantFieldErrs.Session, resBody.Errs.Session)
				}
			}
		})
	}
}
