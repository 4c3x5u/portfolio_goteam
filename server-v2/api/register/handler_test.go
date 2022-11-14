package register

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server-v2/assert"
)

// todo: move to db tests
const usnTaken = "Username is already taken."

func TestHandler(t *testing.T) {
	// handler setup
	creator, validator := &fakeCreatorUser{}, &fakeValidatorReq{}
	sut := NewHandler(creator, validator)

	// test cases below should all return 400
	wantStatusCode := http.StatusBadRequest

	// test cases
	for _, c := range []struct {
		name          string
		reqBody       *ReqBody
		errsValidator *Errs
		errsCreator   *Errs
		wantErrs      *Errs
	}{
		{
			name:          "ErrsValidator",
			reqBody:       &ReqBody{Username: "bobobobobobobobob", Password: "myNOdigitPASSWORD!"},
			errsValidator: &Errs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
			errsCreator:   nil,
			wantErrs:      &Errs{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
		},
		{
			name:          "ErrsCreator",
			reqBody:       &ReqBody{Username: "bob21", Password: "Myp4ssword!"},
			errsValidator: nil,
			errsCreator:   &Errs{Username: []string{usnTaken}},
			wantErrs:      &Errs{Username: []string{usnTaken}},
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
			validator.outErrs = c.errsValidator
			creator.outErrs = c.errsCreator

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
			if c.errsValidator == nil {
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
			assert.EqualArr(t, c.wantErrs.Username, resBody.Errs.Username)
			assert.EqualArr(t, c.wantErrs.Password, resBody.Errs.Password)
		})
	}
}
