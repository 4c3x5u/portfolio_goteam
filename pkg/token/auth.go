package token

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AuthName is the name of the auth token.
const AuthName = "auth-token"

// AuthDurationDefault is the default amount of time that an auth token lasts.
const AuthDurationDefault = 1 * time.Hour

// Auth defines the body of an Auth token.
type Auth struct {
	Username string
	IsAdmin  bool
	TeamID   string
}

// NewAuth creates and returns a new Auth.
func NewAuth(username string, isAdmin bool, teamID string) Auth {
	return Auth{Username: username, TeamID: teamID, IsAdmin: isAdmin}
}

// EncodeAuth encodes an Auth into a JWT string.
func EncodeAuth(exp time.Time, auth Auth) (string, error) {
	tk, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": auth.Username,
		"isAdmin":  auth.IsAdmin,
		"teamID":   auth.TeamID,
		"exp":      exp.Unix(),
	}).SignedString([]byte(os.Getenv(keyName)))
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
			return []byte(os.Getenv(keyName)), nil
		},
	); err != nil {
		return Auth{}, err
	}

	username, ok := claims["username"].(string)
	if !ok {
		return Auth{}, ErrInvalid
	}

	isAdmin, ok := claims["isAdmin"].(bool)
	if !ok {
		return Auth{}, ErrInvalid
	}

	teamID, ok := claims["teamID"].(string)
	if !ok {
		return Auth{}, ErrInvalid
	}

	return NewAuth(username, isAdmin, teamID), nil
}
