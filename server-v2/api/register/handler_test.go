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

func TestRegister(t *testing.T) {
	t.Run("UsernameValidation", func(t *testing.T) {
		for _, c := range []struct {
			caseName string
			username string
			wantErr  []string
		}{
			{
				caseName: "Empty",
				username: "",
				wantErr:  []string{"Username cannot be empty."},
			},
			{
				caseName: "TooShort",
				username: "bob1",
				wantErr:  []string{"Username cannot be shorter than 5 characters."},
			},
			{
				caseName: "TooLong",
				username: "bobobobobobobobob",
				wantErr:  []string{"Username cannot be longer than 15 characters."},
			},
			{
				caseName: "InvalidCharacter",
				username: "bobob!",
				wantErr:  []string{"Username can contain only letters (a-z/A-Z) and digits (0-9)."},
			},
			{
				caseName: "DigitStart",
				username: "1bobob",
				wantErr:  []string{"Username can start only with a letter (a-z/A-Z)."},
			},
			{
				caseName: "TooShort_InvalidCharacter",
				username: "bob!",
				wantErr: []string{
					"Username cannot be shorter than 5 characters.",
					"Username can contain only letters (a-z/A-Z) and digits (0-9).",
				},
			},
			{
				caseName: "TooShort_DigitStart",
				username: "1bob",
				wantErr: []string{
					"Username cannot be shorter than 5 characters.",
					"Username can start only with a letter (a-z/A-Z).",
				},
			},
			{
				caseName: "TooLong_InvalidCharacter",
				username: "bobobobobobobobo!",
				wantErr: []string{
					"Username cannot be longer than 15 characters.",
					"Username can contain only letters (a-z/A-Z) and digits (0-9).",
				},
			},
			{
				caseName: "TooLong_DigitStart",
				username: "1bobobobobobobobo",
				wantErr: []string{
					"Username cannot be longer than 15 characters.",
					"Username can start only with a letter (a-z/A-Z).",
				},
			},
			{
				caseName: "InvalidCharacter_DigitStart",
				username: "1bob!",
				wantErr: []string{
					"Username can contain only letters (a-z/A-Z) and digits (0-9).",
					"Username can start only with a letter (a-z/A-Z).",
				},
			},
			{
				caseName: "TooShort_InvalidCharacter_DigitStart",
				username: "1bo!",
				wantErr: []string{
					"Username cannot be shorter than 5 characters.",
					"Username can contain only letters (a-z/A-Z) and digits (0-9).",
					"Username can start only with a letter (a-z/A-Z).",
				},
			},
			{
				caseName: "TooLong_InvalidCharacter_DigitStart",
				username: "1bobobobobobobob!",
				wantErr: []string{
					"Username cannot be longer than 15 characters.",
					"Username can contain only letters (a-z/A-Z) and digits (0-9).",
					"Username can start only with a letter (a-z/A-Z).",
				},
			},
		} {
			t.Run(c.caseName, func(t *testing.T) {
				// arrange
				req, err := http.NewRequest("POST", "/register", strings.NewReader(fmt.Sprintf(`{
					"username": "%s", 
					"password": "SecureP4ss?", 
					"referrer": ""
				}`, c.username)))
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
				if !assert.EqualArr(gotErr, c.wantErr) {
					t.Logf("\nwant: %+v\ngot: %+v", c.wantErr, gotErr)
					t.Fail()
				}
			})
		}
	})

	t.Run("PasswordValidation", func(t *testing.T) {
		for _, c := range []struct {
			caseName string
			password string
			wantErr  []string
		}{
			{
				caseName: "Empty",
				password: "",
				wantErr:  []string{"Password cannot be empty."},
			},
			{
				caseName: "TooShort",
				password: "mypassw",
				wantErr:  []string{"Password cannot be shorter than 5 characters."},
			},
			{
				caseName: "TooLong",
				password: "mypasswordwhichislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErr:  []string{"Password cannot be longer than 64 characters."},
			},
			{
				caseName: "NoLowercase",
				password: "MYALLUPPERPASSWORD",
				wantErr:  []string{"Password must contain a lowercase letter (a-z)."},
			},
			{
				caseName: "NoUppercase",
				password: "myalllowerpassword",
				wantErr:  []string{"Password must contain an uppercase letter (A-Z)."},
			},
			{
				caseName: "NoDigits",
				password: "myNOdigitPASSWORD",
				wantErr:  []string{"Password must contain a digit (0-9)."},
			},
			{
				caseName: "NoSymbols",
				password: "myNOsymbolP4SSWORD",
				wantErr: []string{
					"Password must contain one of the following special characters: " +
						"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ _ ` { | } ~.",
				},
			},
			{
				caseName: "HasSpaces",
				password: "my SP4CED p4ssword",
				wantErr:  []string{"Password cannot contain spaces."},
			},
			{
				caseName: "NonASCII",
				password: "myNØNÅSCÎÎp4ssword",
				wantErr: []string{
					"Password can contain only letters (a-z/A-Z), digits (0-9), " +
						"and the following special characters: " +
						"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ _ ` { | } ~.",
				},
			},
		} {
			t.Run(c.caseName, func(t *testing.T) {
				// arrange
				req, err := http.NewRequest("POST", "/register", strings.NewReader(fmt.Sprintf(`{
					"username": "mynameisbob", 
					"password": "%s", 
					"referrer": ""
				}`, c.password)))
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
				gotErr := resBody.Errs.Password
				if !assert.EqualArr(gotErr, c.wantErr) {
					t.Logf("\nwant: %+v\ngot: %+v", c.wantErr, gotErr)
					t.Fail()
				}
			})
		}
	})
}
