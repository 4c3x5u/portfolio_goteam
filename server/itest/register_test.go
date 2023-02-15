//go:build itest

package itest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	registerAPI "server/api/register"
	"server/assert"
	"server/auth"
	"server/db"
	"server/log"

	"golang.org/x/crypto/bcrypt"
)

func TestRegister(t *testing.T) {
	sut := registerAPI.NewHandler(
		registerAPI.NewValidator(
			registerAPI.NewUsernameValidator(),
			registerAPI.NewPasswordValidator(),
		),
		db.NewUserSelector(dbConnPool),
		registerAPI.NewPasswordHasher(),
		db.NewUserInserter(dbConnPool),
		auth.NewJWTGenerator(jwtKey),
		log.NewAppLogger(),
	)

	// To validate the JWT returned in the HTTP response.
	jwtValidator := auth.NewJWTValidator(jwtKey)

	for _, c := range []struct {
		name             string
		username         string
		password         string
		wantStatusCode   int
		wantUsernameErrs []string
		wantPasswordErrs []string
	}{
		{
			name:             "UsnEmpty,PwdEmpty",
			username:         "",
			password:         "",
			wantStatusCode:   http.StatusBadRequest,
			wantUsernameErrs: []string{"Username cannot be empty."},
			wantPasswordErrs: []string{"Password cannot be empty."},
		},
		{
			name: "UsnTooShort,UsnInvalidChar,PwdTooShort,PwdNoLower," +
				"PwdNoDigit,PwdNoSpecial",
			username:       "bob!",
			password:       "PASSSSS",
			wantStatusCode: http.StatusBadRequest,
			wantUsernameErrs: []string{
				"Username cannot be shorter than 5 characters.",
				"Username can contain only letters (a-z/A-Z) and digits " +
					"(0-9).",
			},
			wantPasswordErrs: []string{
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
			wantStatusCode: http.StatusBadRequest,
			wantUsernameErrs: []string{
				"Username cannot be longer than 15 characters.",
				"Username can start only with a letter (a-z/A-Z).",
			},
			wantPasswordErrs: []string{
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
			name:             "UsnTaken",
			username:         "bob123",
			password:         "Myp4ssw0rd!",
			wantStatusCode:   http.StatusBadRequest,
			wantUsernameErrs: []string{"Username is already taken."},
			wantPasswordErrs: []string{},
		},
		{
			name:             "Success",
			username:         "bob321",
			password:         "Myp4ssw0rd!",
			wantStatusCode:   http.StatusOK,
			wantUsernameErrs: []string{},
			wantPasswordErrs: []string{},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			reqBody, err := json.Marshal(registerAPI.ReqBody{
				Username: c.username,
				Password: c.password,
			})
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPost, "", bytes.NewReader(reqBody),
			)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			sut.ServeHTTP(w, req)

			res := w.Result()

			if err = assert.Equal(
				c.wantStatusCode, res.StatusCode,
			); err != nil {
				t.Error(err)
			}

			switch c.wantStatusCode {
			case http.StatusBadRequest:
				// assert that the correct validation errors are returned
				var resBody registerAPI.ResBody
				if err = json.NewDecoder(res.Body).Decode(
					&resBody,
				); err != nil {
					t.Fatal(err)
				}
				if err = assert.EqualArr(
					c.wantUsernameErrs,
					resBody.ValidationErrs.Username,
				); err != nil {
					t.Error(err)
				}
				if err = assert.EqualArr(
					c.wantPasswordErrs,
					resBody.ValidationErrs.Password,
				); err != nil {
					t.Error(err)
				}
			case http.StatusOK:
				// assert that a new user is inserted into the database with
				// the correct credentials
				var userID, password string
				err = dbConnPool.QueryRow(
					`SELECT id, password FROM app."user" WHERE id = $1`,
					c.username,
				).Scan(&userID, &password)
				if err != nil {
					t.Fatal(err)
				}
				if err = assert.Equal(c.username, userID); err != nil {
					t.Error(err)
				}
				if err = bcrypt.CompareHashAndPassword(
					[]byte(password), []byte(c.password),
				); err != nil {
					t.Error(err)
				}

				// assert that the returned JWT is valid and has the correct
				// subject
				var jwt string
				for _, ck := range res.Cookies() {
					if ck.Name == "auth-token" {
						jwt = ck.Value
					}
				}
				sub := jwtValidator.Validate(jwt)
				if err = assert.Equal(c.username, sub); err != nil {
					t.Error(err)
				}
			}
		})
	}
}
