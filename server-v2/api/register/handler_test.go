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
			noDigit   = "Password must contain a digit (0-9)."
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
			{
				name:     "TooLong",
				input:    "Myp4sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong},
			},
			{name: "NoLower", input: "MY4LLUPPERPASSWORD!", wantErrs: []string{noLower}},
			{name: "NoUpper", input: "my4lllowerpassword!", wantErrs: []string{noUpper}},
			{name: "NoDigit", input: "myNOdigitPASSWORD!", wantErrs: []string{noDigit}},
			{name: "NoSpecial", input: "myNOspecialP4SSWORD", wantErrs: []string{noSpecial}},
			{name: "HasSpace", input: "my SP4CED p4ssword !", wantErrs: []string{hasSpace}},
			{name: "NonASCII", input: "myNØNÅSCÎÎp4ssword!", wantErrs: []string{nonASCII}},

			// 2-error cases
			{name: "TooShort_NoLower", input: "MYP4SS!", wantErrs: []string{tooShort, noLower}},
			{name: "TooShort_NoUpper", input: "myp4ss!", wantErrs: []string{tooShort, noUpper}},
			{name: "TooShort_NoDigit", input: "MyPass!", wantErrs: []string{tooShort, noDigit}},
			{name: "TooShort_NoSpecial", input: "MyP4ssw", wantErrs: []string{tooShort, noSpecial}},
			{name: "TooShort_HasSpace", input: "My P4s!", wantErrs: []string{tooShort, hasSpace}},
			{name: "TooShort_NonASCII", input: "M¥P4s!2", wantErrs: []string{tooShort, nonASCII}},
			{
				name:     "TooLong_NoLower",
				input:    "MYP4SSWORDWH!CHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH",
				wantErrs: []string{tooLong, noLower},
			},
			{
				name:     "TooLong_NoUpper",
				input:    "myp4sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong, noUpper},
			},
			{
				name:     "TooLong_NoDigit",
				input:    "Mypasswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong, noDigit},
			},
			{
				name:     "TooLong_NoSpecial",
				input:    "Myp4sswordwhichislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong, noSpecial},
			},
			{
				name:     "TooLong_HasSpace",
				input:    "Myp4sswo   rdwh!chislongandimeanreallylongforsomereasonohiknowwhy",
				wantErrs: []string{tooLong, hasSpace},
			},
			{
				name:     "TooLong_NonASCII",
				input:    "Myp4££wordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong, nonASCII},
			},
			{name: "NoLower_NoUpper", input: "4444!!!!", wantErrs: []string{noLower, noUpper}},
			{name: "NoLower_NoDigit", input: "MYP@SSW!", wantErrs: []string{noLower, noDigit}},
			{name: "NoLower_NoSpecial", input: "MYP4SSW1", wantErrs: []string{noLower, noSpecial}},
			{name: "NoLower_HasSpace", input: "MYP4SS !", wantErrs: []string{noLower, hasSpace}},
			{name: "NoLower_NonASCII", input: "MYP4££W!", wantErrs: []string{noLower, nonASCII}},
			{name: "NoUpper_NoDigit", input: "myp@ssw!", wantErrs: []string{noUpper, noDigit}},
			{name: "NoUpper_NoSpecial", input: "myp4ssw1", wantErrs: []string{noUpper, noSpecial}},
			{name: "NoUpper_HasSpace", input: "myp4ss !", wantErrs: []string{noUpper, hasSpace}},
			{name: "NoUpper_NonASCII", input: "myp4££w!", wantErrs: []string{noUpper, nonASCII}},
			{name: "NoDigit_NoSpecial", input: "MyPasswd", wantErrs: []string{noDigit, noSpecial}},
			{name: "NoDigit_HasSpace", input: "MyPass !", wantErrs: []string{noDigit, hasSpace}},
			{name: "NoDigit_NonASCII", input: "MyPa££w!", wantErrs: []string{noDigit, nonASCII}},
			{name: "NoSpecial_HasSpace", input: "My  P4ss", wantErrs: []string{noSpecial, hasSpace}},
			{name: "NoSpecial_NonASCII", input: "MyPa££w1", wantErrs: []string{noSpecial, nonASCII}},
			{name: "HasSpace_NonASCII", input: "MyP4££ !", wantErrs: []string{hasSpace, nonASCII}},

			// 3-error cases
			{name: "TooShort_NoLower_NoUpper", input: "1421!@$", wantErrs: []string{tooShort, noLower, noUpper}},
			{name: "TooShort_NoLower_NoDigit", input: "PASS!@$", wantErrs: []string{tooShort, noLower, noDigit}},
			{name: "TooShort_NoLower_NoSpecial", input: "PASS123", wantErrs: []string{tooShort, noLower, noSpecial}},
			{name: "TooShort_NoLower_HasSpace", input: "PA$ 123", wantErrs: []string{tooShort, noLower, hasSpace}},
			{name: "TooShort_NoLower_NonASCII", input: "PA$£123", wantErrs: []string{tooShort, noLower, nonASCII}},
			{name: "TooShort_NoUpper_NoDigit", input: "pass$$$", wantErrs: []string{tooShort, noUpper, noDigit}},
			{name: "TooShort_NoUpper_NoSpecial", input: "pass123", wantErrs: []string{tooShort, noUpper, noSpecial}},
			{name: "TooShort_NoUpper_HasSpace", input: "pa$ 123", wantErrs: []string{tooShort, noUpper, hasSpace}},
			{name: "TooShort_NoUpper_NonASCII", input: "pa$£123", wantErrs: []string{tooShort, noUpper, nonASCII}},
			{name: "TooShort_NoDigit_NoSpecial", input: "Passwor", wantErrs: []string{tooShort, noDigit, noSpecial}},
			{name: "TooShort_NoDigit_HasSpace", input: "Pa$$ wo", wantErrs: []string{tooShort, noDigit, hasSpace}},
			{name: "TooShort_NoDigit_NonASCII", input: "Pa$$£wo", wantErrs: []string{tooShort, noDigit, nonASCII}},
			{name: "TooShort_NoSpecial_HasSpace", input: "Pa55 wo", wantErrs: []string{tooShort, noSpecial, hasSpace}},
			{name: "TooShort_NoSpecial_NonASCII", input: "Pa55£wo", wantErrs: []string{tooShort, noSpecial, nonASCII}},
			{name: "TooShort_HasSpace_NonASCII", input: "P4$ £wo", wantErrs: []string{tooShort, hasSpace, nonASCII}},
			{
				name:     "TooLong_NoLower_NoUpper",
				input:    "111422222222!3333333333333333333333333333333333333333333333333333",
				wantErrs: []string{tooLong, noLower, noUpper},
			},
			{
				name:     "TooLong_NoLower_NoDigit",
				input:    "MYPASSWORDWH!CHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH",
				wantErrs: []string{tooLong, noLower, noDigit},
			},
			{
				name:     "TooLong_NoLower_NoSpecial",
				input:    "MYP4SSWORDWHICHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH",
				wantErrs: []string{tooLong, noLower, noSpecial},
			},
			{
				name:     "TooLong_NoLower_HasSpace",
				input:    "MYP4SS    WH!CHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH",
				wantErrs: []string{tooLong, noLower, hasSpace},
			},
			{
				name:     "TooLong_NoLower_NonASCII",
				input:    "£YP4SSWORDWH!CHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH",
				wantErrs: []string{tooLong, noLower, nonASCII},
			},
			{
				name:     "TooLong_NoUpper_NoDigit",
				input:    "mypasswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong, noUpper, noDigit},
			},
			{
				name:     "TooLong_NoUpper_NoSpecial",
				input:    "myp4sswordwhichislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong, noUpper, noSpecial},
			},
			{
				name:     "TooLong_NoUpper_HasSpace",
				input:    "myp4ss    wh!chislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong, noUpper, hasSpace},
			},
			{
				name:     "TooLong_NoUpper_NonASCII",
				input:    "£yp4sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong, noUpper, nonASCII},
			},
			{
				name:     "TooLong_NoDigit_NoSpecial",
				input:    "Mypasswordwhichislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong, noDigit, noSpecial},
			},
			{
				name:     "TooLong_NoDigit_HasSpace",
				input:    "Mypass    wh!chislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong, noDigit, hasSpace},
			},
			{
				name:     "TooLong_NoDigit_NonASCII",
				input:    "Myp£sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong, noDigit, nonASCII},
			},
			{
				name:     "TooLong_NoSpecial_HasSpace",
				input:    "Myp4ss    whichislongandimeanreallylongforsomereasonohiknowwhytbh",
				wantErrs: []string{tooLong, noSpecial, hasSpace},
			},
			{
				name:     "TooLong_HasSpace_NonASCII",
				input:    "Myp4ssw£   rdwh!chislongandimeanreallylongforsomereasonohiknowwhy",
				wantErrs: []string{tooLong, hasSpace, nonASCII},
			},
			{name: "NoLower_NoUpper_NoDigit", input: "!!!!!!!!", wantErrs: []string{noLower, noUpper, noDigit}},
			{name: "NoLower_NoUpper_NoSpecial", input: "33333333", wantErrs: []string{noLower, noUpper, noSpecial}},
			{name: "NoLower_NoUpper_HasSpace", input: "444  !!!", wantErrs: []string{noLower, noUpper, hasSpace}},
			{name: "NoLower_NoUpper_NonASCII", input: "£££444!!", wantErrs: []string{noLower, noUpper, nonASCII}},
			{name: "NoLower_NoDigit_NoSpecial", input: "MYPASSWO", wantErrs: []string{noLower, noDigit, noSpecial}},
			{name: "NoLower_NoDigit_HasSpace", input: "MYP@SS !", wantErrs: []string{noLower, noDigit, hasSpace}},
			{name: "NoLower_NoDigit_NonASCII", input: "M£P@SSW!", wantErrs: []string{noLower, noDigit, nonASCII}},
			{name: "NoLower_NoSpecial_HasSpace", input: "MYP4  W1", wantErrs: []string{noLower, noSpecial, hasSpace}},
			{name: "NoLower_NoSpecial_NonASCII", input: "M£P4SSW1", wantErrs: []string{noLower, noSpecial, nonASCII}},
			{name: "NoLower_HasSpace_NonASCII", input: "M£P4SS !", wantErrs: []string{noLower, hasSpace, nonASCII}},
			{name: "NoUpper_NoDigit_NoSpecial", input: "mypasswo", wantErrs: []string{noUpper, noDigit, noSpecial}},
			{name: "NoUpper_NoDigit_HasSpace", input: "myp@ss !", wantErrs: []string{noUpper, noDigit, hasSpace}},
			{name: "NoUpper_NoDigit_NonASCII", input: "m£p@ssw!", wantErrs: []string{noUpper, noDigit, nonASCII}},
			{name: "NoUpper_NoSpecial_HasSpace", input: "myp4ss 1", wantErrs: []string{noUpper, noSpecial, hasSpace}},
			{name: "NoUpper_NoSpecial_NonASCII", input: "m£p4ssw1", wantErrs: []string{noUpper, noSpecial, nonASCII}},
			{name: "NoUpper_HasSpace_NonASCII", input: "m£p4ss !", wantErrs: []string{noUpper, hasSpace, nonASCII}},
			{name: "NoDigit_NoSpecial_HasSpace", input: "MyPass o", wantErrs: []string{noDigit, noSpecial, hasSpace}},
			{name: "NoDigit_NoSpecial_NonASCII", input: "M£Passwd", wantErrs: []string{noDigit, noSpecial, nonASCII}},
			{name: "NoDigit_HasSpace_NonASCII", input: "M£Pass !", wantErrs: []string{noDigit, hasSpace, nonASCII}},
			{name: "NoSpecial_HasSpace_NonASCII", input: "M£  P4ss", wantErrs: []string{noSpecial, hasSpace, nonASCII}},

			// 4-error cases
			{name: "TooShort_NoLower_NoUpper_NoDigits", input: "!@$!@$!", wantErrs: []string{tooShort, noLower, noUpper, noDigit}},
			{name: "TooShort_NoLower_NoUpper_NoSpecial", input: "1421111", wantErrs: []string{tooShort, noLower, noUpper, noSpecial}},
			{name: "TooShort_NoLower_NoUpper_HasSpace", input: "142 !@$", wantErrs: []string{tooShort, noLower, noUpper, hasSpace}},
			{name: "TooShort_NoLower_NoUpper_NonASCII", input: "14£1!@$", wantErrs: []string{tooShort, noLower, noUpper, nonASCII}},
			{name: "TooShort_NoLower_NoDigit_NoSpecial", input: "PASSSSS", wantErrs: []string{tooShort, noLower, noDigit, noSpecial}},
			{name: "TooShort_NoLower_NoDigit_HasSpace", input: "PAS !@$", wantErrs: []string{tooShort, noLower, noDigit, hasSpace}},
			{name: "TooShort_NoLower_NoDigit_NonASCII", input: "P£SS!@$", wantErrs: []string{tooShort, noLower, noDigit, nonASCII}},
			{name: "TooShort_NoLower_NoSpecial_HasSpace", input: "PAS 123", wantErrs: []string{tooShort, noLower, noSpecial, hasSpace}},
			{name: "TooShort_NoLower_NoSpecial_NonASCII", input: "P£SS123", wantErrs: []string{tooShort, noLower, noSpecial, nonASCII}},
			{name: "TooShort_NoLower_HasSpace_NonASCII", input: "P£$ 123", wantErrs: []string{tooShort, noLower, hasSpace, nonASCII}},
			{name: "TooShort_NoUpper_NoDigit_NoSpecial", input: "passsss", wantErrs: []string{tooShort, noUpper, noDigit, noSpecial}},
			{name: "TooShort_NoUpper_NoDigit_HasSpace", input: "pas $$$", wantErrs: []string{tooShort, noUpper, noDigit, hasSpace}},
			{name: "TooShort_NoUpper_NoDigit_NonASCII", input: "p£ss$$$", wantErrs: []string{tooShort, noUpper, noDigit, nonASCII}},
			{name: "TooShort_NoUpper_NoSpecial_HasSpace", input: "pas 123", wantErrs: []string{tooShort, noUpper, noSpecial, hasSpace}},
			{name: "TooShort_NoUpper_NoSpecial_NonASCII", input: "p£ss123", wantErrs: []string{tooShort, noUpper, noSpecial, nonASCII}},
			{name: "TooShort_NoUpper_HasSpace_NonASCII", input: "p£$ 123", wantErrs: []string{tooShort, noUpper, hasSpace, nonASCII}},
			{name: "TooShort_NoDigit_NoSpecial_HasSpace", input: "Pas wor", wantErrs: []string{tooShort, noDigit, noSpecial, hasSpace}},
			{name: "TooShort_NoDigit_NoSpecial_NonASCII", input: "P£sswor", wantErrs: []string{tooShort, noDigit, noSpecial, nonASCII}},
			{name: "TooShort_NoDigit_HasSpace_NonASCII", input: "P£$$ wo", wantErrs: []string{tooShort, noDigit, hasSpace, nonASCII}},
			{name: "TooShort_NoSpecial_HasSpace_NonASCII", input: "P£55 wo", wantErrs: []string{tooShort, noSpecial, hasSpace, nonASCII}},
			{
				name:     "TooLong_NoLower_NoUpper_NoDigits",
				input:    "!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@",
				wantErrs: []string{tooLong, noLower, noUpper, noDigit},
			},
			{
				name:     "TooLong_NoLower_NoUpper_NoSpecial",
				input:    "14211111421111142111114211111421111142111114211114211111421111142",
				wantErrs: []string{tooLong, noLower, noUpper, noSpecial},
			},
			{
				name:     "TooLong_NoLower_NoUpper_HasSpace",
				input:    "142 !@$142 !@$142 !@$142 !@$142 !@$142 !@$142 !@$142 !@$142 !@$14",
				wantErrs: []string{tooLong, noLower, noUpper, hasSpace},
			},
			{
				name:     "TooLong_NoLower_NoUpper_NonASCII",
				input:    "14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14",
				wantErrs: []string{tooLong, noLower, noUpper, nonASCII},
			},
			{
				name:     "TooLong_NoLower_NoDigit_NoSpecial",
				input:    "PASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPA",
				wantErrs: []string{tooLong, noLower, noDigit, noSpecial},
			},
			{
				name:     "TooLong_NoLower_NoDigit_HasSpace",
				input:    "PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PA",
				wantErrs: []string{tooLong, noLower, noDigit, hasSpace},
			},
			{
				name:     "TooLong_NoLower_NoDigit_NonASCII",
				input:    "P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£",
				wantErrs: []string{tooLong, noLower, noDigit, nonASCII},
			},
			{
				name:     "TooLong_NoLower_NoSpecial_HasSpace",
				input:    "PAS 123PAS 123PAS 123PAS 123PAS 123PAS 123PAS 123PAS 123PAS 123PA",
				wantErrs: []string{tooLong, noLower, noSpecial, hasSpace},
			},
			{
				name:     "TooLong_NoLower_NoSpecial_NonASCII",
				input:    "P£SS123P£SS123P£SS123P£SS123P£SS123P£SS123P£SS123P£SS123P£SS123P£",
				wantErrs: []string{tooLong, noLower, noSpecial, nonASCII},
			},
			{
				name:     "TooLong_NoLower_HasSpace_NonASCII",
				input:    "P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£",
				wantErrs: []string{tooLong, noLower, hasSpace, nonASCII},
			},
			{
				name:     "TooLong_NoUpper_NoDigit_NoSpecial",
				input:    "passssspassssspassssspassssspassssspassssspassssspassssspassssspa",
				wantErrs: []string{tooLong, noUpper, noDigit, noSpecial},
			},
			{
				name:     "TooLong_NoUpper_NoDigit_HasSpace",
				input:    "pas $$$pas $$$pas $$$pas $$$pas $$$pas $$$pas $$$pas $$$pas $$$pa",
				wantErrs: []string{tooLong, noUpper, noDigit, hasSpace},
			},
			{
				name:     "TooLong_NoUpper_NoDigit_NonASCII",
				input:    "p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£",
				wantErrs: []string{tooLong, noUpper, noDigit, nonASCII},
			},
			{
				name:     "TooLong_NoUpper_NoSpecial_HasSpace",
				input:    "pas 123pas 123pas 123pas 123pas 123pas 123pas 123pas 123pas 123pa",
				wantErrs: []string{tooLong, noUpper, noSpecial, hasSpace},
			},
			{
				name:     "TooLong_NoUpper_NoSpecial_NonASCII",
				input:    "p£ss123p£ss123p£ss123p£ss123p£ss123p£ss123p£ss123p£ss123p£p£ss123",
				wantErrs: []string{tooLong, noUpper, noSpecial, nonASCII},
			},
			{
				name:     "TooLong_NoUpper_HasSpace_NonASCII",
				input:    "p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£",
				wantErrs: []string{tooLong, noUpper, hasSpace, nonASCII},
			},
			{
				name:     "TooLong_NoDigit_NoSpecial_HasSpace",
				input:    "Pas worPas worPas worPas worPas worPas worPas worPas worPas worPa",
				wantErrs: []string{tooLong, noDigit, noSpecial, hasSpace},
			},
			{
				name:     "TooLong_NoDigit_NoSpecial_NonASCII",
				input:    "P£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£",
				wantErrs: []string{tooLong, noDigit, noSpecial, nonASCII},
			},
			{
				name:     "TooLong_NoDigit_HasSpace_NonASCII",
				input:    "P£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£",
				wantErrs: []string{tooLong, noDigit, hasSpace, nonASCII},
			},
			{
				name:     "TooLong_NoSpecial_HasSpace_NonASCII",
				input:    "P£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP£",
				wantErrs: []string{tooLong, noSpecial, hasSpace, nonASCII},
			},
			{name: "NoLower_NoUpper_NoDigit_HasSpace", input: "!!!  !!!", wantErrs: []string{noLower, noUpper, noDigit, hasSpace}},
			{name: "NoLower_NoUpper_NoDigit_NonASCII", input: "!!!££!!!", wantErrs: []string{noLower, noUpper, noDigit, nonASCII}},
			{name: "NoLower_NoUpper_NoSpecial_HasSpace", input: "333  333", wantErrs: []string{noLower, noUpper, noSpecial, hasSpace}},
			{name: "NoLower_NoUpper_NoSpecial_NonASCII", input: "333££333", wantErrs: []string{noLower, noUpper, noSpecial, nonASCII}},
			{name: "NoLower_NoUpper_HasSpace_NonASCII", input: "4£4  !!!", wantErrs: []string{noLower, noUpper, hasSpace, nonASCII}},
			{name: "NoLower_NoDigit_NoSpecial_HasSpace", input: "MYP  SWO", wantErrs: []string{noLower, noDigit, noSpecial, hasSpace}},
			{name: "NoLower_NoDigit_NoSpecial_NonASCII", input: "MYP££SWO", wantErrs: []string{noLower, noDigit, noSpecial, nonASCII}},
			{name: "NoLower_NoDigit_HasSpace_NonASCII", input: "M£P@SS !", wantErrs: []string{noLower, noDigit, hasSpace, nonASCII}},
			{name: "NoLower_NoSpecial_HasSpace_NonASCII", input: "M£P4  W1", wantErrs: []string{noLower, noSpecial, hasSpace, nonASCII}},
			{name: "NoUpper_NoDigit_NoSpecial_HasSpace", input: "myp  swo", wantErrs: []string{noUpper, noDigit, noSpecial, hasSpace}},
			{name: "NoUpper_NoDigit_NoSpecial_NonASCII", input: "myp££swo", wantErrs: []string{noUpper, noDigit, noSpecial, nonASCII}},
			{name: "NoUpper_NoDigit_HasSpace_NonASCII", input: "m£p@ss !", wantErrs: []string{noUpper, noDigit, hasSpace, nonASCII}},
			{name: "NoUpper_NoSpecial_HasSpace_NonASCII", input: "m£p4  w1", wantErrs: []string{noUpper, noSpecial, hasSpace, nonASCII}},
			{name: "NoDigit_NoSpecial_HasSpace_NonASCII", input: "MyP£ss o", wantErrs: []string{noDigit, noSpecial, hasSpace, nonASCII}},
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
