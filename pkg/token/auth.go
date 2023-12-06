package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// NameAuth is the name of the auth token.
const NameAuth = "auth"

// Auth defines the body of an Auth token.
type Auth struct {
	UserID  string
	IsAdmin bool
	TeamID  string
}

// NewAuth creates and returns a new Auth.
func NewAuth(userID string, isAdmin bool, teamID string) Auth {
	return Auth{UserID: userID, TeamID: teamID, IsAdmin: isAdmin}
}

// EncodeAuth encodes an Auth into a JWT string.
func EncodeAuth(exp time.Time, auth Auth) (string, error) {
	tk, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":  auth.UserID,
		"isAdmin": auth.IsAdmin,
		"teamID":  auth.TeamID,
		"exp":     exp.Unix(),
	}).SignedString([]byte(signKey))
	return tk, err
}

// Decode validates and decodes a raw JWT string into an Auth.
func DecodeAuth(raw string) (Auth, error) {
	if raw == "" {
		return Auth{}, ErrInvalid
	}

	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(
		raw, &claims, func(token *jwt.Token) (any, error) {
			return signKey, nil
		},
	); err != nil {
		return Auth{}, err
	}

	return NewAuth(
		claims["userID"].(string),
		claims["isAdmin"].(bool),
		claims["teamID"].(string),
	), nil
}
