package register

import (
	"net/http"
	"testing"

	"github.com/kxplxn/goteam/server-v2/test"
)

// TestRegister perfomes functional tests on the register endpoint via the
// Handler.
func TestRegister(t *testing.T) {
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

	reqBody := make(map[string]string)
	reqBody["username"] = "bobby"
	reqBody["password"] = "S3curePa$$"
	resBody := &ResBody{}

	test.NewRoute("/register", http.MethodPost, NewHandler(), reqBody, resBody, []*test.RouteSuite{
		test.NewRouteSuite("UsernameValidation", "username", []*test.RouteCase{
			// 1-error cases
			test.NewRouteCase("Empty", "", []string{usnEmpty}),
			test.NewRouteCase("TooShort", "bob1", []string{usnTooShort}),
			test.NewRouteCase("TooLong", "bobobobobobobobob", []string{usnTooLong}),
			test.NewRouteCase("InvalidCharacter", "bobob!", []string{usnInvalidChar}),
			test.NewRouteCase("DigitStart", "1bobob", []string{usnDigitStart}),

			// 2-error cases
			test.NewRouteCase("TooShort_InvalidCharacter", "bob!", []string{usnTooShort, usnInvalidChar}),
			test.NewRouteCase("TooShort_DigitStart", "1bob", []string{usnTooShort, usnDigitStart}),
			test.NewRouteCase("TooLong_InvalidCharacter", "bobobobobobobobo!", []string{usnTooLong, usnInvalidChar}),
			test.NewRouteCase("TooLong_DigitStart", "1bobobobobobobobo", []string{usnTooLong, usnDigitStart}),
			test.NewRouteCase("InvalidCharacter_DigitStart", "1bob!", []string{usnInvalidChar, usnDigitStart}),

			// 3-error cases
			test.NewRouteCase("TooShort_InvalidCharacter_DigitStart", "1bo!", []string{usnTooShort, usnInvalidChar, usnDigitStart}),
			test.NewRouteCase("TooLong_InvalidCharacter_DigitStart", "1bobobobobobobob!", []string{usnTooLong, usnInvalidChar, usnDigitStart}),
		}, http.StatusBadRequest),

		test.NewRouteSuite("PasswordValidation", "password", []*test.RouteCase{
			// 1-error cases
			test.NewRouteCase("Empty", "", []string{pwdEmpty}),
			test.NewRouteCase("TooShort", "Myp4ss!", []string{pwdTooShort}),
			test.NewRouteCase("TooLong", "Myp4sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong}),
			test.NewRouteCase("NoLower", "MY4LLUPPERPASSWORD!", []string{pwdNoLower}),
			test.NewRouteCase("NoUpper", "my4lllowerpassword!", []string{pwdNoUpper}),
			test.NewRouteCase("NoDigit", "myNOdigitPASSWORD!", []string{pwdNoDigit}),
			test.NewRouteCase("NoSpecial", "myNOspecialP4SSWORD", []string{pwdNoSpecial}),
			test.NewRouteCase("HasSpace", "my SP4CED p4ssword !", []string{pwdHasSpace}),
			test.NewRouteCase("NonASCII", "myNØNÅSCÎÎp4ssword!", []string{pwdNonASCII}),

			// 2-error cases
			test.NewRouteCase("TooShort_NoLower", "MYP4SS!", []string{pwdTooShort, pwdNoLower}),
			test.NewRouteCase("TooShort_NoUpper", "myp4ss!", []string{pwdTooShort, pwdNoUpper}),
			test.NewRouteCase("TooShort_NoDigit", "MyPass!", []string{pwdTooShort, pwdNoDigit}),
			test.NewRouteCase("TooShort_NoSpecial", "MyP4ssw", []string{pwdTooShort, pwdNoSpecial}),
			test.NewRouteCase("TooShort_HasSpace", "My P4s!", []string{pwdTooShort, pwdHasSpace}),
			test.NewRouteCase("TooShort_NonASCII", "M¥P4s!2", []string{pwdTooShort, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower", "MYP4SSWORDWH!CHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH", []string{pwdTooLong, pwdNoLower}),
			test.NewRouteCase("TooLong_NoUpper", "myp4sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong, pwdNoUpper}),
			test.NewRouteCase("TooLong_NoDigit", "Mypasswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong, pwdNoDigit}),
			test.NewRouteCase("TooLong_NoSpecial", "Myp4sswordwhichislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong, pwdNoSpecial}),
			test.NewRouteCase("TooLong_HasSpace", "Myp4sswo   rdwh!chislongandimeanreallylongforsomereasonohiknowwhy", []string{pwdTooLong, pwdHasSpace}),
			test.NewRouteCase("TooLong_NonASCII", "Myp4££wordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoUpper", "4444!!!!", []string{pwdNoLower, pwdNoUpper}),
			test.NewRouteCase("NoLower_NoDigit", "MYP@SSW!", []string{pwdNoLower, pwdNoDigit}),
			test.NewRouteCase("NoLower_NoSpecial", "MYP4SSW1", []string{pwdNoLower, pwdNoSpecial}),
			test.NewRouteCase("NoLower_HasSpace", "MYP4SS !", []string{pwdNoLower, pwdHasSpace}),
			test.NewRouteCase("NoLower_NonASCII", "MYP4££W!", []string{pwdNoLower, pwdNonASCII}),
			test.NewRouteCase("NoUpper_NoDigit", "myp@ssw!", []string{pwdNoUpper, pwdNoDigit}),
			test.NewRouteCase("NoUpper_NoSpecial", "myp4ssw1", []string{pwdNoUpper, pwdNoSpecial}),
			test.NewRouteCase("NoUpper_HasSpace", "myp4ss !", []string{pwdNoUpper, pwdHasSpace}),
			test.NewRouteCase("NoUpper_NonASCII", "myp4££w!", []string{pwdNoUpper, pwdNonASCII}),
			test.NewRouteCase("NoDigit_NoSpecial", "MyPasswd", []string{pwdNoDigit, pwdNoSpecial}),
			test.NewRouteCase("NoDigit_HasSpace", "MyPass !", []string{pwdNoDigit, pwdHasSpace}),
			test.NewRouteCase("NoDigit_NonASCII", "MyPa££w!", []string{pwdNoDigit, pwdNonASCII}),
			test.NewRouteCase("NoSpecial_HasSpace", "My  P4ss", []string{pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("NoSpecial_NonASCII", "MyPa££w1", []string{pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("HasSpace_NonASCII", "MyP4££ !", []string{pwdHasSpace, pwdNonASCII}),

			// 3-error cases
			test.NewRouteCase("TooShort_NoLower_NoUpper", "1421!@$", []string{pwdTooShort, pwdNoLower, pwdNoUpper}),
			test.NewRouteCase("TooShort_NoLower_NoDigit", "PASS!@$", []string{pwdTooShort, pwdNoLower, pwdNoDigit}),
			test.NewRouteCase("TooShort_NoLower_NoSpecial", "PASS123", []string{pwdTooShort, pwdNoLower, pwdNoSpecial}),
			test.NewRouteCase("TooShort_NoLower_HasSpace", "PA$ 123", []string{pwdTooShort, pwdNoLower, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoLower_NonASCII", "PA$£123", []string{pwdTooShort, pwdNoLower, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoUpper_NoDigit", "pass$$$", []string{pwdTooShort, pwdNoUpper, pwdNoDigit}),
			test.NewRouteCase("TooShort_NoUpper_NoSpecial", "pass123", []string{pwdTooShort, pwdNoUpper, pwdNoSpecial}),
			test.NewRouteCase("TooShort_NoUpper_HasSpace", "pa$ 123", []string{pwdTooShort, pwdNoUpper, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoUpper_NonASCII", "pa$£123", []string{pwdTooShort, pwdNoUpper, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoDigit_NoSpecial", "Passwor", []string{pwdTooShort, pwdNoDigit, pwdNoSpecial}),
			test.NewRouteCase("TooShort_NoDigit_HasSpace", "Pa$$ wo", []string{pwdTooShort, pwdNoDigit, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoDigit_NonASCII", "Pa$$£wo", []string{pwdTooShort, pwdNoDigit, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoSpecial_HasSpace", "Pa55 wo", []string{pwdTooShort, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoSpecial_NonASCII", "Pa55£wo", []string{pwdTooShort, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooShort_HasSpace_NonASCII", "P4$ £wo", []string{pwdTooShort, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoUpper", "111422222222!3333333333333333333333333333333333333333333333333333", []string{pwdTooLong, pwdNoLower, pwdNoUpper}),
			test.NewRouteCase("TooLong_NoLower_NoDigit", "MYPASSWORDWH!CHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH", []string{pwdTooLong, pwdNoLower, pwdNoDigit}),
			test.NewRouteCase("TooLong_NoLower_NoSpecial", "MYP4SSWORDWHICHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH", []string{pwdTooLong, pwdNoLower, pwdNoSpecial}),
			test.NewRouteCase("TooLong_NoLower_HasSpace", "MYP4SS    WH!CHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH", []string{pwdTooLong, pwdNoLower, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoLower_NonASCII", "£YP4SSWORDWH!CHISLONGANDIMEANREALLYLONGFORSOMEREASONOHIKNOWWHYTBH", []string{pwdTooLong, pwdNoLower, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoUpper_NoDigit", "mypasswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong, pwdNoUpper, pwdNoDigit}),
			test.NewRouteCase("TooLong_NoUpper_NoSpecial", "myp4sswordwhichislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong, pwdNoUpper, pwdNoSpecial}),
			test.NewRouteCase("TooLong_NoUpper_HasSpace", "myp4ss    wh!chislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong, pwdNoUpper, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoUpper_NonASCII", "£yp4sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong, pwdNoUpper, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoDigit_NoSpecial", "Mypasswordwhichislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong, pwdNoDigit, pwdNoSpecial}),
			test.NewRouteCase("TooLong_NoDigit_HasSpace", "Mypass    wh!chislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong, pwdNoDigit, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoDigit_NonASCII", "Myp£sswordwh!chislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong, pwdNoDigit, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoSpecial_HasSpace", "Myp4ss    whichislongandimeanreallylongforsomereasonohiknowwhytbh", []string{pwdTooLong, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooLong_HasSpace_NonASCII", "Myp4ssw£   rdwh!chislongandimeanreallylongforsomereasonohiknowwhy", []string{pwdTooLong, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoUpper_NoDigit", "!!!!!!!!", []string{pwdNoLower, pwdNoUpper, pwdNoDigit}),
			test.NewRouteCase("NoLower_NoUpper_NoSpecial", "33333333", []string{pwdNoLower, pwdNoUpper, pwdNoSpecial}),
			test.NewRouteCase("NoLower_NoUpper_HasSpace", "444  !!!", []string{pwdNoLower, pwdNoUpper, pwdHasSpace}),
			test.NewRouteCase("NoLower_NoUpper_NonASCII", "£££444!!", []string{pwdNoLower, pwdNoUpper, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoDigit_NoSpecial", "MYPASSWO", []string{pwdNoLower, pwdNoDigit, pwdNoSpecial}),
			test.NewRouteCase("NoLower_NoDigit_HasSpace", "MYP@SS !", []string{pwdNoLower, pwdNoDigit, pwdHasSpace}),
			test.NewRouteCase("NoLower_NoDigit_NonASCII", "M£P@SSW!", []string{pwdNoLower, pwdNoDigit, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoSpecial_HasSpace", "MYP4  W1", []string{pwdNoLower, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("NoLower_NoSpecial_NonASCII", "M£P4SSW1", []string{pwdNoLower, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("NoLower_HasSpace_NonASCII", "M£P4SS !", []string{pwdNoLower, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoUpper_NoDigit_NoSpecial", "mypasswo", []string{pwdNoUpper, pwdNoDigit, pwdNoSpecial}),
			test.NewRouteCase("NoUpper_NoDigit_HasSpace", "myp@ss !", []string{pwdNoUpper, pwdNoDigit, pwdHasSpace}),
			test.NewRouteCase("NoUpper_NoDigit_NonASCII", "m£p@ssw!", []string{pwdNoUpper, pwdNoDigit, pwdNonASCII}),
			test.NewRouteCase("NoUpper_NoSpecial_HasSpace", "myp4ss 1", []string{pwdNoUpper, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("NoUpper_NoSpecial_NonASCII", "m£p4ssw1", []string{pwdNoUpper, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("NoUpper_HasSpace_NonASCII", "m£p4ss !", []string{pwdNoUpper, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoDigit_NoSpecial_HasSpace", "MyPass o", []string{pwdNoDigit, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("NoDigit_NoSpecial_NonASCII", "M£Passwd", []string{pwdNoDigit, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("NoDigit_HasSpace_NonASCII", "M£Pass !", []string{pwdNoDigit, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoSpecial_HasSpace_NonASCII", "M£  P4ss", []string{pwdNoSpecial, pwdHasSpace, pwdNonASCII}),

			// 4-error cases
			test.NewRouteCase("TooShort_NoLower_NoUpper_NoDigit", "!@$!@$!", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit}),
			test.NewRouteCase("TooShort_NoLower_NoUpper_NoSpecial", "1421111", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoSpecial}),
			test.NewRouteCase("TooShort_NoLower_NoUpper_HasSpace", "142 !@$", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoLower_NoUpper_NonASCII", "14£1!@$", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoLower_NoDigit_NoSpecial", "PASSSSS", []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdNoSpecial}),
			test.NewRouteCase("TooShort_NoLower_NoDigit_HasSpace", "PAS !@$", []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoLower_NoDigit_NonASCII", "P£SS!@$", []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoLower_NoSpecial_HasSpace", "PAS 123", []string{pwdTooShort, pwdNoLower, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoLower_NoSpecial_NonASCII", "P£SS123", []string{pwdTooShort, pwdNoLower, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoLower_HasSpace_NonASCII", "P£$ 123", []string{pwdTooShort, pwdNoLower, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoUpper_NoDigit_NoSpecial", "passsss", []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNoSpecial}),
			test.NewRouteCase("TooShort_NoUpper_NoDigit_HasSpace", "pas $$$", []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoUpper_NoDigit_NonASCII", "p£ss$$$", []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoUpper_NoSpecial_HasSpace", "pas 123", []string{pwdTooShort, pwdNoUpper, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoUpper_NoSpecial_NonASCII", "p£ss123", []string{pwdTooShort, pwdNoUpper, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoUpper_HasSpace_NonASCII", "p£$ 123", []string{pwdTooShort, pwdNoUpper, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoDigit_NoSpecial_HasSpace", "Pas wor", []string{pwdTooShort, pwdNoDigit, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoDigit_NoSpecial_NonASCII", "P£sswor", []string{pwdTooShort, pwdNoDigit, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoDigit_HasSpace_NonASCII", "P£$$ wo", []string{pwdTooShort, pwdNoDigit, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoSpecial_HasSpace_NonASCII", "P£55 wo", []string{pwdTooShort, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_NoDigit", "!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@$!@$!!@", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_NoSpecial", "14211111421111142111114211111421111142111114211114211111421111142", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoSpecial}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_HasSpace", "142 !@$142 !@$142 !@$142 !@$142 !@$142 !@$142 !@$142 !@$142 !@$14", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_NonASCII", "14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14£1!@$14", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoDigit_NoSpecial", "PASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPASSSSSPA", []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdNoSpecial}),
			test.NewRouteCase("TooLong_NoLower_NoDigit_HasSpace", "PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PAS !@$PA", []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoLower_NoDigit_NonASCII", "P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£SS!@$P£", []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoSpecial_HasSpace", "PAS 123PAS 123PAS 123PAS 123PAS 123PAS 123PAS 123PAS 123PAS 123PA", []string{pwdTooLong, pwdNoLower, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoLower_NoSpecial_NonASCII", "P£SS123P£SS123P£SS123P£SS123P£SS123P£SS123P£SS123P£SS123P£SS123P£", []string{pwdTooLong, pwdNoLower, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_HasSpace_NonASCII", "P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£$ 123P£", []string{pwdTooLong, pwdNoLower, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoUpper_NoDigit_NoSpecial", "passssspassssspassssspassssspassssspassssspassssspassssspassssspa", []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNoSpecial}),
			test.NewRouteCase("TooLong_NoUpper_NoDigit_HasSpace", "pas $$$pas $$$pas $$$pas $$$pas $$$pas $$$pas $$$pas $$$pas $$$pa", []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoUpper_NoDigit_NonASCII", "p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£ss$$$p£", []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoUpper_NoSpecial_HasSpace", "pas 123pas 123pas 123pas 123pas 123pas 123pas 123pas 123pas 123pa", []string{pwdTooLong, pwdNoUpper, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoUpper_NoSpecial_NonASCII", "p£ss123p£ss123p£ss123p£ss123p£ss123p£ss123p£ss123p£ss123p£p£ss123", []string{pwdTooLong, pwdNoUpper, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoUpper_HasSpace_NonASCII", "p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£", []string{pwdTooLong, pwdNoUpper, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoDigit_NoSpecial_HasSpace", "Pas worPas worPas worPas worPas worPas worPas worPas worPas worPa", []string{pwdTooLong, pwdNoDigit, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoDigit_NoSpecial_NonASCII", "P£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£ssworP£", []string{pwdTooLong, pwdNoDigit, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoDigit_HasSpace_NonASCII", "P£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£$$ woP£", []string{pwdTooLong, pwdNoDigit, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoSpecial_HasSpace_NonASCII", "P£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP£55 woP£", []string{pwdTooLong, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoUpper_NoDigit_HasSpace", "!!!  !!!", []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace}),
			test.NewRouteCase("NoLower_NoUpper_NoDigit_NonASCII", "!!!££!!!", []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoUpper_NoSpecial_HasSpace", "333  333", []string{pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("NoLower_NoUpper_NoSpecial_NonASCII", "333££333", []string{pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoUpper_HasSpace_NonASCII", "4£4  !!!", []string{pwdNoLower, pwdNoUpper, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoDigit_NoSpecial_HasSpace", "MYP  SWO", []string{pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("NoLower_NoDigit_NoSpecial_NonASCII", "MYP££SWO", []string{pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoDigit_HasSpace_NonASCII", "M£P@SS !", []string{pwdNoLower, pwdNoDigit, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoSpecial_HasSpace_NonASCII", "M£P4  W1", []string{pwdNoLower, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoUpper_NoDigit_NoSpecial_HasSpace", "myp  swo", []string{pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("NoUpper_NoDigit_NoSpecial_NonASCII", "myp££swo", []string{pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("NoUpper_NoDigit_HasSpace_NonASCII", "m£p@ss !", []string{pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoUpper_NoSpecial_HasSpace_NonASCII", "m£p4  w1", []string{pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoDigit_NoSpecial_HasSpace_NonASCII", "MyP£ss o", []string{pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),

			// 5-error cases
			test.NewRouteCase("TooShort_NoLower_NoUpper_NoDigit_HasSpace", "!@   $!", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoLower_NoUpper_NoDigit_NonASCII", "!@£££$!", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoLower_NoUpper_NoSpecial_HasSpace", "14   11", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoLower_NoUpper_NoSpecial_NonASCII", "14£££11", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoLower_NoUpper_HasSpace_NonASCII", "1£2 !@$", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoLower_NoDigit_NoSpecial_HasSpace", "PAS SSS", []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoLower_NoDigit_NoSpecial_NonASCII", "PAS£SSS", []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoLower_NoDigit_HasSpace_NonASCII", "P£S !@$", []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoLower_NoSpecial_HasSpace_NonASCII", "P£S 123", []string{pwdTooShort, pwdNoLower, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoUpper_NoDigit_NoSpecial_HasSpace", "pas sss", []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoUpper_NoDigit_NoSpecial_NonASCII", "pas£sss", []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoUpper_NoDigit_HasSpace_NonASCII", "p£s $$$", []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoUpper_NoSpecial_HasSpace_NonASCII", "p£s 123", []string{pwdTooShort, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoDigit_NoSpecial_HasSpace_NonASCII", "P£s wor", []string{pwdTooShort, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_NoDigit_HasSpace", "!@   $!!@   $!!@   $!!@   $!!@   $!!@   $!!@   $!!@   $!!@   $!!@", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_NoDigit_NonASCII", "!@£££$!!@£££$!!@£££$!!@£££$!!@£££$!!@£££$!!@£££$!!@£££$!!@£££$!!@", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_NoSpecial_HasSpace", "14   1114   1114   1114   1114   1114   1114   1114   1114   1114", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_NoSpecial_NonASCII", "14£££1114£££1114£££1114£££1114£££1114£££1114£££1114£££1114£££1114", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_HasSpace_NonASCII", "1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£2 !@$1£", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoDigit_NoSpecial_HasSpace", "PAS SSSPAS SSSPAS SSSPAS SSSPAS SSSPAS SSSPAS SSSPAS SSSPAS SSSPA", []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoLower_NoDigit_NoSpecial_NonASCII", "PAS£SSSPAS£SSSPAS£SSSPAS£SSSPAS£SSSPAS£SSSPAS£SSSPAS£SSSPAS£SSSPA", []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoDigit_HasSpace_NonASCII", "P£S !@$P£S !@$P£S !@$P£S !@$P£S !@$P£S !@$P£S !@$P£S !@$P£S !@$P£", []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoSpecial_HasSpace_NonASCII", "P£S 123P£S 123P£S 123P£S 123P£S 123P£S 123P£S 123P£S 123P£S 123P£", []string{pwdTooLong, pwdNoLower, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoUpper_NoDigit_NoSpecial_HasSpace", "pas ssspas ssspas ssspas ssspas ssspas ssspas ssspas ssspas ssspa", []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoUpper_NoDigit_NoSpecial_NonASCII", "pas£ssspas£ssspas£ssspas£ssspas£ssspas£ssspas£ssspas£ssspas£ssspa", []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoUpper_NoDigit_HasSpace_NonASCII", "p£s $$$p£s $$$p£s $$$p£s $$$p£s $$$p£s $$$p£s $$$p£s $$$p£s $$$p£", []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoUpper_NoSpecial_HasSpace_NonASCII", "p£s 123p£s 123p£s 123p£s 123p£s 123p£s 123p£s 123p£s 123p£s 123p£", []string{pwdTooLong, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoDigit_NoSpecial_HasSpace_NonASCII", "P£s worP£s worP£s worP£s worP£s worP£s worP£s worP£s worP£s worP£", []string{pwdTooLong, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoUpper_NoDigit_NoSpecial_HasSpace", "        ", []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("NoLower_NoUpper_NoDigit_NoSpecial_NonASCII", "££££££££", []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoUpper_NoDigit_HasSpace_NonASCII", "!£!  !!!", []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoUpper_NoSpecial_HasSpace_NonASCII", "3£3  333", []string{pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoDigit_NoSpecial_HasSpace_NonASCII", "M£P  SWO", []string{pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoUpper_NoDigit_NoSpecial_HasSpace_NonASCII", "m£p  swo", []string{pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),

			// 6-error cases
			test.NewRouteCase("TooShort_NoLower_NoUpper_NoDigit_NoSpecial_HasSpace", "       ", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooShort_NoLower_NoUpper_NoDigit_NoSpecial_NonASCII", "£££££££", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoLower_NoUpper_NoDigit_HasSpace_NonASCII", "!£   $!", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoLower_NoUpper_NoSpecial_HasSpace_NonASCII", "1£   11", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoLower_NoDigit_NoSpecial_HasSpace_NonASCII", "P£S SSS", []string{pwdTooShort, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooShort_NoUpper_NoDigit_NoSpecial_HasSpace_NonASCII", "p£s sss", []string{pwdTooShort, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_NoDigit_NoSpecial_HasSpace", "                                                                 ", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_NoDigit_NoSpecial_NonASCII", "£££££££££££££££££££££££££££££££££££££££££££££££££££££££££££££££££", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_NoDigit_HasSpace_NonASCII", "!£   $!!£   $!!£   $!!£   $!!£   $!!£   $!!£   $!!£   $!!£   $!!£", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_NoSpecial_HasSpace_NonASCII", "1£   111£   111£   111£   111£   111£   111£   111£   111£   111£", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoDigit_NoSpecial_HasSpace_NonASCII", "P£S SSSP£S SSSP£S SSSP£S SSSP£S SSSP£S SSSP£S SSSP£S SSSP£S SSSP£", []string{pwdTooLong, pwdNoLower, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoUpper_NoDigit_NoSpecial_HasSpace_NonASCII", "p£s sssp£s sssp£s sssp£s sssp£s sssp£s sssp£s sssp£s sssp£s sssp£", []string{pwdTooLong, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("NoLower_NoUpper_NoDigit_NoSpecial_HasSpace_NonASCII", "   ££   ", []string{pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),

			// 7-error cases
			test.NewRouteCase("TooShort_NoLower_NoUpper_NoDigit_NoSpecial_HasSpace_NonASCII", "   £   ", []string{pwdTooShort, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
			test.NewRouteCase("TooLong_NoLower_NoUpper_NoDigit_NoSpecial_HasSpace_NonASCII", "   £      £      £      £      £      £      £      £      £     ", []string{pwdTooLong, pwdNoLower, pwdNoUpper, pwdNoDigit, pwdNoSpecial, pwdHasSpace, pwdNonASCII}),
		}, http.StatusBadRequest),
	}).Run(t)
}
