// Package cookie contains code for working with HTTP cookies.
package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// TokenGenerator describes a type that can be used to generate an
// authentication token string that is valid until the expiry time for the
// given subject (i.e. username).
type TokenGenerator interface {
	Generate(sub string, exp time.Time) (string, error)
}

// JWTGenerator can be used to generate a JWT token that is valid until
// the expiry time for the given subject (i.e. username).
type JWTGenerator struct{ key []byte }

// NewJWTGenerator creates and returns a new JWTCookieGenerator.
func NewJWTGenerator(key string) JWTGenerator {
	return JWTGenerator{key: []byte(key)}
}

// Generate generates a JWT token as a *http.Cookie that is valid until the
// expiry time for the given subject (i.e. username)
func (g JWTGenerator) Generate(sub string, exp time.Time) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub, "exp": exp.Unix(),
	}).SignedString(g.key)
}

