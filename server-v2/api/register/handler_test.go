package register

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/kxplxn/goteam/server-v2/relay"
)

func TestRegister(t *testing.T) {
	t.Run("UsernameValidation", func(t *testing.T) {
		for _, c := range []struct {
			caseName string
			username string
			wantErrs []string
		}{
			{
				caseName: "TooShort",
				username: "bob",
				wantErrs: []string{"Username cannot be shorter than 5 characters."},
			},
			{
				caseName: "TooLong",
				username: "bobobobobobobobob",
				wantErrs: []string{"Username cannot be longer than 15 characters."},
			},
			{
				caseName: "InvalidCharacter",
				username: "bobob!",
				wantErrs: []string{"Username can contain only letters and digits."},
			},
			{
				caseName: "DigitStart",
				username: "1bobob",
				wantErrs: []string{"Username can start only with a letter."},
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
				handler := NewHandler(relay.NewAPILogger())

				// act
				handler.ServeHTTP(w, req)

				// assert
				res := w.Result()
				if res.StatusCode != http.StatusBadRequest {
					t.Logf("\nwant: %d\ngot: %d", http.StatusBadRequest, res.StatusCode)
					t.Fail()
				}
				resBody := &Res{}
				if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
					t.Fatal(err)
				}
				if !cmp.Equal(resBody.Errs.Username, c.wantErrs) {
					t.Logf("\nwant: %+v\ngot: %+v", c.wantErrs, resBody.Errs.Username)
					t.Fail()
				}
			})
		}
	})
}
