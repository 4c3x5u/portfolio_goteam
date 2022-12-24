package register

import (
	"testing"

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

// This suite of tests is only for ensuring valdator returns the correct Errs
// object based on the errors that field validators return. Each possible
// outcome from each field validator is tested in the following test suites.
func TestValidator(t *testing.T) {
	fakeValidatorUsername := &fakeValidatorStr{}
	fakeValidatorPassword := &fakeValidatorStr{}

	sut := NewValidatorReq(fakeValidatorUsername, fakeValidatorPassword)

	for _, c := range []struct {
		name         string
		reqBody      *Req
		errsUsername []string
		errsPassword []string
	}{
		{
			name:         "UsnEmpty_PwdEmpty",
			reqBody:      &Req{Username: "", Password: ""},
			errsUsername: []string{usnEmpty},
			errsPassword: []string{pwdEmpty},
		},
		{
			name:         "UsnTooShort_UsnInvalidChar_PwdEmpty",
			reqBody:      &Req{Username: "bob!", Password: "myNØNÅSCÎÎp4ssword!"},
			errsUsername: []string{usnTooShort, usnInvalidChar},
			errsPassword: []string{pwdNonASCII},
		},
		{
			name:         "UsnDigitStart_PwdTooLong_PwdNoDigit",
			reqBody:      &Req{Username: "1bobob", Password: "MyPass!"},
			errsUsername: []string{usnDigitStart},
			errsPassword: []string{pwdTooShort, pwdNoDigit},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			// arrange field validator return values
			fakeValidatorUsername.outErrs = c.errsUsername
			fakeValidatorPassword.outErrs = c.errsPassword

			// act
			res := sut.Validate(c.reqBody)

			// assert that ValidatorReq passed in the arguments to field
			// validators correctly
			assert.Equal(t, c.reqBody.Username, fakeValidatorUsername.inArg)
			assert.Equal(t, c.reqBody.Password, fakeValidatorPassword.inArg)

			// assert that ValidatorReq returned the correct errors based on
			// field validator return values
			assert.EqualArr(t, c.errsUsername, res.Username)
			assert.EqualArr(t, c.errsPassword, res.Password)
		})
	}
}

func TestValidatorUsername(t *testing.T) {
	sut := NewValidatorUsername()

	for _, c := range []struct {
		name     string
		username string
		wantErrs []string
	}{
		// 1-error cases
		{name: "Empty", username: "", wantErrs: []string{usnEmpty}},
		{name: "TooShort", username: "bob1", wantErrs: []string{usnTooShort}},
		{name: "TooLong", username: "bobobobobobobobob", wantErrs: []string{usnTooLong}},
		{name: "InvalidCharacter", username: "bobob!", wantErrs: []string{usnInvalidChar}},
		{name: "DigitStart", username: "1bobob", wantErrs: []string{usnDigitStart}},

		// 2-error cases
		{name: "TooShort_InvalidCharacter", username: "bob!", wantErrs: []string{usnTooShort, usnInvalidChar}},
		{name: "TooShort_DigitStart", username: "1bob", wantErrs: []string{usnTooShort, usnDigitStart}},
		{name: "TooLong_InvalidCharacter", username: "bobobobobobobobo!", wantErrs: []string{usnTooLong, usnInvalidChar}},
		{name: "TooLong_DigitStart", username: "1bobobobobobobobo", wantErrs: []string{usnTooLong, usnDigitStart}},
		{name: "InvalidCharacter_DigitStart", username: "1bob!", wantErrs: []string{usnInvalidChar, usnDigitStart}},

		// 3-error cases
		{name: "TooShort_InvalidCharacter_DigitStart", username: "1bo!", wantErrs: []string{usnTooShort, usnInvalidChar, usnDigitStart}},
		{name: "TooLong_InvalidCharacter_DigitStart", username: "1bobobobobobobob!", wantErrs: []string{usnTooLong, usnInvalidChar, usnDigitStart}},
	} {
		t.Run(c.name, func(t *testing.T) {
			gotErrs := sut.Validate(c.username)
			assert.EqualArr(t, c.wantErrs, gotErrs)
		})
	}
}

