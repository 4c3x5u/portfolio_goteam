// Package cookie contains code for working with HTTP cookies.
package auth

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// CookieGenerator describes a type that can be used to generate an
// authentication cookie that is valid until the expiry time for the given
// subject (i.e. username).
type CookieGenerator interface {
	Generate(sub string, exp time.Time) (*http.Cookie, error)
}

// CookieName is the name of the cookie that the auth token is stored id.
const CookieName = "authToken"

// JWTCookieGenerator can be used to generate a JWT token that is valid until
// the expiry time for the given subject (i.e. username).
type JWTCookieGenerator struct{ key []byte }

// NewJWTCookieGenerator creates and returns a new JWTCookieGenerator.
func NewJWTCookieGenerator(key string) JWTCookieGenerator {
	return JWTCookieGenerator{key: []byte(key)}
}

// Generate generates a JWT token as a *http.Cookie that is valid until the
// expiry time for the given subject (i.e. username)
func (g JWTCookieGenerator) Generate(sub string, exp time.Time) (*http.Cookie, error) {
	if token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": exp.Unix(),
	}).SignedString(g.key); err != nil {
		return nil, err
	} else {
		return &http.Cookie{
			Name:    "authToken",
			Value:   token,
			Expires: exp.UTC(),
		}, nil
	}
}
