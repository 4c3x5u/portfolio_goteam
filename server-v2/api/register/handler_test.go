package register

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/kxplxn/goteam/server-v2/assert"
)

func TestHandler(t *testing.T) {
	t.Run("Errs", func(t *testing.T) {
		const (
			usnTooShort = "Username cannot be shorter than 5 characters."
			usnTooLong  = "Username cannot be longer than 15 characters."
			usnTaken    = "Username is already taken."

			pwdTooShort = "Password cannot be shorter than 8 characters."
			pwdNoDigit  = "Password must contain a digit (0-9)."
		)

		// create an empty request body - since we're using fakes for user
		// creator and request validator that do not actually process the
		// request, it's irrevelant what it contains
		reqBody, err := json.Marshal(&ReqBody{})
		if err != nil {
			t.Fatal(err)
		}

		// handler setup
		creator, validator := &fakeCreatorUser{}, &fakeValidatorReq{}
		sut := NewHandler(creator, validator)

		// test cases below should all return 400
		wantStatusCode := http.StatusBadRequest

		// test cases
		for _, c := range []struct {
			name          string
			errsValidator *ErrsValidation
			errsCreator   *ErrsValidation
			wantErrs      *ErrsValidation
		}{
			{
				name:          "Validator1",
				errsValidator: &ErrsValidation{Username: []string{usnTooShort}, Password: []string{pwdTooShort}},
				errsCreator:   nil,
				wantErrs:      &ErrsValidation{Username: []string{usnTooShort}, Password: []string{pwdTooShort}},
			},
			{
				name:          "Validator2",
				errsValidator: &ErrsValidation{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
				errsCreator:   nil,
				wantErrs:      &ErrsValidation{Username: []string{usnTooLong}, Password: []string{pwdNoDigit}},
			},
			{
				name:          "Creator",
				errsValidator: nil,
				errsCreator:   &ErrsValidation{Username: []string{usnTaken}},
				wantErrs:      &ErrsValidation{Username: []string{usnTaken}},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				// set the errs return on fake validator and creator (arrange)
				validator.errs = c.errsValidator
				creator.errs = c.errsCreator

				// create request (arrange)
				req, err := http.NewRequest("POST", "/register", bytes.NewReader(reqBody))
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				// send request (act)
				sut.ServeHTTP(w, req)

				// make assertions on the status code and response body (assert)
				res := w.Result()
				assert.Equal(t, wantStatusCode, res.StatusCode)
				resBody := &ResBody{}
				if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}
				assert.EqualArr(t, c.wantErrs.Username, resBody.ErrsValidation.Username)
				assert.EqualArr(t, c.wantErrs.Password, resBody.ErrsValidation.Password)
			})
		}
	})
}