func TestValidatorPassword(t *testing.T) {
	sut := NewValidatorPassword()

	for _, c := range []struct {
		name     string
		password string
		wantErrs []string
	}{
		// 1-error cases
		{name: "Empty", password: "", wantErrs: []string{pwdEmpty}},
		{name: "TooShort", password: "Myp4ss!", wantErrs: []string{pwdTooShort}},
		{name: "TooLong", password: "Myp4sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong}},
		{name: "NoLower", password: "MY4LLUPPERPASSWORD!", wantErrs: []string{pwdNoLower}},
		{name: "NoUpper", password: "my4lllowerpassword!", wantErrs: []string{pwdNoUpper}},
		{name: "NoDigit", password: "myNOdigitPASSWORD!", wantErrs: []string{pwdNoDigit}},
		{name: "NoSpecial", password: "myNOspecialP4SSWORD", wantErrs: []string{pwdNoSpecial}},
		{name: "HasSpace", password: "my SP4CED p4ssword !", wantErrs: []string{pwdHasSpace}},
		{name: "NonASCII", password: "myNØNÅSCÎÎp4ssword!", wantErrs: []string{pwdNonASCII}},

		// 2-error cases
		{name: "TooShort_NoLower", password: "MYP4SS!", wantErrs: []string{pwdTooShort, pwdNoLower}},
		{name: "TooShort_NoUpper", password: "myp4ss!", wantErrs: []string{pwdTooShort, pwdNoUpper}},
		{name: "TooShort_NoDigit", password: "MyPass!", wantErrs: []string{pwdTooShort, pwdNoDigit}},
		{name: "TooShort_NoSpecial", password: "MyP4ssw", wantErrs: []string{pwdTooShort, pwdNoSpecial}},
		{name: "TooShort_HasSpace", password: "My P4s!", wantErrs: []string{pwdTooShort, pwdHasSpace}},
		{name: "TooShort_NonASCII", password: "M¥P4s!2", wantErrs: []string{pwdTooShort, pwdNonASCII}},
		{name: "TooLong_NoLower", password: "MYP4SSWORDWH!CHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH", wantErrs: []string{pwdTooLong, pwdNoLower}},
		{name: "TooLong_NoUpper", password: "myp4sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong, pwdNoUpper}},
		{name: "TooLong_NoDigit", password: "Mypasswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong, pwdNoDigit}},
		{name: "TooLong_NoSpecial", password: "Myp4sswordwhichislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong, pwdNoSpecial}},
		{name: "TooLong_HasSpace", password: "Myp4sswo   rdwh!chislongandimeanreallylongforsomereasonohiknowwhy", wantErrs: []string{pwdTooLong, pwdHasSpace}},
		{name: "TooLong_NonASCII", password: "Myp4££wordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong, pwdNonASCII}},
		{name: "NoLower_NoUpper", password: "4444!!!!", wantErrs: []string{pwdNoLower, pwdNoUpper}},
		{name: "NoLower_NoDigit", password: "MYP@SSW!", wantErrs: []string{pwdNoLower, pwdNoDigit}},
		{name: "NoLower_NoSpecial", password: "MYP4SSW1", wantErrs: []string{pwdNoLower, pwdNoSpecial}},
		{name: "NoLower_HasSpace", password: "MYP4SS !", wantErrs: []string{pwdNoLower, pwdHasSpace}},
		{name: "NoLower_NonASCII", password: "MYP4££W!", wantErrs: []string{pwdNoLower, pwdNonASCII}},
		{name: "NoUpper_NoDigit", password: "myp@ssw!", wantErrs: []string{pwdNoUpper, pwdNoDigit}},
		{name: "NoUpper_NoSpecial", password: "myp4ssw1", wantErrs: []string{pwdNoUpper, pwdNoSpecial}},
		{name: "NoUpper_HasSpace", password: "myp4ss !", wantErrs: []string{pwdNoUpper, pwdHasSpace}},
		{name: "NoUpper_NonASCII", password: "myp4££w!", wantErrs: []string{pwdNoUpper, pwdNonASCII}},
		{name: "NoDigit_NoSpecial", password: "MyPasswd", wantErrs: []string{pwdNoDigit, pwdNoSpecial}},
		{name: "NoDigit_HasSpace", password: "MyPass !", wantErrs: []string{pwdNoDigit, pwdHasSpace}},
		{name: "NoDigit_NonASCII", password: "MyPa££w!", wantErrs: []string{pwdNoDigit, pwdNonASCII}},
		{name: "NoSpecial_HasSpace", password: "My  P4ss", wantErrs: []string{pwdNoSpecial, pwdHasSpace}},
		{name: "NoSpecial_NonASCII", password: "MyPa££w1", wantErrs: []string{pwdNoSpecial, pwdNonASCII}},
		{name: "HasSpace_NonASCII", password: "MyP4££ !", wantErrs: []string{pwdHasSpace, pwdNonASCII}},

		// 3-error cases
		{name: "TooShort_NoLower_NoUpper", password: "1421!@$", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper}},
		{name: "TooShort_NoLower_NoDigit", password: "PASS!@$", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoDigit}},
		{name: "TooShort_NoLower_NoSpecial", password: "PASS123", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoSpecial}},
		{name: "TooShort_NoLower_HasSpace", password: "PA$ 123", wantErrs: []string{pwdTooShort, pwdNoLower, pwdHasSpace}},
		{name: "TooShort_NoLower_NonASCII", password: "PA$£123", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNonASCII}},
		{name: "TooShort_NoUpper_NoDigit", password: "pass$$$", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoDigit}},
		{name: "TooShort_NoUpper_NoSpecial", password: "pass123", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoSpecial}},
		{name: "TooShort_NoUpper_HasSpace", password: "pa$ 123", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdHasSpace}},
		{name: "TooShort_NoUpper_NonASCII", password: "pa$£123", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNonASCII}},
		{name: "TooShort_NoDigit_NoSpecial", password: "Passwor", wantErrs: []string{pwdTooShort, pwdNoDigit, pwdNoSpecial}},
		{name: "TooShort_NoDigit_HasSpace", password: "Pa$$ wo", wantErrs: []string{pwdTooShort, pwdNoDigit, pwdHasSpace}},
		{name: "TooShort_NoDigit_NonASCII", password: "Pa$$£wo", wantErrs: []string{pwdTooShort, pwdNoDigit, pwdNonASCII}},
		{name: "TooShort_NoSpecial_HasSpace", password: "Pa55 wo", wantErrs: []string{pwdTooShort, pwdNoSpecial, pwdHasSpace}},
		{name: "TooShort_NoSpecial_NonASCII", password: "Pa55£wo", wantErrs: []string{pwdTooShort, pwdNoSpecial, pwdNonASCII}},
		{name: "TooShort_HasSpace_NonASCII", password: "P4$ £wo", wantErrs: []string{pwdTooShort, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoLower_NoUpper", password: "111422222222!3333333333333333333333333333333333333333333333333333", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper}},
		{name: "TooLong_NoLower_NoDigit", password: "MYPASSWORDWH!CHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoDigit}},
		{name: "TooLong_NoLower_NoSpecial", password: "MYP4SSWORDWHICHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoSpecial}},
		{name: "TooLong_NoLower_HasSpace", password: "MYP4SS    WH!CHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH", wantErrs: []string{pwdTooLong, pwdNoLower, pwdHasSpace}},
		{name: "TooLong_NoLower_NonASCII", password: "£YP4SSWORDWH!CHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNonASCII}},
		{name: "TooLong_NoUpper_NoDigit", password: "mypasswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoDigit}},
		{name: "TooLong_NoUpper_NoSpecial", password: "myp4sswordwhichislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoSpecial}},
		{name: "TooLong_NoUpper_HasSpace", password: "myp4ss    wh!chislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdHasSpace}},
		{name: "TooLong_NoUpper_NonASCII", password: "£yp4sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNonASCII}},
		{name: "TooLong_NoDigit_NoSpecial", password: "Mypasswordwhichislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong, pwdNoDigit, pwdNoSpecial}},
		{name: "TooLong_NoDigit_HasSpace", password: "Mypass    wh!chislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong, pwdNoDigit, pwdHasSpace}},
		{name: "TooLong_NoDigit_NonASCII", password: "Myp£sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong, pwdNoDigit, pwdNonASCII}},
		{name: "TooLong_NoSpecial_HasSpace", password: "Myp4ss    whichislongandimeanreallylongforsomereasonohiknowwhytbh", wantErrs: []string{pwdTooLong, pwdNoSpecial, pwdHasSpace}},
		{name: "TooLong_HasSpace_NonASCII", password: "Myp4ssw£   rdwh!chislongandimeanreallylongforsomereasonohiknowwhy", wantErrs: []string{pwdTooLong, pwdHasSpace, pwdNonASCII}},
		{name: "NoLower_NoUpper_NoDigit", password: "!!!!!!!!", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoDigit}},
		{name: "NoLower_NoUpper_NoSpecial", password: "33333333", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoSpecial}},
		{name: "NoLower_NoUpper_HasSpace", password: "444  !!!", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdHasSpace}},
		{name: "NoLower_NoUpper_NonASCII", password: "£££444!!", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNonASCII}},
		{name: "NoLower_NoDigit_NoSpecial", password: "MYPASSWO", wantErrs: []string{pwdNoLower, pwdNoDigit, pwdNoSpecial}},
		{name: "NoLower_NoDigit_HasSpace", password: "MYP@SS !", wantErrs: []string{pwdNoLower, pwdNoDigit, pwdHasSpace}},
		{name: "NoLower_NoDigit_NonASCII", password: "M£P@SSW!", wantErrs: []string{pwdNoLower, pwdNoDigit, pwdNonASCII}},
		{name: "NoLower_NoSpecial_HasSpace", password: "MYP4  W1", wantErrs: []string{pwdNoLower, pwdNoSpecial, pwdHasSpace}},
		{name: "NoLower_NoSpecial_NonASCII", password: "M£P4SSW1", wantErrs: []string{pwdNoLower, pwdNoSpecial, pwdNonASCII}},
		{name: "NoLower_HasSpace_NonASCII", password: "M£P4SS !", wantErrs: []string{pwdNoLower, pwdHasSpace, pwdNonASCII}},
		{name: "NoUpper_NoDigit_NoSpecial", password: "mypasswo", wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdNoSpecial}},
		{name: "NoUpper_NoDigit_HasSpace", password: "myp@ss !", wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdHasSpace}},
		{name: "NoUpper_NoDigit_NonASCII", password: "m£p@ssw!", wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdNonASCII}},
		{name: "NoUpper_NoSpecial_HasSpace", password: "myp4ss 1", wantErrs: []string{pwdNoUpper, pwdNoSpecial, pwdHasSpace}},
		{name: "NoUpper_NoSpecial_NonASCII", password: "m£p4ssw1", wantErrs: []string{pwdNoUpper, pwdNoSpecial, pwdNonASCII}},
		{name: "NoUpper_HasSpace_NonASCII", password: "m£p4ss !", wantErrs: []string{pwdNoUpper, pwdHasSpace, pwdNonASCII}},
		{name: "NoDigit_NoSpecial_HasSpace", password: "MyPass o", wantErrs: []string{pwdNoDigit, pwdNoSpecial, pwdHasSpace}},
		{name: "NoDigit_NoSpecial_NonASCII", password: "M£Passwd", wantErrs: []string{pwdNoDigit, pwdNoSpecial, pwdNonASCII}},
		{name: "NoDigit_HasSpace_NonASCII", password: "M£Pass !", wantErrs: []string{pwdNoDigit, pwdHasSpace, pwdNonASCII}},
		{name: "NoSpecial_HasSpace_NonASCII", password: "M£  P4ss", wantErrs: []string{pwdNoSpecial, pwdHasSpace, pwdNonASCII}},

		// 4-error cases
		{name: "TooShort_NoLower_NoUpper_NoDigit", password: "!@$!@$!", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit}},
		{name: "TooShort_NoLower_NoUpper_NoSpecial", password: "1421111", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoSpecial}},
		{name: "TooShort_NoLower_NoUpper_HasSpace", password: "142 !@$", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdHasSpace}},
		{name: "TooShort_NoLower_NoUpper_NonASCII", password: "14£1!@$", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNonASCII}},
		{name: "TooShort_NoLower_NoDigit_NoSpecial", password: "PASSSSS", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdNoSpecial}},
		{name: "TooShort_NoLower_NoDigit_HasSpace", password: "PAS !@$", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdHasSpace}},
		{name: "TooShort_NoLower_NoDigit_NonASCII", password: "P£SS!@$", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdNonASCII}},
		{name: "TooShort_NoLower_NoSpecial_HasSpace", password: "PAS 123", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoSpecial, pwdHasSpace}},
		{name: "TooShort_NoLower_NoSpecial_NonASCII", password: "P£SS123", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoSpecial, pwdNonASCII}},
		{name: "TooShort_NoLower_HasSpace_NonASCII", password: "P£$ 123", wantErrs: []string{pwdTooShort, pwdNoLower, pwdHasSpace, pwdNonASCII}},
		{name: "TooShort_NoUpper_NoDigit_NoSpecial", password: "passsss", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNoSpecial}},
		{name: "TooShort_NoUpper_NoDigit_HasSpace", password: "pas $$$", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdHasSpace}},
		{name: "TooShort_NoUpper_NoDigit_NonASCII", password: "p£ss$$$", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNonASCII}},
		{name: "TooShort_NoUpper_NoSpecial_HasSpace", password: "pas 123", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoSpecial, pwdHasSpace}},
		{name: "TooShort_NoUpper_NoSpecial_NonASCII", password: "p£ss123", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoSpecial, pwdNonASCII}},
		{name: "TooShort_NoUpper_HasSpace_NonASCII", password: "p£$ 123", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdHasSpace, pwdNonASCII}},
		{name: "TooShort_NoDigit_NoSpecial_HasSpace", password: "Pas wor", wantErrs: []string{pwdTooShort, pwdNoDigit, pwdNoSpecial, pwdHasSpace}},
		{name: "TooShort_NoDigit_NoSpecial_NonASCII", password: "P£sswor", wantErrs: []string{pwdTooShort, pwdNoDigit, pwdNoSpecial, pwdNonASCII}},
		{name: "TooShort_NoDigit_HasSpace_NonASCII", password: "P£$$ wo", wantErrs: []string{pwdTooShort, pwdNoDigit, pwdHasSpace, pwdNonASCII}},
		{name: "TooShort_NoSpecial_HasSpace_NonASCII", password: "P£55 wo", wantErrs: []string{pwdTooShort, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoLower_NoUpper_NoDigit", password: "!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit}},
		{name: "TooLong_NoLower_NoUpper_NoSpecial", password: "14211111421111142111114211111421111142111114211114211111421111142", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoSpecial}},
		{name: "TooLong_NoLower_NoUpper_HasSpace", password: "142 !@$142 !@$142 !@$142 !@$142 !@$142 !@$142 !@$142 !@$142 !@$14", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdHasSpace}},
		{name: "TooLong_NoLower_NoUpper_NonASCII", password: "14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNonASCII}},
		{name: "TooLong_NoLower_NoDigit_NoSpecial", password: "PASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPA", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdNoSpecial}},
		{name: "TooLong_NoLower_NoDigit_HasSpace", password: "PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PA", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdHasSpace}},
		{name: "TooLong_NoLower_NoDigit_NonASCII", password: "P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdNonASCII}},
		{name: "TooLong_NoLower_NoSpecial_HasSpace", password: "PAS 123PAS 123PAS 123PAS 123PAS 123PAS 123PAS 123PAS 123PAS 123PA", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoSpecial, pwdHasSpace}},
		{name: "TooLong_NoLower_NoSpecial_NonASCII", password: "P£SS123P£SS123P£SS123P£SS123P£SS123P£SS123P£SS123P£SS123P£SS123P£", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoSpecial, pwdNonASCII}},
		{name: "TooLong_NoLower_HasSpace_NonASCII", password: "P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£", wantErrs: []string{pwdTooLong, pwdNoLower, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoUpper_NoDigit_NoSpecial", password: "passssspassssspassssspassssspassssspassssspassssspassssspassssspa", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNoSpecial}},
		{name: "TooLong_NoUpper_NoDigit_HasSpace", password: "pas $$$pas $$$pas $$$pas $$$pas $$$pas $$$pas $$$pas $$$pas $$$pa", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdHasSpace}},
		{name: "TooLong_NoUpper_NoDigit_NonASCII", password: "p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNonASCII}},
		{name: "TooLong_NoUpper_NoSpecial_HasSpace", password: "pas 123pas 123pas 123pas 123pas 123pas 123pas 123pas 123pas 123pa", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoSpecial, pwdHasSpace}},
		{name: "TooLong_NoUpper_NoSpecial_NonASCII", password: "p£ss123p£ss123p£ss123p£ss123p£ss123p£ss123p£ss123p£ss123p£p£ss123", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoSpecial, pwdNonASCII}},
		{name: "TooLong_NoUpper_HasSpace_NonASCII", password: "p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoDigit_NoSpecial_HasSpace", password: "Pas worPas worPas worPas worPas worPas worPas worPas worPas worPa", wantErrs: []string{pwdTooLong, pwdNoDigit, pwdNoSpecial, pwdHasSpace}},
		{name: "TooLong_NoDigit_NoSpecial_NonASCII", password: "P£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£", wantErrs: []string{pwdTooLong, pwdNoDigit, pwdNoSpecial, pwdNonASCII}},
		{name: "TooLong_NoDigit_HasSpace_NonASCII", password: "P£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£", wantErrs: []string{pwdTooLong, pwdNoDigit, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoSpecial_HasSpace_NonASCII", password: "P£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP£", wantErrs: []string{pwdTooLong, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "NoLower_NoUpper_NoDigit_HasSpace", password: "!!!  !!!", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace}},
		{name: "NoLower_NoUpper_NoDigit_NonASCII", password: "!!!££!!!", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNonASCII}},
		{name: "NoLower_NoUpper_NoSpecial_HasSpace", password: "333  333", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace}},
		{name: "NoLower_NoUpper_NoSpecial_NonASCII", password: "333££333", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdNonASCII}},
		{name: "NoLower_NoUpper_HasSpace_NonASCII", password: "4£4  !!!", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdHasSpace, pwdNonASCII}},
		{name: "NoLower_NoDigit_NoSpecial_HasSpace", password: "MYP  SWO", wantErrs: []string{pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace}},
		{name: "NoLower_NoDigit_NoSpecial_NonASCII", password: "MYP££SWO", wantErrs: []string{pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdNonASCII}},
		{name: "NoLower_NoDigit_HasSpace_NonASCII", password: "M£P@SS !", wantErrs: []string{pwdNoLower, pwdNoDigit, pwdHasSpace, pwdNonASCII}},
		{name: "NoLower_NoSpecial_HasSpace_NonASCII", password: "M£P4  W1", wantErrs: []string{pwdNoLower, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "NoUpper_NoDigit_NoSpecial_HasSpace", password: "myp  swo", wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace}},
		{name: "NoUpper_NoDigit_NoSpecial_NonASCII", password: "myp££swo", wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII}},
		{name: "NoUpper_NoDigit_HasSpace_NonASCII", password: "m£p@ss !", wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII}},
		{name: "NoUpper_NoSpecial_HasSpace_NonASCII", password: "m£p4  w1", wantErrs: []string{pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "NoDigit_NoSpecial_HasSpace_NonASCII", password: "MyP£ss o", wantErrs: []string{pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},

		// 5-error cases
		{name: "TooShort_NoLower_NoUpper_NoDigit_HasSpace", password: "!@   $!", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace}},
		{name: "TooShort_NoLower_NoUpper_NoDigit_NonASCII", password: "!@£££$!", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNonASCII}},
		{name: "TooShort_NoLower_NoUpper_NoSpecial_HasSpace", password: "14   11", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace}},
		{name: "TooShort_NoLower_NoUpper_NoSpecial_NonASCII", password: "14£££11", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdNonASCII}},
		{name: "TooShort_NoLower_NoUpper_HasSpace_NonASCII", password: "1£2 !@$", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdHasSpace, pwdNonASCII}},
		{name: "TooShort_NoLower_NoDigit_NoSpecial_HasSpace", password: "PAS SSS", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace}},
		{name: "TooShort_NoLower_NoDigit_NoSpecial_NonASCII", password: "PAS£SSS", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdNonASCII}},
		{name: "TooShort_NoLower_NoDigit_HasSpace_NonASCII", password: "P£S !@$", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdHasSpace, pwdNonASCII}},
		{name: "TooShort_NoLower_NoSpecial_HasSpace_NonASCII", password: "P£S 123", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "TooShort_NoUpper_NoDigit_NoSpecial_HasSpace", password: "pas sss", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace}},
		{name: "TooShort_NoUpper_NoDigit_NoSpecial_NonASCII", password: "pas£sss", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII}},
		{name: "TooShort_NoUpper_NoDigit_HasSpace_NonASCII", password: "p£s $$$", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII}},
		{name: "TooShort_NoUpper_NoSpecial_HasSpace_NonASCII", password: "p£s 123", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "TooShort_NoDigit_NoSpecial_HasSpace_NonASCII", password: "P£s wor", wantErrs: []string{pwdTooShort, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoLower_NoUpper_NoDigit_HasSpace", password: "!@   $!!@   $!!@   $!!@   $!!@   $!!@   $!!@   $!!@   $!!@   $!!@", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace}},
		{name: "TooLong_NoLower_NoUpper_NoDigit_NonASCII", password: "!@£££$!!@£££$!!@£££$!!@£££$!!@£££$!!@£££$!!@£££$!!@£££$!!@£££$!!@", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNonASCII}},
		{name: "TooLong_NoLower_NoUpper_NoSpecial_HasSpace", password: "14   1114   1114   1114   1114   1114   1114   1114   1114   1114", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace}},
		{name: "TooLong_NoLower_NoUpper_NoSpecial_NonASCII", password: "14£££1114£££1114£££1114£££1114£££1114£££1114£££1114£££1114£££1114", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdNonASCII}},
		{name: "TooLong_NoLower_NoUpper_HasSpace_NonASCII", password: "1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoLower_NoDigit_NoSpecial_HasSpace", password: "PAS SSSPAS SSSPAS SSSPAS SSSPAS SSSPAS SSSPAS SSSPAS SSSPAS SSSPA", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace}},
		{name: "TooLong_NoLower_NoDigit_NoSpecial_NonASCII", password: "PAS£SSSPAS£SSSPAS£SSSPAS£SSSPAS£SSSPAS£SSSPAS£SSSPAS£SSSPAS£SSSPA", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdNonASCII}},
		{name: "TooLong_NoLower_NoDigit_HasSpace_NonASCII", password: "P£S !@$P£S !@$P£S !@$P£S !@$P£S !@$P£S !@$P£S !@$P£S !@$P£S !@$P£", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoLower_NoSpecial_HasSpace_NonASCII", password: "P£S 123P£S 123P£S 123P£S 123P£S 123P£S 123P£S 123P£S 123P£S 123P£", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoUpper_NoDigit_NoSpecial_HasSpace", password: "pas ssspas ssspas ssspas ssspas ssspas ssspas ssspas ssspas ssspa", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace}},
		{name: "TooLong_NoUpper_NoDigit_NoSpecial_NonASCII", password: "pas£ssspas£ssspas£ssspas£ssspas£ssspas£ssspas£ssspas£ssspas£ssspa", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII}},
		{name: "TooLong_NoUpper_NoDigit_HasSpace_NonASCII", password: "p£s $$$p£s $$$p£s $$$p£s $$$p£s $$$p£s $$$p£s $$$p£s $$$p£s $$$p£", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoUpper_NoSpecial_HasSpace_NonASCII", password: "p£s 123p£s 123p£s 123p£s 123p£s 123p£s 123p£s 123p£s 123p£s 123p£", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoDigit_NoSpecial_HasSpace_NonASCII", password: "P£s worP£s worP£s worP£s worP£s worP£s worP£s worP£s worP£s worP£", wantErrs: []string{pwdTooLong, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "NoLower_NoUpper_NoDigit_NoSpecial_HasSpace", password: "        ", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace}},
		{name: "NoLower_NoUpper_NoDigit_NoSpecial_NonASCII", password: "££££££££", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII}},
		{name: "NoLower_NoUpper_NoDigit_HasSpace_NonASCII", password: "!£!  !!!", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII}},
		{name: "NoLower_NoUpper_NoSpecial_HasSpace_NonASCII", password: "3£3  333", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "NoLower_NoDigit_NoSpecial_HasSpace_NonASCII", password: "M£P  SWO", wantErrs: []string{pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "NoUpper_NoDigit_NoSpecial_HasSpace_NonASCII", password: "m£p  swo", wantErrs: []string{pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},

		// 6-error cases
		{name: "TooShort_NoLower_NoUpper_NoDigit_NoSpecial_HasSpace", password: "       ", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace}},
		{name: "TooShort_NoLower_NoUpper_NoDigit_NoSpecial_NonASCII", password: "£££££££", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII}},
		{name: "TooShort_NoLower_NoUpper_NoDigit_HasSpace_NonASCII", password: "!£   $!", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII}},
		{name: "TooShort_NoLower_NoUpper_NoSpecial_HasSpace_NonASCII", password: "1£   11", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "TooShort_NoLower_NoDigit_NoSpecial_HasSpace_NonASCII", password: "P£S SSS", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "TooShort_NoUpper_NoDigit_NoSpecial_HasSpace_NonASCII", password: "p£s sss", wantErrs: []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoLower_NoUpper_NoDigit_NoSpecial_HasSpace", password: "                                                                 ", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace}},
		{name: "TooLong_NoLower_NoUpper_NoDigit_NoSpecial_NonASCII", password: "£££££££££££££££££££££££££££££££££££££££££££££££££££££££££££££££££", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII}},
		{name: "TooLong_NoLower_NoUpper_NoDigit_HasSpace_NonASCII", password: "!£   $!!£   $!!£   $!!£   $!!£   $!!£   $!!£   $!!£   $!!£   $!!£", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoLower_NoUpper_NoSpecial_HasSpace_NonASCII", password: "1£   111£   111£   111£   111£   111£   111£   111£   111£   111£", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoLower_NoDigit_NoSpecial_HasSpace_NonASCII", password: "P£S SSSP£S SSSP£S SSSP£S SSSP£S SSSP£S SSSP£S SSSP£S SSSP£S SSSP£", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoUpper_NoDigit_NoSpecial_HasSpace_NonASCII", password: "p£s sssp£s sssp£s sssp£s sssp£s sssp£s sssp£s sssp£s sssp£s sssp£", wantErrs: []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "NoLower_NoUpper_NoDigit_NoSpecial_HasSpace_NonASCII", password: "   ££   ", wantErrs: []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},

		// 7-error cases
		{name: "TooShort_NoLower_NoUpper_NoDigit_NoSpecial_HasSpace_NonASCII", password: "   £   ", wantErrs: []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
		{name: "TooLong_NoLower_NoUpper_NoDigit_NoSpecial_HasSpace_NonASCII", password: "   £      £      £      £      £      £      £      £      £     ", wantErrs: []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}},
	} {
		t.Run(c.name, func(t *testing.T) {
			gotErrs := sut.Validate(c.password)
			assert.EqualArr(t, c.wantErrs, gotErrs)
		})
	}
}
