//go:build itest

package itest

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"testing"

	_ "github.com/lib/pq"
)

// TestRegister tests the /register route to assert that it behaves correctly
// basesd on the request sent.
func TestRegister(t *testing.T) {
	url := serverURL + "/register"

	// Redeclare contract to not depend on server/api/register.
	type ReqBody struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	type ValidationErrs struct {
		Username []string `json:"username"`
		Password []string `json:"password"`
		Auth     string   `json:"auth"`
	}
	type ResBody struct {
		Msg  string         `json:"message"`
		Errs ValidationErrs `json:"errors"`
	}

	t.Run("ValidationErrs", func(t *testing.T) {
		t.Log("dbConnStr: " + dbConnStr)

		db, err := sql.Open("postgres", dbConnStr)
		if err != nil {
			t.Fatal(err.Error())
		}
		existingUserID := "bob123"
		_, err = db.Exec(
			`INSERT INTO app."user"(id, password) VALUES ($1, $2)`,
			existingUserID, "somepassword",
		)
		if err != nil {
			t.Fatal(err)
		}

		for _, c := range []struct {
			name         string
			username     string
			password     string
			usernameErrs []string
			passwordErrs []string
		}{
			{
				name:         "UsnEmpty,PwdEmpty",
				username:     "",
				password:     "",
				usernameErrs: []string{"Username cannot be empty."},
				passwordErrs: []string{"Password cannot be empty."},
			},
			{
				name: "UsnTooShort,UsnInvalidChar,PwdTooShort,PwdNoLower," +
					"PwdNoDigit,PwdNoSpecial",
				username: "bob!",
				password: "PASSSSS",
				usernameErrs: []string{
					"Username cannot be shorter than 5 characters.",
					"Username can contain only letters (a-z/A-Z) and digits " +
						"(0-9).",
				},
				passwordErrs: []string{
					"Password cannot be shorter than 8 characters.",
					"Password must contain a lowercase letter (a-z).",
					"Password must contain a digit (0-9).",
					"Password must contain one of the following special " +
						"characters: ! \" # $ % & ' ( ) * + , - . / : ; < = " +
						"> ? [ \\ ] ^ _ ` { | } ~.",
				},
			},
			{
				name: "UsnTooLong,UsnDigitStart,PwdTooLong,PwdNoUpper," +
					"PwdHasSpace,PwdNonASCII",
				username: "1bobobobobobobobo",
				password: "p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p" +
					"£$ 123p£$ 123p£$ 123p£",
				usernameErrs: []string{
					"Username cannot be longer than 15 characters.",
					"Username can start only with a letter (a-z/A-Z).",
				},
				passwordErrs: []string{
					"Password cannot be longer than 64 characters.",
					"Password must contain an uppercase letter (A-Z).",
					"Password cannot contain spaces.",
					"Password can contain only letters (a-z/A-Z), digits " +
						"(0-9), and the following special characters: " +
						"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ " +
						"_ ` { | } ~.",
				},
			},
			{
				name:         "UsnTaken",
				username:     existingUserID,
				password:     "Myp4ssw0rd!",
				usernameErrs: []string{"Username is already taken."},
				passwordErrs: []string{},
			},
		} {
			t.Run(c.name, func(t *testing.T) {
				reqBody, err := json.Marshal(ReqBody{
					Username: c.username,
					Password: c.password,
				})
				if err != nil {
					t.Fatal(err)
				}
				res, err := http.Post(
					url, "application/json", bytes.NewBuffer(reqBody),
				)
				if err != nil {
					t.Fatal(err)
				}

				if err := assertEqual(
					http.StatusBadRequest, res.StatusCode,
				); err != nil {
					t.Error(err)
				}

				var resBody ResBody
				if err := json.NewDecoder(res.Body).Decode(
					&resBody,
				); err != nil {
					t.Fatal(err)
				}

				if err := assertEqualArr(
					c.usernameErrs,
					resBody.Errs.Username,
				); err != nil {
					t.Error(err)
				}

				if err := assertEqualArr(
					c.passwordErrs,
					resBody.Errs.Password,
				); err != nil {
					t.Error(err)
				}
			})
		}
	})
}
