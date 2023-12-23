package cookie

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AuthName is the name of the auth token.
const AuthName = "auth-token"

// Auth defines the body of an Auth token.
type Auth struct {
	Username string
	IsAdmin  bool
	TeamID   string
}

// NewAuth creates and returns a new Auth.
func NewAuth(username string, isAdmin bool, teamID string) Auth {
	return Auth{Username: username, IsAdmin: isAdmin, TeamID: teamID}
}

// EncoderAuth defines a type that can be used to encode an auth token.
type EncoderAuth struct {
	key []byte
	dur time.Duration
}

// NewAuthEncoder creates and returns a new AuthEncoder.
func NewAuthEncoder(jwtKey []byte, duration time.Duration) EncoderAuth {
	return EncoderAuth{key: jwtKey, dur: duration}
}

// Encode encodes an Auth into a JWT string.
func (e EncoderAuth) Encode(auth Auth) (http.Cookie, error) {
	exp := time.Now().Add(e.dur)

	tk, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": auth.Username,
		"isAdmin":  auth.IsAdmin,
		"teamID":   auth.TeamID,
		"exp":      exp.Unix(),
	}).SignedString(e.key)
	if err != nil {
		return http.Cookie{}, err
	}

	return http.Cookie{
		Name:     AuthName,
		Value:    tk,
		Expires:  exp.UTC(),
		SameSite: http.SameSiteNoneMode,
		Secure:   true,
	}, nil
}

// AuthDecoder defines a type that can be used to decode an auth token.
type AuthDecoder struct{ key []byte }

// NewAuthDecoder creates and returns a new AuthDecoder.
func NewAuthDecoder(key []byte) AuthDecoder { return AuthDecoder{key: key} }

// Decode validates and decodes a raw JWT string into an Auth.
func (d AuthDecoder) Decode(ck http.Cookie) (Auth, error) {
	if ck.Value == "" {
		return Auth{}, ErrInvalid
	}

	claims := jwt.MapClaims{}
	if _, err := jwt.ParseWithClaims(
		ck.Value, &claims, func(token *jwt.Token) (any, error) {
			return d.key, nil
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
