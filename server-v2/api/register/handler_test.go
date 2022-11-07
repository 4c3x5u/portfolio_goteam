package register

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/kxplxn/goteam/server-v2/assert"
)

// ValidationTestCase defines the values that are commonly necessary between
// validation tests.
type ValidationTestCase struct {
	name     string
	input    string
	wantErrs []string
}

// TestRegister perfomes functional tests on the register endpoint via the
// Handler.
func TestRegister(t *testing.T) {
	t.Run("UsernameValidation", func(t *testing.T) {
		const (
			empty       = "Username cannot be empty."
			tooShort    = "Username cannot be shorter than 5 characters."
			tooLong     = "Username cannot be longer than 15 characters."
			invalidChar = "Username can contain only letters (a-z/A-Z) and digits (0-9)."
			digitStart  = "Username can start only with a letter (a-z/A-Z)."
		)
		for _, c := range []ValidationTestCase{
			// 1-error cases
			{name: "Empty", input: "", wantErrs: []string{empty}},
			{name: "TooShort", input: "bob1", wantErrs: []string{tooShort}},
			{name: "TooLong", input: "bobobobobobobobob", wantErrs: []string{tooLong}},
			{name: "InvalidCharacter", input: "bobob!", wantErrs: []string{invalidChar}},
			{name: "DigitStart", input: "1bobob", wantErrs: []string{digitStart}},

			// 2-error cases
			{name: "TooShort_InvalidCharacter", input: "bob!", wantErrs: []string{tooShort, invalidChar}},
			{name: "TooShort_DigitStart", input: "1bob", wantErrs: []string{tooShort, digitStart}},
			{name: "TooLong_InvalidCharacter", input: "bobobobobobobobo!", wantErrs: []string{tooLong, invalidChar}},
			{name: "TooLong_DigitStart", input: "1bobobobobobobobo", wantErrs: []string{tooLong, digitStart}},
			{name: "InvalidCharacter_DigitStart", input: "1bob!", wantErrs: []string{invalidChar, digitStart}},

			// 3-error cases
			{name: "TooShort_InvalidCharacter_DigitStart", input: "1bo!", wantErrs: []string{tooShort, invalidChar, digitStart}},
			{name: "TooLong_InvalidCharacter_DigitStart", input: "1bobobobobobobob!", wantErrs: []string{tooLong, invalidChar, digitStart}},
		} {
			t.Run(c.name, func(t *testing.T) {
				// arrange
				req, err := http.NewRequest("POST", "/register", strings.NewReader(fmt.Sprintf(`{
					"username": "%s", 
					"password": "SecureP4ss?", 
					"referrer": ""
				}`, c.input)))
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				handler := NewHandler()

				// act
				handler.ServeHTTP(w, req)

				// assert
				res := w.Result()
				gotStatusCode, wantStatusCode := res.StatusCode, http.StatusBadRequest
				if gotStatusCode != wantStatusCode {
					t.Logf("\nwant: %d\ngot: %d", http.StatusBadRequest, res.StatusCode)
					t.Fail()
				}
				resBody := &ResBody{}
				if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}
				gotErr := resBody.Errs.Username
				if !assert.EqualArr(gotErr, c.wantErrs) {
					t.Logf("\nwant: %+v\ngot: %+v", c.wantErrs, gotErr)
					t.Fail()
				}
			})
		}
	})

	t.Run("PasswordValidation", func(t *testing.T) {
		const (
			empty     = "Password cannot be empty."
			tooShort  = "Password cannot be shorter than 8 characters."
			tooLong   = "Password cannot be longer than 64 characters."
			noLower   = "Password must contain a lowercase letter (a-z)."
			noUpper   = "Password must contain an uppercase letter (A-Z)."
			noDigits  = "Password must contain a digit (0-9)."
			noSpecial = "Password must contain one of the following special characters: " +
				"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ _ ` { | } ~."
			hasSpace = "Password cannot contain spaces."
			nonASCII = "Password can contain only letters (a-z/A-Z), digits (0-9), " +
				"and the following special characters: " +
				"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ _ ` { | } ~."
		)
		for _, c := range []ValidationTestCase{
			// 1-error cases
			{name: "Empty", input: "", wantErrs: []string{empty}},
			{name: "TooShort", input: "Myp4ss!", wantErrs: []string{tooShort}},
			{name: "TooLong", input: "Myp4sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{tooLong}},
			{name: "NoLower", input: "MY4LLUPPERPASSWORD!", wantErrs: []string{noLower}},
			{name: "NoUpper", input: "my4lllowerpassword!", wantErrs: []string{noUpper}},
			{name: "NoDigits", input: "myNOdigitPASSWORD!", wantErrs: []string{noDigits}},
			{name: "NoSpecial", input: "myNOspecialP4SSWORD", wantErrs: []string{noSpecial}},
			{name: "HasSpace", input: "my SP4CED p4ssword !", wantErrs: []string{hasSpace}},
			{name: "NonASCII", input: "myNØNÅSCÎÎp4ssword!", wantErrs: []string{nonASCII}},

			// 2-error cases
			{name: "TooShort_NoLower", input: "MYP4SS!", wantErrs: []string{tooShort, noLower}},
			{name: "TooShort_NoUpper", input: "myp4ss!", wantErrs: []string{tooShort, noUpper}},
			{name: "TooShort_NoDigits", input: "MyPass!", wantErrs: []string{tooShort, noDigits}},
			{name: "TooShort_NoSpecial", input: "MyP4ssw", wantErrs: []string{tooShort, noSpecial}},
			{name: "TooShort_HasSpace", input: "My P4s!", wantErrs: []string{tooShort, hasSpace}},
			{name: "TooShort_NonASCII", input: "M¥π4ss!", wantErrs: []string{tooShort, nonASCII}},
		} {
			t.Run(c.name, func(t *testing.T) {
				// arrange
				req, err := http.NewRequest("POST", "/register", strings.NewReader(fmt.Sprintf(`{
					"username": "mynameisbob", 
					"password": "%s", 
					"referrer": ""
				}`, c.input)))
				if err != nil {
					t.Fatal(err)
				}
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()
				handler := NewHandler()

				// act
				handler.ServeHTTP(w, req)

				// assert
				res := w.Result()
				gotStatusCode, wantStatusCode := res.StatusCode, http.StatusBadRequest
				if gotStatusCode != wantStatusCode {
					t.Logf("\nwant: %d\ngot: %d", wantStatusCode, gotStatusCode)
					t.Fail()
				}
				resBody := &ResBody{}
				if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}
				gotErrs := resBody.Errs.Password
				if !assert.EqualArr(gotErrs, c.wantErrs) {
					t.Logf("\nwant: %+v\ngot: %+v", c.wantErrs, gotErrs)
					t.Fail()
				}
			})
		}
	})
}
