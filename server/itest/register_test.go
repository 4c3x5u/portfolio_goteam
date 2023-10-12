//go:build itest

package itest

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	registerAPI "github.com/kxplxn/goteam/server/api/register"
	"github.com/kxplxn/goteam/server/assert"
	"github.com/kxplxn/goteam/server/auth"
	userTable "github.com/kxplxn/goteam/server/dbaccess/user"
	"github.com/kxplxn/goteam/server/log"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

// TestRegisterHandler tests the http.Handler for the register API route and
// asserts that it behaves correctly during various execution paths.
func TestRegisterHandler(t *testing.T) {
	sut := registerAPI.NewHandler(
		registerAPI.NewValidator(
			registerAPI.NewUsernameValidator(),
			registerAPI.NewPasswordValidator(),
		),
		userTable.NewSelector(db),
		registerAPI.NewPasswordHasher(),
		userTable.NewInserter(db),
		auth.NewJWTGenerator(jwtKey),
		log.New(),
	)

	// Used in status 400 cases to assert on username and password error messages.
	assertOnValidationErrs := func(
		wantUsernameErrs, wantPasswordErrs []string,
	) func(*testing.T, *http.Response) {
		return func(t *testing.T, res *http.Response) {
			var resBody registerAPI.ResBody
			if err := json.NewDecoder(res.Body).Decode(
				&resBody,
			); err != nil {
				t.Fatal(err)
			}
			if err := assert.EqualArr(
				wantUsernameErrs,
				resBody.Errs.Username,
			); err != nil {
				t.Error(err)
			}
			if err := assert.EqualArr(
				wantPasswordErrs,
				resBody.Errs.Password,
			); err != nil {
				t.Error(err)
			}
		}
	}

	for _, c := range []struct {
		name           string
		username       string
		password       string
		wantStatusCode int
		assertFunc     func(*testing.T, *http.Response)
	}{
		{
			name:           "UsnEmpty,PwdEmpty",
			username:       "",
			password:       "",
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
			username:       "bob123",
			password:       "Myp4ssw0rd!",
			wantStatusCode: http.StatusBadRequest,
			assertFunc: assertOnValidationErrs(
				[]string{"Username is already taken."}, []string{},
			),
		},
		{
			name:           "Success",
			username:       "bob321",
			password:       "Myp4ssw0rd!",
			wantStatusCode: http.StatusOK,
			assertFunc: func(t *testing.T, res *http.Response) {
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
				if err := assert.True(cookie.Secure); err != nil {
					t.Error(err)
				}
				if err := assert.Equal(
					http.SameSiteNoneMode, cookie.SameSite,
				); err != nil {
					t.Error(err)
				}

				claims := jwt.RegisteredClaims{}
				if _, err = jwt.ParseWithClaims(
					cookie.Value, &claims, func(token *jwt.Token) (any, error) {
						return []byte(jwtKey), nil
					},
				); err != nil {
					t.Fatal(err)
				}
				if err = assert.Equal("bob321", claims.Subject); err != nil {
					t.Error(err)
				}
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

			// Run case-specific assertions.
			c.assertFunc(t, res)
		})
	}
}
