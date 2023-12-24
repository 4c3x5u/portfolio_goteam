//go:build utest

package registerapi

import (
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
)

// These constants are declared separately from the strings that are actually
// used in validators so that the corresponding assertions fail when the strings
// used are accidentally edited.
const (
	idEmpty        = "Username cannot be empty."
	idTooShort     = "Username cannot be shorter than 5 characters."
	idTooLong      = "Username cannot be longer than 15 characters."
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

// TestUserValidator tests the UserValidator's Validate method to ensure that it returns
// whatever error is returned to it by UsernameValidator and PasswordValidator.
func TestUserValidator(t *testing.T) {
	fakeIDValidator := &fakeStringValidator{}
	fakePasswordValidator := &fakeStringValidator{}

	sut := NewUserValidator(fakeIDValidator, fakePasswordValidator)

	for _, c := range []struct {
		name         string
		reqBody      PostReq
		usernameErrs []string
		passwordErrs []string
	}{
		{
			name:         "UsnEmpty,PwdEmpty",
			reqBody:      PostReq{Username: "", Password: ""},
			usernameErrs: []string{idEmpty},
			passwordErrs: []string{pwdEmpty},
		},
		{
			name:         "UsnTooShort,UsnInvalidChar,PwdEmpty",
			reqBody:      PostReq{Username: "bob!", Password: "myNØNÅSCÎÎp4ssword!"},
			usernameErrs: []string{idTooShort, usnInvalidChar},
			passwordErrs: []string{pwdNonASCII},
		},
		{
			name:         "UsnDigitStart,PwdTooLong,PwdNoDigit",
			reqBody:      PostReq{Username: "1bobob", Password: "MyPass!"},
			usernameErrs: []string{usnDigitStart},
			passwordErrs: []string{pwdTooShort, pwdNoDigit},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			fakeIDValidator.errs = c.usernameErrs
			fakePasswordValidator.errs = c.passwordErrs

			res := sut.Validate(c.reqBody)

			assert.AllEqual(t.Error, res.Username, c.usernameErrs)
			assert.AllEqual(t.Error, res.Password, c.passwordErrs)
		})
	}
}

// TestUsernameValidator tests the UsernameValidator to assert that it returns
// the correct error strings based on the username passed to it.
func TestUsernameValidator(t *testing.T) {
	sut := NewUsernameValidator()

	for _, c := range []struct {
		name     string
		username string
		wantErrs []string
	}{
		// 1-error cases
		{name: "Empty", username: "", wantErrs: []string{idEmpty}},
		{name: "TooShort", username: "bob1", wantErrs: []string{idTooShort}},
		{
			name:     "TooLong",
			username: "bobobobobobobobob",
			wantErrs: []string{idTooLong},
		},
		{
			name:     "InvalidCharacter",
			username: "bobob!",
			wantErrs: []string{usnInvalidChar},
		},
		{
			name:     "DigitStart",
			username: "1bobob",
			wantErrs: []string{usnDigitStart},
		},

		// 2-error cases
		{
			name:     "TooShort,InvalidCharacter",
			username: "bob!",
			wantErrs: []string{idTooShort, usnInvalidChar},
		},
		{
			name:     "TooShort,DigitStart",
			username: "1bob",
			wantErrs: []string{idTooShort, usnDigitStart},
		},
		{
			name:     "TooLong,InvalidCharacter",
			username: "bobobobobobobobo!",
			wantErrs: []string{idTooLong, usnInvalidChar},
		},
		{
			name:     "TooLong,DigitStart",
			username: "1bobobobobobobobo",
			wantErrs: []string{idTooLong, usnDigitStart},
		},
		{
			name:     "InvalidCharacter,DigitStart",
			username: "1bob!",
			wantErrs: []string{usnInvalidChar, usnDigitStart},
		},

		// 3-error cases
		{
			name:     "TooShort,InvalidCharacter,DigitStart",
			username: "1bo!",
			wantErrs: []string{idTooShort, usnInvalidChar, usnDigitStart},
		},
		{
			name:     "TooLong,InvalidCharacter,DigitStart",
			username: "1bobobobobobobob!",
			wantErrs: []string{idTooLong, usnInvalidChar, usnDigitStart},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			errs := sut.Validate(c.username)
			assert.AllEqual(t.Error, errs, c.wantErrs)
		})
	}
}

