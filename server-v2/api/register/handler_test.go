package register

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegister(t *testing.T) {
	t.Run("UsernameValidation", func(t *testing.T) {
		for _, c := range []struct {
			caseName string
			username string
			wantErr  string
		}{
			{
				caseName: "Empty",
				username: "",
				wantErr:  "Username cannot be empty.",
			},
			{
				caseName: "TooShort",
				username: "bob",
				wantErr:  "Username cannot be shorter than 5 characters.",
			},
			{
				caseName: "TooLong",
				username: "bobobobobobobobob",
				wantErr:  "Username cannot be longer than 15 characters.",
			},
			{
				caseName: "InvalidCharacter",
				username: "bobob!",
				wantErr:  "Username can contain only letters and digits.",
			},
			{
				caseName: "DigitStart",
				username: "1bobob",
				wantErr:  "Username can start only with a letter.",
			},
		} {
			t.Run(c.caseName, func(t *testing.T) {
				// arrange
				req, err := http.NewRequest("POST", "/register", strings.NewReader(fmt.Sprintf(`{
					"username": "%s", 
					"password": "securepass1!", 
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
				if res.StatusCode != http.StatusBadRequest {
					t.Logf("\nwant: %d\ngot: %d", http.StatusBadRequest, res.StatusCode)
					t.Fail()
				}
				resBody := &ResBody{}
				if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}
				gotErr := resBody.Errs.Username
				if gotErr != c.wantErr {
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
			wantErr  string
		}{
			{
				caseName: "Empty",
				password: "",
				wantErr:  "Password cannot be empty.",
			},
			{
				caseName: "TooShort",
				password: "mypassw",
				wantErr:  "Password cannot be shorter than 5 characters.",
			},
			{
				caseName: "TooLong",
				password: "mypasswordwhichislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErr:  "Password cannot be longer than 64 characters.",
			},
			{
				caseName: "NoLowercase",
				password: "MYALLUPPERPASSWORD",
				wantErr:  "Password must contain a lowercase letter (a-z).",
			},
			{
				caseName: "NoUppercase",
				password: "myalllowerpassword",
				wantErr:  "Password must contain an uppercase letter (A-Z).",
			},
			{
				caseName: "NoDigits",
				password: "myNOdigitPASSWORD",
				wantErr:  "Password must contain a digit (0-9).",
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
				if res.StatusCode != http.StatusBadRequest {
					t.Logf("\nwant: %d\ngot: %d", http.StatusBadRequest, res.StatusCode)
					t.Fail()
				}
				resBody := &ResBody{}
				if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}
				gotErr := resBody.Errs.Password
				if gotErr != c.wantErr {
					t.Logf("\nwant: %+v\ngot: %+v", c.wantErr, gotErr)
					t.Fail()
				}
			})
		}
	})
}
