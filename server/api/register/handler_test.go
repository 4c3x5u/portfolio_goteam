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
	creator, validator := &fakeCreatorUser{}, &fakeValidatorReq{}
	sut := NewHandler(creator, validator)

	// test cases below should all return 400
	wantStatusCode := http.StatusBadRequest

	// test cases
	for _, c := range []struct {
		name            string
		reqBody         *ReqBody
		outErrValidator *Errs
		outErrCreator   error
		wantErrsHandler *Errs
	}{
		{
			name:            "ErrsValidator",
			reqBody:         &ReqBody{Username: "bobobobobobobobob", Password: "myNOdigitPASSWORD!"},
			outErrValidator: &Errs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
			outErrCreator:   nil,
			wantErrsHandler: &Errs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
		},
		{
			name:            "ErrsCreator",
			reqBody:         &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			outErrValidator: nil,
			outErrCreator:   errCreatorUsernameTaken,
			wantErrsHandler: &Errs{Username: []string{errHandlerUsernameTaken}},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// parse response body - done only to assert tha the creator and
			// the validator receives the correct input based on the request
			// passed in
			reqBody, err := json.Marshal(c.reqBody)
			if err != nil {
				t.Fatal(err)
			}

			// set the wantOutErrs return on fake validator and creator (arrange)
			validator.outErrs = c.outErrValidator
			creator.outErr = c.outErrCreator

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

			// When errors occur on validator, the handler code will
			// terminate and creator will not be called, causing these
			// assertions to fail. Only make them if the validator is set
			// to return nil Errs.
			if c.outErrValidator == nil {
				assert.Equal(t, c.reqBody.Username, creator.inUsername)
				assert.Equal(t, c.reqBody.Password, creator.inPassword)
			}

			// make assertions on the status code and response body (assert)
			res := w.Result()
			assert.Equal(t, wantStatusCode, res.StatusCode)
			resBody := &ResBody{}
			if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
				t.Fatal(err)
			}
			assert.EqualArr(t, c.wantErrsHandler.Username, resBody.Errs.Username)
			assert.EqualArr(t, c.wantErrsHandler.Password, resBody.Errs.Password)
		})
	}
}