// TestPasswordValidator tests the PasswordValidator to assert that it returns
// the correct error strings based on the password passed to it.
func TestValidatorPassword(t *testing.T) {
	sut := NewPasswordValidator()

	for _, c := range []struct {
		name     string
		password string
		wantErrs []string
	}{
		// 1-error cases
		{
			name:     "Empty",
			password: "", wantErrs: []string{pwdEmpty},
		},
		{
			name:     "TooShort",
			password: "Myp4ss!", wantErrs: []string{pwdTooShort},
		},
		{
			name: "TooLong",
			password: "Myp4sswordwh!chislongandimeanreallylongforsomereasonohikno" +
				"wwhytbh",
			wantErrs: []string{pwdTooLong},
		},
		{
			name:     "NoLower",
			password: "MY4LLUPPERPASSWORD!",
			wantErrs: []string{pwdNoLower},
		},
		{
			name:     "NoUpper",
			password: "my4lllowerpassword!",
			wantErrs: []string{pwdNoUpper},
		},
		{
			name:     "NoDigit",
			password: "myNOdigitPASSWORD!",
			wantErrs: []string{pwdNoDigit},
		},
		{
			name:     "NoSpecial",
			password: "myNOspecialP4SSWORD",
			wantErrs: []string{pwdNoSpecial},
		},
		{
			name:     "HasSpace",
			password: "my SP4CED p4ssword !",
			wantErrs: []string{pwdHasSpace},
		},
		{
			name:     "NonASCII",
			password: "myNØNÅSCÎÎp4ssword!",
			wantErrs: []string{pwdNonASCII},
		},

		// 2-error cases
		{
			name:     "TooShort,NoLower",
			password: "MYP4SS!",
			wantErrs: []string{pwdTooShort, pwdNoLower},
		},
		{
			name:     "TooShort,NoUpper",
			password: "myp4ss!",
			wantErrs: []string{pwdTooShort, pwdNoUpper},
		},
		{
			name:     "TooShort,NoDigit",
			password: "MyPass!",
			wantErrs: []string{pwdTooShort, pwdNoDigit},
		},
		{
			name:     "TooShort,NoSpecial",
			password: "MyP4ssw",
			wantErrs: []string{pwdTooShort, pwdNoSpecial},
		},
		{
			name:     "TooShort,HasSpace",
			password: "My P4s!",
			wantErrs: []string{pwdTooShort, pwdHasSpace},
		},
		{
			name:     "TooShort,NonASCII",
			password: "M¥P4s!2",
			wantErrs: []string{pwdTooShort, pwdNonASCII},
		},
		{
			name: "TooLong,NoLower",
			password: "MYP4SSWORDWH!CHISLONGANDIMEANREALLY" +
				"LONGFORSOMEREASONOHIKNOWWHYTBH",
			wantErrs: []string{pwdTooLong, pwdNoLower},
		},
		{
			name: "TooLong,NoUpper",
			password: "myp4sswordwh!chislongandimeanreally" +
				"longforsomereasonohiknowwhytbh",
			wantErrs: []string{pwdTooLong, pwdNoUpper},
		},
		{
			name: "TooLong,NoDigit",
			password: "Mypasswordwh!chislongandimeanreally" +
				"longforsomereasonohiknowwhytbh",
			wantErrs: []string{pwdTooLong, pwdNoDigit},
		},
		{
			name: "TooLong,NoSpecial",
			password: "Myp4sswordwhichislongandimeanreally" +
				"longforsomereasonohiknowwhytbh",
			wantErrs: []string{pwdTooLong, pwdNoSpecial},
		},
		{
			name: "TooLong,HasSpace",
			password: "Myp4sswo   rdwh!chislongandimeanreally" +
				"longforsomereasonohiknowwhy",
			wantErrs: []string{pwdTooLong, pwdHasSpace},
		},
		{
			name: "TooLong,NonASCII",
			password: "Myp4££wordwh!chislongandimeanreally" +
				"longforsomereasonohiknowwhytbh",
			wantErrs: []string{pwdTooLong, pwdNonASCII},
		},
		{
			name:     "NoLower,NoUpper",
			password: "4444!!!!",
			wantErrs: []string{pwdNoLower, pwdNoUpper},
		},
		{
			name:     "NoLower,NoDigit",
			password: "MYP@SSW!",
			wantErrs: []string{pwdNoLower, pwdNoDigit},
		},
		{
			name:     "NoLower,NoSpecial",
			password: "MYP4SSW1",
			wantErrs: []string{pwdNoLower, pwdNoSpecial},
		},
		{
			name:     "NoLower,HasSpace",
			password: "MYP4SS !",
			wantErrs: []string{pwdNoLower, pwdHasSpace},
		},
		{
			name:     "NoLower,NonASCII",
			password: "MYP4££W!",
			wantErrs: []string{pwdNoLower, pwdNonASCII},
		},
		{
			name:     "NoUpper,NoDigit",
			password: "myp@ssw!",
			wantErrs: []string{pwdNoUpper, pwdNoDigit},
		},
		{
			name:     "NoUpper,NoSpecial",
			password: "myp4ssw1",
			wantErrs: []string{pwdNoUpper, pwdNoSpecial},
		},
		{
			name:     "NoUpper,HasSpace",
			password: "myp4ss !",
			wantErrs: []string{pwdNoUpper, pwdHasSpace},
		},
		{
			name:     "NoUpper,NonASCII",
			password: "myp4££w!",
			wantErrs: []string{pwdNoUpper, pwdNonASCII},
		},
		{
			name:     "NoDigit,NoSpecial",
			password: "MyPasswd",
			wantErrs: []string{pwdNoDigit, pwdNoSpecial},
		},
		{
			name:     "NoDigit,HasSpace",
			password: "MyPass !",
			wantErrs: []string{pwdNoDigit, pwdHasSpace},
		},
		{
			name:     "NoDigit,NonASCII",
			password: "MyPa££w!",
			wantErrs: []string{pwdNoDigit, pwdNonASCII},
		},
		{
			name:     "NoSpecial,HasSpace",
			password: "My  P4ss",
			wantErrs: []string{pwdNoSpecial, pwdHasSpace},
		},
		{
			name:     "NoSpecial,NonASCII",
			password: "MyPa££w1",
			wantErrs: []string{pwdNoSpecial, pwdNonASCII},
		},
		{
			name:     "HasSpace,NonASCII",
			password: "MyP4££ !",
			wantErrs: []string{pwdHasSpace, pwdNonASCII},
		},

		// 3-error cases
		{
			name:     "TooShort,NoLower,NoUpper",
			password: "1421!@$",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper},
		},
		{
			name:     "TooShort,NoLower,NoDigit",
			password: "PASS!@$",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoDigit},
		},
		{
			name:     "TooShort,NoLower,NoSpecial",
			password: "PASS123",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoSpecial},
		},
		{
			name:     "TooShort,NoLower,HasSpace",
			password: "PA$ 123",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdHasSpace},
		},
		{
			name:     "TooShort,NoLower,NonASCII",
			password: "PA$£123",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNonASCII},
		},
		{
			name:     "TooShort,NoUpper,NoDigit",
			password: "pass$$$",
			wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoDigit},
		},
		{
			name:     "TooShort,NoUpper,NoSpecial",
			password: "pass123",
			wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoSpecial},
		},
		{
			name:     "TooShort,NoUpper,HasSpace",
			password: "pa$ 123",
			wantErrs: []string{pwdTooShort, pwdNoUpper, pwdHasSpace},
		},
		{
			name:     "TooShort,NoUpper,NonASCII",
			password: "pa$£123",
			wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNonASCII},
		},
		{
			name:     "TooShort,NoDigit,NoSpecial",
			password: "Passwor",
			wantErrs: []string{pwdTooShort, pwdNoDigit, pwdNoSpecial},
		},
		{
			name:     "TooShort,NoDigit,HasSpace",
			password: "Pa$$ wo",
			wantErrs: []string{pwdTooShort, pwdNoDigit, pwdHasSpace},
		},
		{
			name:     "TooShort,NoDigit,NonASCII",
			password: "Pa$$£wo",
			wantErrs: []string{pwdTooShort, pwdNoDigit, pwdNonASCII},
		},
		{
			name:     "TooShort,NoSpecial,HasSpace",
			password: "Pa55 wo",
			wantErrs: []string{pwdTooShort, pwdNoSpecial, pwdHasSpace},
		},
		{
			name:     "TooShort,NoSpecial,NonASCII",
			password: "Pa55£wo",
			wantErrs: []string{pwdTooShort, pwdNoSpecial, pwdNonASCII},
		},
		{
			name:     "TooShort,HasSpace,NonASCII",
			password: "P4$ £wo",
			wantErrs: []string{pwdTooShort, pwdHasSpace, pwdNonASCII},
		},
		{
			name: "TooLong,NoLower,NoUpper",
			password: "111422222222!3333333333333333333333333333" +
				"333333333333333333333333",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper},
		},
		{
			name: "TooLong,NoLower,NoDigit",
			password: "MYPASSWORDWH!CHISLONGANDIMEANREALLY" +
				"LONGFORSOMEREASONOHIKNOWWHYTBH",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoDigit},
		},
		{
			name: "TooLong,NoLower,NoSpecial",
			password: "MYP4SSWORDWHICHISLONGANDIMEANREALLY" +
				"LONGFORSOMEREASONOHIKNOWWHYTBH",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoSpecial},
		},
		{
			name: "TooLong,NoLower,HasSpace",
			password: "MYP4SS    WH!CHISLONGANDIMEANREALLY" +
				"LONGFORSOMEREASONOHIKNOWWHYTBH",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdHasSpace},
		},
		{
			name: "TooLong,NoLower,NonASCII",
			password: "£YP4SSWORDWH!CHISLONGANDIMEANREALLY" +
				"LONGFORSOMEREASONOHIKNOWWHYTBH",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNonASCII},
		},
		{
			name: "TooLong,NoUpper,NoDigit",
			password: "mypasswordwh!chislongandimeanreally" +
				"longforsomereasonohiknowwhytbh",
			wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoDigit},
		},
		{
			name: "TooLong,NoUpper,NoSpecial",
			password: "myp4sswordwhichislongandimeanreally" +
				"longforsomereasonohiknowwhytbh",
			wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoSpecial},
		},
		{
			name: "TooLong,NoUpper,HasSpace",
			password: "myp4ss    wh!chislongandimeanreally" +
				"longforsomereasonohiknowwhytbh",
			wantErrs: []string{pwdTooLong, pwdNoUpper, pwdHasSpace},
		},
		{
			name: "TooLong,NoUpper,NonASCII",
			password: "£yp4sswordwh!chislongandimeanreally" +
				"longforsomereasonohiknowwhytbh",
			wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNonASCII},
		},
		{
			name: "TooLong,NoDigit,NoSpecial",
			password: "Mypasswordwhichislongandimeanreally" +
				"longforsomereasonohiknowwhytbh",
			wantErrs: []string{pwdTooLong, pwdNoDigit, pwdNoSpecial},
		},
		{
			name: "TooLong,NoDigit,HasSpace",
			password: "Mypass    wh!chislongandimeanreally" +
				"longforsomereasonohiknowwhytbh",
			wantErrs: []string{pwdTooLong, pwdNoDigit, pwdHasSpace},
		},
		{
			name: "TooLong,NoDigit,NonASCII",
			password: "Myp£sswordwh!chislongandimeanreally" +
				"longforsomereasonohiknowwhytbh",
			wantErrs: []string{pwdTooLong, pwdNoDigit, pwdNonASCII},
		},
		{
			name: "TooLong,NoSpecial,HasSpace",
			password: "Myp4ss    whichislongandimeanreally" +
				"longforsomereasonohiknowwhytbh",
			wantErrs: []string{pwdTooLong, pwdNoSpecial, pwdHasSpace},
		},
		{
			name: "TooLong,HasSpace,NonASCII",
			password: "Myp4ssw£   rdwh!chislongandimeanreally" +
				"longforsomereasonohiknowwhy",
			wantErrs: []string{pwdTooLong, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "NoLower,NoUpper,NoDigit",
			password: "!!!!!!!!",
			wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoDigit},
		},
		{
			name:     "NoLower,NoUpper,NoSpecial",
			password: "33333333",
			wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoSpecial},
		},
		{
			name:     "NoLower,NoUpper,HasSpace",
			password: "444  !!!",
			wantErrs: []string{pwdNoLower, pwdNoUpper, pwdHasSpace},
		},
		{
			name:     "NoLower,NoUpper,NonASCII",
			password: "£££444!!",
			wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNonASCII},
		},
		{
			name:     "NoLower,NoDigit,NoSpecial",
			password: "MYPASSWO",
			wantErrs: []string{pwdNoLower, pwdNoDigit, pwdNoSpecial},
		},
		{
			name:     "NoLower,NoDigit,HasSpace",
			password: "MYP@SS !",
			wantErrs: []string{pwdNoLower, pwdNoDigit, pwdHasSpace},
		},
		{
			name:     "NoLower,NoDigit,NonASCII",
			password: "M£P@SSW!",
			wantErrs: []string{pwdNoLower, pwdNoDigit, pwdNonASCII},
		},
		{
			name:     "NoLower,NoSpecial,HasSpace",
			password: "MYP4  W1",
			wantErrs: []string{pwdNoLower, pwdNoSpecial, pwdHasSpace},
		},
		{
			name:     "NoLower,NoSpecial,NonASCII",
			password: "M£P4SSW1",
			wantErrs: []string{pwdNoLower, pwdNoSpecial, pwdNonASCII},
		},
		{
			name:     "NoLower,HasSpace,NonASCII",
			password: "M£P4SS !",
			wantErrs: []string{pwdNoLower, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "NoUpper,NoDigit,NoSpecial",
			password: "mypasswo",
			wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdNoSpecial},
		},
		{
			name:     "NoUpper,NoDigit,HasSpace",
			password: "myp@ss !",
			wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdHasSpace},
		},
		{
			name:     "NoUpper,NoDigit,NonASCII",
			password: "m£p@ssw!",
			wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdNonASCII},
		},
		{
			name:     "NoUpper,NoSpecial,HasSpace",
			password: "myp4ss 1",
			wantErrs: []string{pwdNoUpper, pwdNoSpecial, pwdHasSpace},
		},
		{
			name:     "NoUpper,NoSpecial,NonASCII",
			password: "m£p4ssw1",
			wantErrs: []string{pwdNoUpper, pwdNoSpecial, pwdNonASCII},
		},
		{
			name:     "NoUpper,HasSpace,NonASCII",
			password: "m£p4ss !",
			wantErrs: []string{pwdNoUpper, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "NoDigit,NoSpecial,HasSpace",
			password: "MyPass o",
			wantErrs: []string{pwdNoDigit, pwdNoSpecial, pwdHasSpace},
		},
		{
			name:     "NoDigit,NoSpecial,NonASCII",
			password: "M£Passwd",
			wantErrs: []string{pwdNoDigit, pwdNoSpecial, pwdNonASCII},
		},
		{
			name:     "NoDigit,HasSpace,NonASCII",
			password: "M£Pass !",
			wantErrs: []string{pwdNoDigit, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "NoSpecial,HasSpace,NonASCII",
			password: "M£  P4ss",
			wantErrs: []string{pwdNoSpecial, pwdHasSpace, pwdNonASCII},
		},

		// 4-error cases
		{
			name:     "TooShort,NoLower,NoUpper,NoDigit",
			password: "!@$!@$!",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit},
		},
		{
			name:     "TooShort,NoLower,NoUpper,NoSpecial",
			password: "1421111",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoSpecial},
		},
		{
			name:     "TooShort,NoLower,NoUpper,HasSpace",
			password: "142 !@$",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdHasSpace},
		},
		{
			name:     "TooShort,NoLower,NoUpper,NonASCII",
			password: "14£1!@$",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNonASCII},
		},
		{
			name:     "TooShort,NoLower,NoDigit,NoSpecial",
			password: "PASSSSS",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdNoSpecial},
		},
		{
			name:     "TooShort,NoLower,NoDigit,HasSpace",
			password: "PAS !@$",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdHasSpace},
		},
		{
			name:     "TooShort,NoLower,NoDigit,NonASCII",
			password: "P£SS!@$",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdNonASCII},
		},
		{
			name:     "TooShort,NoLower,NoSpecial,HasSpace",
			password: "PAS 123",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoSpecial, pwdHasSpace},
		},
		{
			name:     "TooShort,NoLower,NoSpecial,NonASCII",
			password: "P£SS123",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoSpecial, pwdNonASCII},
		},
		{
			name:     "TooShort,NoLower,HasSpace,NonASCII",
			password: "P£$ 123",
			wantErrs: []string{pwdTooShort, pwdNoLower, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "TooShort,NoUpper,NoDigit,NoSpecial",
			password: "passsss",
			wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNoSpecial},
		},
		{
			name:     "TooShort,NoUpper,NoDigit,HasSpace",
			password: "pas $$$",
			wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdHasSpace},
		},
		{
			name:     "TooShort,NoUpper,NoDigit,NonASCII",
			password: "p£ss$$$",
			wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNonASCII},
		},
		{
			name:     "TooShort,NoUpper,NoSpecial,HasSpace",
			password: "pas 123",
			wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoSpecial, pwdHasSpace},
		},
		{
			name:     "TooShort,NoUpper,NoSpecial,NonASCII",
			password: "p£ss123",
			wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoSpecial, pwdNonASCII},
		},
		{
			name:     "TooShort,NoUpper,HasSpace,NonASCII",
			password: "p£$ 123",
			wantErrs: []string{pwdTooShort, pwdNoUpper, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "TooShort,NoDigit,NoSpecial,HasSpace",
			password: "Pas wor",
			wantErrs: []string{pwdTooShort, pwdNoDigit, pwdNoSpecial, pwdHasSpace},
		},
		{
			name:     "TooShort,NoDigit,NoSpecial,NonASCII",
			password: "P£sswor",
			wantErrs: []string{pwdTooShort, pwdNoDigit, pwdNoSpecial, pwdNonASCII},
		},
		{
			name:     "TooShort,NoDigit,HasSpace,NonASCII",
			password: "P£$$ wo",
			wantErrs: []string{pwdTooShort, pwdNoDigit, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "TooShort,NoSpecial,HasSpace,NonASCII",
			password: "P£55 wo",
			wantErrs: []string{pwdTooShort, pwdNoSpecial, pwdHasSpace, pwdNonASCII},
		},
		{
			name: "TooLong,NoLower,NoUpper,NoDigit",
			password: "!@$!@$!!@$!@$!!@$!@$!!@$!@$!!" +
				"@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit},
		},
		{
			name: "TooLong,NoLower,NoUpper,NoSpecial",
			password: "142111114211111421111142111114211" +
				"11142111114211114211111421111142",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoSpecial},
		},
		{
			name: "TooLong,NoLower,NoUpper,HasSpace",
			password: "142 !@$142 !@$142 !@$142 !@$142 " +
				"!@$142 !@$142 !@$142 !@$142 !@$14",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdHasSpace},
		},
		{
			name: "TooLong,NoLower,NoUpper,NonASCII",
			password: "14£1!@$14£1!@$14£1!@$14£1!@$14" +
				"£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNonASCII},
		},
		{
			name: "TooLong,NoLower,NoDigit,NoSpecial",
			password: "PASSSSSPASSSSSPASSSSSPASSSSSPASSS" +
				"SSPASSSSSPASSSSSPASSSSSPASSSSSPA",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdNoSpecial},
		},
		{
			name: "TooLong,NoLower,NoDigit,HasSpace",
			password: "PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS " +
				"!@$PAS !@$PAS !@$PAS !@$PA",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdHasSpace},
		},
		{
			name: "TooLong,NoLower,NoDigit,NonASCII",
			password: "P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@" +
				"$P£SS!@$P£SS!@$P£SS!@$P£",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdNonASCII},
		},
		{
			name: "TooLong,NoLower,NoSpecial,HasSpace",
			password: "PAS 123PAS 123PAS 123PAS 123PAS 123PAS " +
				"123PAS 123PAS 123PAS 123PA",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoSpecial, pwdHasSpace},
		},
		{
			name: "TooLong,NoLower,NoSpecial,NonASCII",
			password: "P£SS123P£SS123P£SS123P£SS123P£SS123P£SS1" +
				"23P£SS123P£SS123P£SS123P£",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoSpecial, pwdNonASCII},
		},
		{
			name: "TooLong,NoLower,HasSpace,NonASCII",
			password: "P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P" +
				"£$ 123P£$ 123P£$ 123P£",
			wantErrs: []string{pwdTooLong, pwdNoLower, pwdHasSpace, pwdNonASCII},
		},
		{
			name: "TooLong,NoUpper,NoDigit,NoSpecial",
			password: "passssspassssspassssspassssspassssspas" +
				"sssspassssspassssspassssspa",
			wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNoSpecial},
		},
		{
			name: "TooLong,NoUpper,NoDigit,HasSpace",
			password: "pas $$$pas $$$pas $$$pas $$$pas $$$pas " +
				"$$$pas $$$pas $$$pas $$$pa",
			wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdHasSpace},
		},
		{
			name: "TooLong,NoUpper,NoDigit,NonASCII",
			password: "p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$" +
				"$p£ss$$$p£ss$$$p£ss$$$p£",
			wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNonASCII},
		},
		{
			name: "TooLong,NoUpper,NoSpecial,HasSpace",
			password: "pas 123pas 123pas 123pas 123pas 123pas 123pas " +
				"123pas 123pas 123pa",
			wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoSpecial, pwdHasSpace},
		},
		{
			name: "TooLong,NoUpper,NoSpecial,NonASCII",
			password: "p£ss123p£ss123p£ss123p£ss123p£ss123p£ss123p£" +
				"ss123p£ss123p£p£ss123",
			wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoSpecial, pwdNonASCII},
		},
		{
			name: "TooLong,NoUpper,HasSpace,NonASCII",
			password: "p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p" +
				"£$ 123p£$ 123p£$ 123p£",
			wantErrs: []string{pwdTooLong, pwdNoUpper, pwdHasSpace, pwdNonASCII},
		},
		{
			name: "TooLong,NoDigit,NoSpecial,HasSpace",
			password: "Pas worPas worPas worPas worPas worPas " +
				"worPas worPas worPas worPa",
			wantErrs: []string{pwdTooLong, pwdNoDigit, pwdNoSpecial, pwdHasSpace},
		},
		{
			name: "TooLong,NoDigit,NoSpecial,NonASCII",
			password: "P£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£" +
				"ssworP£ssworP£ssworP£",
			wantErrs: []string{pwdTooLong, pwdNoDigit, pwdNoSpecial, pwdNonASCII},
		},
		{
			name: "TooLong,NoDigit,HasSpace,NonASCII",
			password: "P£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ " +
				"woP£$$ woP£$$ woP£$$ woP£",
			wantErrs: []string{pwdTooLong, pwdNoDigit, pwdHasSpace, pwdNonASCII},
		},
		{
			name: "TooLong,NoSpecial,HasSpace,NonASCII",
			password: "P£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP" +
				"£55 woP£55 woP£55 woP£",
			wantErrs: []string{pwdTooLong, pwdNoSpecial, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "NoLower,NoUpper,NoDigit,HasSpace",
			password: "!!!  !!!",

			wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace},
		},
		{
			name:     "NoLower,NoUpper,NoDigit,NonASCII",
			password: "!!!££!!!",
			wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNonASCII},
		},
		{
			name:     "NoLower,NoUpper,NoSpecial,HasSpace",
			password: "333  333",
			wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace},
		},
		{
			name:     "NoLower,NoUpper,NoSpecial,NonASCII",
			password: "333££333",
			wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdNonASCII},
		},
		{
			name:     "NoLower,NoUpper,HasSpace,NonASCII",
			password: "4£4  !!!",
			wantErrs: []string{pwdNoLower, pwdNoUpper, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "NoLower,NoDigit,NoSpecial,HasSpace",
			password: "MYP  SWO",
			wantErrs: []string{pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace},
		},
		{
			name:     "NoLower,NoDigit,NoSpecial,NonASCII",
			password: "MYP££SWO",
			wantErrs: []string{pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdNonASCII},
		},
		{
			name:     "NoLower,NoDigit,HasSpace,NonASCII",
			password: "M£P@SS !",
			wantErrs: []string{pwdNoLower, pwdNoDigit, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "NoLower,NoSpecial,HasSpace,NonASCII",
			password: "M£P4  W1",
			wantErrs: []string{pwdNoLower, pwdNoSpecial, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "NoUpper,NoDigit,NoSpecial,HasSpace",
			password: "myp  swo",
			wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace},
		},
		{
			name:     "NoUpper,NoDigit,NoSpecial,NonASCII",
			password: "myp££swo",
			wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII},
		},
		{
			name:     "NoUpper,NoDigit,HasSpace,NonASCII",
			password: "m£p@ss !",
			wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "NoUpper,NoSpecial,HasSpace,NonASCII",
			password: "m£p4  w1",
			wantErrs: []string{pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII},
		},
		{
			name:     "NoDigit,NoSpecial,HasSpace,NonASCII",
			password: "MyP£ss o",
			wantErrs: []string{pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII},
		},

		// 5-error cases
		{
			name:     "TooShort,NoLower,NoUpper,NoDigit,HasSpace",
			password: "!@   $!",
			wantErrs: []string{
				pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace,
			},
		},
		{
			name:     "TooShort,NoLower,NoUpper,NoDigit,NonASCII",
			password: "!@£££$!",
			wantErrs: []string{
				pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoLower,NoUpper,NoSpecial,HasSpace",
			password: "14   11",
			wantErrs: []string{
				pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace,
			},
		},
		{
			name:     "TooShort,NoLower,NoUpper,NoSpecial,NonASCII",
			password: "14£££11",
			wantErrs: []string{
				pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoLower,NoUpper,HasSpace,NonASCII",
			password: "1£2 !@$",
			wantErrs: []string{
				pwdTooShort, pwdNoLower, pwdNoUpper, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoLower,NoDigit,NoSpecial,HasSpace",
			password: "PAS SSS",
			wantErrs: []string{
				pwdTooShort, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace,
			},
		},
		{
			name:     "TooShort,NoLower,NoDigit,NoSpecial,NonASCII",
			password: "PAS£SSS",
			wantErrs: []string{
				pwdTooShort, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoLower,NoDigit,HasSpace,NonASCII",
			password: "P£S !@$",
			wantErrs: []string{
				pwdTooShort, pwdNoLower, pwdNoDigit, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoLower,NoSpecial,HasSpace,NonASCII",
			password: "P£S 123",
			wantErrs: []string{
				pwdTooShort, pwdNoLower, pwdNoSpecial, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoUpper,NoDigit,NoSpecial,HasSpace",
			password: "pas sss",
			wantErrs: []string{
				pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace,
			},
		},
		{
			name:     "TooShort,NoUpper,NoDigit,NoSpecial,NonASCII",
			password: "pas£sss",
			wantErrs: []string{
				pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoUpper,NoDigit,HasSpace,NonASCII",
			password: "p£s $$$",
			wantErrs: []string{
				pwdTooShort, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoUpper,NoSpecial,HasSpace,NonASCII",
			password: "p£s 123",
			wantErrs: []string{
				pwdTooShort, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoDigit,NoSpecial,HasSpace,NonASCII",
			password: "P£s wor",
			wantErrs: []string{
				pwdTooShort, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoLower,NoUpper,NoDigit,HasSpace",
			password: "!@   $!!@   $!!@   $!!@   $!!@   " +
				"$!!@   $!!@   $!!@   $!!@   $!!@",
			wantErrs: []string{
				pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace,
			},
		},
		{
			name: "TooLong,NoLower,NoUpper,NoDigit,NonASCII",
			password: "!@£££$!!@£££$!!@£££$!!@£££$!!" +
				"@£££$!!@£££$!!@£££$!!@£££$!!@£££$!!@",
			wantErrs: []string{
				pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoLower,NoUpper,NoSpecial,HasSpace",
			password: "14   1114   1114   1114   1114   1114   " +
				"1114   1114   1114   1114",
			wantErrs: []string{
				pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace,
			},
		},
		{
			name: "TooLong,NoLower,NoUpper,NoSpecial,NonASCII",
			password: "14£££1114£££1114£££1114£££1114££" +
				"£1114£££1114£££1114£££1114£££1114",
			wantErrs: []string{
				pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoLower,NoUpper,HasSpace,NonASCII",
			password: "1£2 !@$1£2 !@$1£2 !@$1£2 " +
				"!@$1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£",
			wantErrs: []string{
				pwdTooLong, pwdNoLower, pwdNoUpper, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoLower,NoDigit,NoSpecial,HasSpace",
			password: "PAS SSSPAS SSSPAS SSSPAS SSSPAS " +
				"SSSPAS SSSPAS SSSPAS SSSPAS SSSPA",
			wantErrs: []string{
				pwdTooLong, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace,
			},
		},
		{
			name: "TooLong,NoLower,NoDigit,NoSpecial,NonASCII",
			password: "PAS£SSSPAS£SSSPAS£SSSPAS£SSSPAS£SSSP" +
				"AS£SSSPAS£SSSPAS£SSSPAS£SSSPA",
			wantErrs: []string{
				pwdTooLong, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoLower,NoDigit,HasSpace,NonASCII",
			password: "P£S !@$P£S !@$P£S !@$P£S !@$P£S !@$P£S " +
				"!@$P£S !@$P£S !@$P£S !@$P£",
			wantErrs: []string{
				pwdTooLong, pwdNoLower, pwdNoDigit, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoLower,NoSpecial,HasSpace,NonASCII",
			password: "P£S 123P£S 123P£S 123P£S 123P£S 123P£S 123P" +
				"£S 123P£S 123P£S 123P£",
			wantErrs: []string{
				pwdTooLong, pwdNoLower, pwdNoSpecial, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoUpper,NoDigit,NoSpecial,HasSpace",
			password: "pas ssspas ssspas ssspas ssspas ssspas " +
				"ssspas ssspas ssspas ssspa",
			wantErrs: []string{
				pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace,
			},
		},
		{
			name: "TooLong,NoUpper,NoDigit,NoSpecial,NonASCII",
			password: "pas£ssspas£ssspas£ssspas£ssspas£ssspas£" +
				"ssspas£ssspas£ssspas£ssspa",
			wantErrs: []string{
				pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoUpper,NoDigit,HasSpace,NonASCII",
			password: "p£s $$$p£s $$$p£s $$$p£s $$$p£s $$$p£s " +
				"$$$p£s $$$p£s $$$p£s $$$p£",
			wantErrs: []string{
				pwdTooLong, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoUpper,NoSpecial,HasSpace,NonASCII",
			password: "p£s 123p£s 123p£s 123p£s 123p£s 123p£s " +
				"123p£s 123p£s 123p£s 123p£",
			wantErrs: []string{
				pwdTooLong, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoDigit,NoSpecial,HasSpace,NonASCII",
			password: "P£s worP£s worP£s worP£s worP£s worP£s worP£s " +
				"worP£s worP£s worP£",
			wantErrs: []string{
				pwdTooLong, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name:     "NoLower,NoUpper,NoDigit,NoSpecial,HasSpace",
			password: "        ",
			wantErrs: []string{
				pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace,
			},
		},
		{
			name:     "NoLower,NoUpper,NoDigit,NoSpecial,NonASCII",
			password: "££££££££",
			wantErrs: []string{
				pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII,
			},
		},
		{
			name:     "NoLower,NoUpper,NoDigit,HasSpace,NonASCII",
			password: "!£!  !!!",
			wantErrs: []string{
				pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name:     "NoLower,NoUpper,NoSpecial,HasSpace,NonASCII",
			password: "3£3  333",
			wantErrs: []string{
				pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name:     "NoLower,NoDigit,NoSpecial,HasSpace,NonASCII",
			password: "M£P  SWO",
			wantErrs: []string{
				pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII,
			},
		},
		{
			name:     "NoUpper,NoDigit,NoSpecial,HasSpace,NonASCII",
			password: "m£p  swo",
			wantErrs: []string{
				pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII,
			},
		},

		// 6-error cases
		{
			name:     "TooShort,NoLower,NoUpper,NoDigit,NoSpecial,HasSpace",
			password: "       ",
			wantErrs: []string{
				pwdTooShort,
				pwdNoLower,
				pwdNoUpper,
				pwdNoDigit,
				pwdNoSpecial,
				pwdHasSpace,
			},
		},
		{
			name:     "TooShort,NoLower,NoUpper,NoDigit,NoSpecial,NonASCII",
			password: "£££££££",
			wantErrs: []string{
				pwdTooShort,
				pwdNoLower,
				pwdNoUpper,
				pwdNoDigit,
				pwdNoSpecial,
				pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoLower,NoUpper,NoDigit,HasSpace,NonASCII",
			password: "!£   $!",
			wantErrs: []string{
				pwdTooShort,
				pwdNoLower,
				pwdNoUpper,
				pwdNoDigit,
				pwdHasSpace,
				pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoLower,NoUpper,NoSpecial,HasSpace,NonASCII",
			password: "1£   11",
			wantErrs: []string{
				pwdTooShort,
				pwdNoLower,
				pwdNoUpper,
				pwdNoSpecial,
				pwdHasSpace,
				pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoLower,NoDigit,NoSpecial,HasSpace,NonASCII",
			password: "P£S SSS",
			wantErrs: []string{
				pwdTooShort,
				pwdNoLower,
				pwdNoDigit,
				pwdNoSpecial,
				pwdHasSpace,
				pwdNonASCII,
			},
		},
		{
			name:     "TooShort,NoUpper,NoDigit,NoSpecial,HasSpace,NonASCII",
			password: "p£s sss",
			wantErrs: []string{
				pwdTooShort,
				pwdNoUpper,
				pwdNoDigit,
				pwdNoSpecial,
				pwdHasSpace,
				pwdNonASCII,
			},
		},
		{
			name:     "TooLong,NoLower,NoUpper,NoDigit,NoSpecial,HasSpace",
			password: "                                                                 ",
			wantErrs: []string{
				pwdTooLong,
				pwdNoLower,
				pwdNoUpper,
				pwdNoDigit,
				pwdNoSpecial,
				pwdHasSpace,
			},
		},
		{
			name: "TooLong,NoLower,NoUpper,NoDigit,NoSpecial,NonASCII",
			password: "££££££££££££££££££££££££££££££££££££££££££" +
				"£££££££££££££££££££££££",
			wantErrs: []string{
				pwdTooLong,
				pwdNoLower,
				pwdNoUpper,
				pwdNoDigit,
				pwdNoSpecial,
				pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoLower,NoUpper,NoDigit,HasSpace,NonASCII",
			password: "!£   $!!£   $!!£   $!!£   $!!£   $!!£   " +
				"$!!£   $!!£   $!!£   $!!£",
			wantErrs: []string{
				pwdTooLong,
				pwdNoLower,
				pwdNoUpper,
				pwdNoDigit,
				pwdHasSpace,
				pwdNonASCII,
			},
		},
		{
			name:     "TooLong,NoLower,NoUpper,NoSpecial,HasSpace,NonASCII",
			password: "1£   111£   111£   111£   111£   111£   111£   111£   111£   111£",
			wantErrs: []string{
				pwdTooLong,
				pwdNoLower,
				pwdNoUpper,
				pwdNoSpecial,
				pwdHasSpace,
				pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoLower,NoDigit,NoSpecial,HasSpace,NonASCII",
			password: "P£S SSSP£S SSSP£S SSSP£S SSSP£S SSSP£S " +
				"SSSP£S SSSP£S SSSP£S SSSP£",
			wantErrs: []string{
				pwdTooLong,
				pwdNoLower,
				pwdNoDigit,
				pwdNoSpecial,
				pwdHasSpace,
				pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoUpper,NoDigit,NoSpecial,HasSpace,NonASCII",
			password: "p£s sssp£s sssp£s sssp£s sssp£s sssp£s sssp£s " +
				"sssp£s sssp£s sssp£",
			wantErrs: []string{
				pwdTooLong,
				pwdNoUpper,
				pwdNoDigit,
				pwdNoSpecial,
				pwdHasSpace,
				pwdNonASCII,
			},
		},
		{
			name:     "NoLower,NoUpper,NoDigit,NoSpecial,HasSpace,NonASCII",
			password: "   ££   ",
			wantErrs: []string{
				pwdNoLower,
				pwdNoUpper,
				pwdNoDigit,
				pwdNoSpecial,
				pwdHasSpace,
				pwdNonASCII,
			},
		},

		// 7-error cases
		{
			name:     "TooShort,NoLower,NoUpper,NoDigit,NoSpecial,HasSpace,NonASCII",
			password: "   £   ",
			wantErrs: []string{
				pwdTooShort,
				pwdNoLower,
				pwdNoUpper,
				pwdNoDigit,
				pwdNoSpecial,
				pwdHasSpace,
				pwdNonASCII,
			},
		},
		{
			name: "TooLong,NoLower,NoUpper,NoDigit,NoSpecial,HasSpace,NonASCII",
			password: "   £      £      £      £      £      £      " +
				"£      £      £     ",
			wantErrs: []string{
				pwdTooLong,
				pwdNoLower,
				pwdNoUpper,
				pwdNoDigit,
				pwdNoSpecial,
				pwdHasSpace,
				pwdNonASCII,
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			gotErrs := sut.Validate(c.password)

			assert.AllEqual(t.Error, c.wantErrs, gotErrs)
		})
	}
}
