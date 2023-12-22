//go:build itest

package api

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	registerAPI "github.com/kxplxn/goteam/internal/user/register"
	"github.com/kxplxn/goteam/pkg/assert"
	"github.com/kxplxn/goteam/pkg/db/usertable"
	"github.com/kxplxn/goteam/pkg/log"
	"github.com/kxplxn/goteam/pkg/token"
)

func TestRegisterAPI(t *testing.T) {
	sut := registerAPI.NewPostHandler(
		registerAPI.NewUserValidator(
			registerAPI.NewUsernameValidator(),
			registerAPI.NewPasswordValidator(),
		),
		token.DecodeInvite,
		registerAPI.NewPasswordHasher(),
		usertable.NewInserter(svcDynamo),
		token.EncodeAuth,
		log.New(),
	)

	// Used in status 400 cases to assert on username and password error messages.
	assertOnValidationErrs := func(
		wantUsernameErrs, wantPasswordErrs []string,
	) func(*testing.T, *http.Response, string) {
		return func(t *testing.T, res *http.Response, _ string) {
			var resBody registerAPI.PostResp
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
			var resBody registerAPI.PostResp
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
		inviteToken    string
		wantStatusCode int
		assertFunc     func(*testing.T, *http.Response, string)
	}{
		{
			name:           "UsnEmpty,PwdEmpty",
			username:       "",
			password:       "",
			inviteToken:    "",
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
			inviteToken:    "",
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
			inviteToken:    "",
			wantStatusCode: http.StatusBadRequest,
			assertFunc: assertOnValidationErrs(
				[]string{"Username is already taken."}, []string{},
			),
		},
		{
			name:           "InviteInvalid",
			username:       "bob321",
			password:       "Myp4ssw0rd!",
			inviteToken:    "10249812049182",
			wantStatusCode: http.StatusBadRequest,
			assertFunc:     assertOnResErr("Invalid invite token."),
		},
		{
			name:           "Success",
			username:       "bob321",
			password:       "Myp4ssw0rd!",
			inviteToken:    "",
			wantStatusCode: http.StatusOK,
			assertFunc: func(t *testing.T, res *http.Response, _ string) {
				// might take some time for post to create user so tr once
				// a second 5 times just in case.
				out, err := svcDynamo.GetItem(
					context.Background(), &dynamodb.GetItemInput{
						TableName: &userTableName,
						Key: map[string]types.AttributeValue{
							"Username": &types.AttributeValueMemberS{
								Value: "bob321",
							},
						},
					},
				)

				var user usertable.User
				attributevalue.UnmarshalMap(out.Item, &user)

				if err != nil {
					t.Fatal(err)
				}
				if err = bcrypt.CompareHashAndPassword(
					user.Password, []byte("Myp4ssw0rd!"),
				); err != nil {
					t.Error(err)
				}

				// assert that the returned JWT is valid and has the correct
				// subject
				cookie := res.Cookies()[0]
				assert.True(t.Error, cookie.Secure)
				assert.Equal(t.Error, cookie.SameSite, http.SameSiteNoneMode)
				claims := jwt.MapClaims{}
				if _, err = jwt.ParseWithClaims(
					cookie.Value, &claims, func(token *jwt.Token) (any, error) {
						return []byte(jwtKey), nil
					},
				); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t.Error, claims["username"].(string), "bob321")

				// no invite token was sent - therefore user must be put as
				// admin and given a random guid as team ID
				assert.Equal(t.Error, claims["isAdmin"].(bool), true)
				_, err = uuid.Parse(claims["teamID"].(string))
				assert.Nil(t.Error, err)

				exp := claims["exp"].(float64)
				if exp > float64(time.Now().Add(1*time.Hour).Unix()) {
					t.Error("expiry was more than an hour")
				}
				if exp < float64(time.Now().Add(59*time.Minute).Unix()) {
					t.Error("expiry was less than an hour")
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
			req := httptest.NewRequest(
				http.MethodPost, "/", bytes.NewReader(reqBody),
			)
			if c.inviteToken != "" {
				req.AddCookie(&http.Cookie{
					Name:  "invite-token",
					Value: c.inviteToken,
				})
			}
			w := httptest.NewRecorder()

			sut.Handle(w, req, "")

			res := w.Result()
			assert.Equal(t.Error, res.StatusCode, c.wantStatusCode)
			c.assertFunc(t, res, "")
		})
	}
}
