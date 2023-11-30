//go:build itest

package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	registerAPI "github.com/kxplxn/goteam/internal/api/register"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/auth"
	teamTable "github.com/kxplxn/goteam/pkg/dbaccess/team"
	userTable "github.com/kxplxn/goteam/pkg/dbaccess/user"
	"github.com/kxplxn/goteam/pkg/log"
)

// TestRegisterHandler tests the http.Handler for the register API route and
// asserts that it behaves correctly during various execution paths.
func TestRegisterHandler(t *testing.T) {
	sut := registerAPI.NewPOSTHandler(
		registerAPI.NewUserValidator(
			registerAPI.NewUsernameValidator(),
			registerAPI.NewPasswordValidator(),
		),
		registerAPI.NewInviteCodeValidator(),
		teamTable.NewSelectorByInvCode(db),
		userTable.NewSelector(db),
		registerAPI.NewPasswordHasher(),
		userTable.NewInserter(db),
		auth.NewJWTGenerator(jwtKey),
		log.New(),
	)

	// Used in status 400 cases to assert on username and password error messages.
	assertOnValidationErrs := func(
		wantUsernameErrs, wantPasswordErrs []string,
	) func(*testing.T, *http.Response, string) {
		return func(t *testing.T, res *http.Response, _ string) {
			var resBody registerAPI.POSTResp
			if err := json.NewDecoder(res.Body).Decode(
				&resBody,
			); err != nil {
				t.Fatal(err)
			}
			assert.AllEqual(t.Error,
				resBody.ValidationErrs.Username,
				wantUsernameErrs,
			)
			assert.AllEqual(t.Error,
				resBody.ValidationErrs.Password,
				wantPasswordErrs,
			)
		}
	}

	assertOnResErr := func(
		wantErrMsg string,
	) func(*testing.T, *http.Response, string) {
		return func(t *testing.T, res *http.Response, _ string) {
			var resBody registerAPI.POSTResp
			if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
				t.Fatal(err)
			}
			assert.Equal(t.Error, resBody.Err, wantErrMsg)
		}
	}

	for _, c := range []struct {
		name           string
		username       string
		password       string
		inviteCode     string
		wantStatusCode int
		assertFunc     func(*testing.T, *http.Response, string)
	}{
		{
			name:           "UsnEmpty,PwdEmpty",
			username:       "",
			password:       "",
			inviteCode:     "",
			wantStatusCode: http.StatusBadRequest,
			assertFunc: assertOnValidationErrs(
				[]string{"Username cannot be empty."},
				[]string{"Password cannot be empty."},
			),
		},
		{
			name: "UsnTooShort,UsnInvalidChar,PwdTooShort,PwdNoLower," +
				"PwdNoDigit,PwdNoSpecial",
			username:       "bob!",
			password:       "PASSSSS",
			inviteCode:     "",
			wantStatusCode: http.StatusBadRequest,
			assertFunc: assertOnValidationErrs(
				[]string{
					"Username cannot be shorter than 5 characters.",
					"Username can contain only letters (a-z/A-Z) and digits " +
						"(0-9).",
				},
				[]string{
					"Password cannot be shorter than 8 characters.",
					"Password must contain a lowercase letter (a-z).",
					"Password must contain a digit (0-9).",
					"Password must contain one of the following special " +
						"characters: ! \" # $ % & ' ( ) * + , - . / : ; < = " +
						"> ? [ \\ ] ^ _ ` { | } ~.",
				},
			),
		},
		{
			name: "UsnTooLong,UsnDigitStart,PwdTooLong,PwdNoUpper," +
				"PwdHasSpace,PwdNonASCII",
			username: "1bobobobobobobobo",
			password: "p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p£$ 123p" +
				"£$ 123p£$ 123p£$ 123p£",
			wantStatusCode: http.StatusBadRequest,
			assertFunc: assertOnValidationErrs(
				[]string{
					"Username cannot be longer than 15 characters.",
					"Username can start only with a letter (a-z/A-Z).",
				},
				[]string{
					"Password cannot be longer than 64 characters.",
					"Password must contain an uppercase letter (A-Z).",
					"Password cannot contain spaces.",
					"Password can contain only letters (a-z/A-Z), digits " +
						"(0-9), and the following special characters: " +
						"! \" # $ % & ' ( ) * + , - . / : ; < = > ? [ \\ ] ^ " +
						"_ ` { | } ~.",
				},
			),
		},
		{
			name:           "UsnTaken",
			username:       "team1Member",
			password:       "Myp4ssw0rd!",
			inviteCode:     "",
			wantStatusCode: http.StatusBadRequest,
			assertFunc: assertOnValidationErrs(
				[]string{"Username is already taken."}, []string{},
			),
		},
		{
			name:           "InviteCodeInvalid",
			username:       "bob321",
			password:       "Myp4ssw0rd!",
			inviteCode:     "10249812049182",
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     assertOnResErr("Invalid invite code."),
		},
		{
			name:           "TeamNotFound",
			username:       "bob321",
			password:       "Myp4ssw0rd!",
			inviteCode:     "ca9512b0-d448-46da-9a25-2a6b7d4b405e",
			wantStatusCode: http.StatusNotFound,
			assertFunc:     assertOnResErr("Team not found."),
		},
		{
			name:           "Success",
			username:       "bob321",
			password:       "Myp4ssw0rd!",
			inviteCode:     "",
			wantStatusCode: http.StatusOK,
			assertFunc: func(t *testing.T, res *http.Response, _ string) {
				// assert that a new user is inserted into the database with
				// the correct credentials
				var password string
				err := db.QueryRow(
					`SELECT password FROM app."user" WHERE username = $1`,
					"bob321",
				).Scan(&password)
				if err != nil {
					t.Fatal(err)
				}
				if err = bcrypt.CompareHashAndPassword(
					[]byte(password), []byte("Myp4ssw0rd!"),
				); err != nil {
					t.Error(err)
				}

				// assert that the returned JWT is valid and has the correct
				// subject
				cookie := res.Cookies()[0]
				assert.True(t.Error, cookie.Secure)
				assert.Equal(t.Error, cookie.SameSite, http.SameSiteNoneMode)
				claims := jwt.RegisteredClaims{}
				if _, err = jwt.ParseWithClaims(
					cookie.Value, &claims, func(token *jwt.Token) (any, error) {
						return []byte(jwtKey), nil
					},
				); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t.Error, claims.Subject, "bob321")
			},
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			reqBody, err := json.Marshal(map[string]string{
				"username": c.username,
				"password": c.password,
			})
			if err != nil {
				t.Fatal(err)
			}
			req, err := http.NewRequest(
				http.MethodPost,
				"?inviteCode="+c.inviteCode,
				bytes.NewReader(reqBody),
			)
			if err != nil {
				t.Fatal(err)
			}
			w := httptest.NewRecorder()

			sut.Handle(w, req, "")

			res := w.Result()

			assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)
			// Run case-specific assertions.
			c.assertFunc(t, res, "")
		})
	}
}
