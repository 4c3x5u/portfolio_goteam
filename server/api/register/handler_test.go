package register

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"server/assert"
)

func TestHandler(t *testing.T) {
	// handler setup
	var (
		validator   = &fakeValidatorReq{}
		existorUser = &fakeExistorUser{}
		hasherPwd   = &fakeHasherPwd{}
		creatorUser = &fakeCreatorUser{}
	)

	for _, c := range []struct {
		name             string
		reqBody          *ReqBody
		validatorOutErrs *Errs
		existorOutExists bool
		existorOutErr    error
		hasherOutHash    []byte
		hasherOutErr     error
		wantStatusCode   int
		wantHandlerErrs  *Errs
	}{
		{
			name:             "ErrsValidator",
			reqBody:          &ReqBody{Username: "bobobobobobobobob", Password: "myNOdigitPASSWORD!"},
			validatorOutErrs: &Errs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
			existorOutExists: false,
			existorOutErr:    nil,
			hasherOutHash:    []byte{},
			hasherOutErr:     nil,
			wantStatusCode:   http.StatusBadRequest,
			wantHandlerErrs:  &Errs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
		},
		{
			name:             "ErrsCreator",
			reqBody:          &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			validatorOutErrs: nil,
			existorOutExists: true,
			existorOutErr:    nil,
			hasherOutHash:    []byte{},
			hasherOutErr:     nil,
			wantStatusCode:   http.StatusBadRequest,
			wantHandlerErrs:  &Errs{Username: []string{errHandlerUsernameTaken}},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			sut := NewHandler(validator, existorUser, hasherPwd, creatorUser)

			// parse response body - done only to assert tha the creator and
			// the validator receives the correct input based on the request
			// passed in
			reqBody, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			// set pre-determinate return values for Handler dependencies
			validator.outErrs = c.validatorOutErrs
			existorUser.outExists = c.existorOutExists
			existorUser.outErr = c.existorOutErr
			hasherPwd.outHash = c.hasherOutHash
			hasherPwd.outErr = c.hasherOutErr

			// create request (arrange)
			req, err := http.NewRequest("POST", "/register", bytes.NewReader(reqBody))
			if err != nil {
				t.Fatal(err)
			}
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			// send request (act)
			sut.ServeHTTP(w, req)

			// assert that the handler correctly passed in the arguments
			// for request validator and user creator
			assert.Equal(t, c.reqBody.Username, validator.inReqBody.Username)
			assert.Equal(t, c.reqBody.Password, validator.inReqBody.Password)

			// There are multiple breakpoints to the assertions here. Each stage
			// of the Handler is either ran or not ran based on the error
			// returns on Handler's dependencies. The conditionals below serve
			// to only run the assertions up onto the point where the Handler
			// exist execution.
			if c.validatorOutErrs == nil {
				assert.Equal(t, c.reqBody.Username, existorUser.inUsername)
				if existorUser.outErr == nil && existorUser.outExists != true {
					assert.Equal(t, c.reqBody.Password, hasherPwd.inPlaintext)
					if c.hasherOutErr == nil {
						creatorInPwd, err := json.Marshal(creatorUser.inArgs[1])
						if err != nil {
							t.Fatal(err)
						}
						hasherOutHash, err := json.Marshal(c.hasherOutHash)
						if err != nil {
							t.Fatal(err)
						}
						assert.EqualArr(t, hasherOutHash, creatorInPwd)
					}
				}
			}

			// make assertions on the status code and response body
			res := w.Result()
			resBody := &ResBody{}
			if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, c.wantStatusCode, res.StatusCode)
			assert.EqualArr(t, c.wantHandlerErrs.Username, resBody.Errs.Username)
			assert.EqualArr(t, c.wantHandlerErrs.Password, resBody.Errs.Password)
		})
	}
}
