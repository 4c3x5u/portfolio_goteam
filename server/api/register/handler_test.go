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
		validator     = &fakeValidatorReq{}
		existorUser   = &fakeExistorUser{}
		hasherPwd     = &fakeHasherPwd{}
		creatorUser   = &fakeCreatorUser{}
		keeperSession = &fakeKeeperSession{}
	)

	for _, c := range []struct {
		name            string
		reqBody         *Req
		outErrValidator *Errs
		outResExistor   bool
		outErrExistor   error
		outResHasher    []byte
		outErrHasher    error
		outErrCreator   error
		wantStatusCode  int
		wantFieldErrs   *Errs
	}{
		{
			name:            "ErrValidator",
			reqBody:         &Req{Username: "bobobobobobobobob", Password: "myNOdigitPASSWORD!"},
			outErrValidator: &Errs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
			outResExistor:   false,
			outErrExistor:   nil,
			outResHasher:    nil,
			outErrHasher:    nil,
			outErrCreator:   nil,
			wantStatusCode:  http.StatusBadRequest,
			wantFieldErrs:   &Errs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
		},
		{
			name:            "ResExistorTrue",
			reqBody:         &Req{Username: "bob21", Password: "Myp4ssword!"},
			outErrValidator: nil,
			outResExistor:   true,
			outErrExistor:   nil,
			outResHasher:    nil,
			outErrHasher:    nil,
			outErrCreator:   nil,
			wantStatusCode:  http.StatusBadRequest,
			wantFieldErrs:   &Errs{Username: []string{errFieldUsernameTaken}},
		},
		{
			name:            "ErrExistor",
			reqBody:         &Req{Username: "bob2121", Password: "Myp4ssword!"},
			outErrValidator: nil,
			outResExistor:   false,
			outErrExistor:   errors.New("existor fatal error"),
			outResHasher:    nil,
			outErrHasher:    nil,
			outErrCreator:   nil,
			wantStatusCode:  http.StatusInternalServerError,
			wantFieldErrs:   nil,
		},
		{
			name:            "ErrHasher",
			reqBody:         &Req{Username: "bob2121", Password: "Myp4ssword!"},
			outErrValidator: nil,
			outResExistor:   false,
			outErrExistor:   nil,
			outResHasher:    nil,
			outErrHasher:    errors.New("hasher fatal error"),
			outErrCreator:   nil,
			wantStatusCode:  http.StatusInternalServerError,
			wantFieldErrs:   nil,
		},
		{
			name:            "ErrCreator",
			reqBody:         &Req{Username: "bob2121", Password: "Myp4ssword!"},
			outErrValidator: nil,
			outResExistor:   false,
			outErrExistor:   nil,
			outResHasher:    nil,
			outErrHasher:    nil,
			outErrCreator:   errors.New("creator fatal error"),
			wantStatusCode:  http.StatusInternalServerError,
			wantFieldErrs:   nil,
		},
		{
			name:            "ResHandlerOK",
			reqBody:         &Req{Username: "bob2121", Password: "Myp4ssword!"},
			outErrValidator: nil,
			outResExistor:   false,
			outErrExistor:   nil,
			outResHasher:    nil,
			outErrHasher:    nil,
			outErrCreator:   nil,
			wantStatusCode:  http.StatusOK,
			wantFieldErrs:   nil,
		},
		// TODO: Expand – stages? Curried function that takes in *testing.T and
		//       whatever else arg needed to make its assertions. Simpler.
		// TODO: Abstract a Logger to make assertions on logged messages?
	} {
		t.Run(c.name, func(t *testing.T) {
			sut := NewHandler(validator, existorUser, hasherPwd, creatorUser, keeperSession)

			// parse response body - done only to assert tha the creator and
			// the validator receives the correct input based on the request
			// passed in
			reqBody, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			// set pre-determinate return values for Handler dependencies
			validator.outErrs = c.outErrValidator
			existorUser.outExists = c.outResExistor
			existorUser.outErr = c.outErrExistor
			hasherPwd.outHash = c.outResHasher
			hasherPwd.outErr = c.outErrHasher
			creatorUser.outErr = c.outErrCreator

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
			assert.Equal(t, c.reqBody.Username, validator.inReqBody.Username)
			assert.Equal(t, c.reqBody.Password, validator.inReqBody.Password)
			if c.outErrValidator == nil {
				assert.Equal(t, c.reqBody.Username, existorUser.inUsername)
				if existorUser.outErr == nil && existorUser.outExists == false {
					assert.Equal(t, c.reqBody.Password, hasherPwd.inPlaintext)
					if c.outErrHasher == nil {
						inCreatorPwd, err := json.Marshal(creatorUser.inArgs[1])
						if err != nil {
							t.Fatal(err)
						}
						outHasherRes, err := json.Marshal(c.outResHasher)
						if err != nil {
							t.Fatal(err)
						}
						assert.EqualArr(t, inCreatorPwd, outHasherRes)
					}
				}
			}

			// assert on status code
			res := w.Result()
			assert.Equal(t, c.wantStatusCode, res.StatusCode)

			// assert on response body – however, there are some cases such as
			// internal server errors where an empty res body is returned and
			// these assertions are not viable
			if c.outErrExistor == nil && c.outErrHasher == nil && c.outErrCreator == nil {
				resBody := &Res{}
				if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}

				if c.wantFieldErrs != nil {
					assert.EqualArr(t, c.wantFieldErrs.Username, resBody.Errs.Username)
					assert.EqualArr(t, c.wantFieldErrs.Password, resBody.Errs.Password)
				}
			}
		})
	}
}
