//go:build itest

package itest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	registerAPI "server/api/register"
	"server/assert"
)

const (
	usnEmpty       = "Username cannot be empty."
	usnTooShort    = "Username cannot be shorter than 5 characters."
	usnTooLong     = "Username cannot be longer than 15 characters."
	usnInvalidChar = "Username can contain only letters (a-z/A-Z) and digits (0-9)."
	usnDigitStart  = "Username can start only with a letter (a-z/A-Z)."

	pwdEmpty     = "Password cannot be empty."
	pwdTooShort  = "Password cannot be shorter than 8 characters."
	pwdTooLong   = "Password cannot be longer than 64 characters."
	pwdNoLower   = "Password must contain a lowercase letter (a-z)."
	pwdNoUpper   = "Password must contain an uppercase letter (A-Z)."
	pwdNoDigit   = "Password must contain a digit (0-9)."
	pwdNoSpecial = "Password must contain one of the following special characters: " +
		"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ _ ` { | } ~."
	pwdHasSpace = "Password cannot contain spaces."
	pwdNonASCII = "Password can contain only letters (a-z/A-Z), digits (0-9), " +
		"and the following special characters: " +
		"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ _ ` { | } ~."
)

// TestRegister tests the /register route to assert that it behaves correctly
// basesd on the request sent.
func TestRegister(t *testing.T) {
	url := serverURL + "/register"

	t.Run("ValidationErrs", func(t *testing.T) {
		t.Run("UsnEmpty,PwdEmpty", func(t *testing.T) {
			reqBody := registerAPI.ReqBody{
				Username: "",
				Password: "",
			}
			reqBodyJSON, err := json.Marshal(reqBody)
			if err != nil {
				t.Fatal(err)
			}
			res, err := http.Post(
				url, "application/json", bytes.NewBuffer(reqBodyJSON),
			)
			if err != nil {
				t.Fatal(err)
			}

			if err := assert.Equal(
				http.StatusBadRequest, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			resBody := registerAPI.ResBody{}
			if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
				t.Fatal(err)
			}

			if err := assert.EqualArr(
				[]string{usnEmpty},
				resBody.ValidationErrs.Username,
			); err != nil {
				t.Error(err)
			}

			if err := assert.EqualArr(
				[]string{pwdEmpty},
				resBody.ValidationErrs.Password,
			); err != nil {
				t.Error(err)
			}
		})
	})
}
