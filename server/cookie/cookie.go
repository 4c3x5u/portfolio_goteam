// Package cookie contains code for working with HTTP cookies.
package cookie

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// AuthGenerator describes a type that can be used to generate an
// authentication cookie that is valid until the expiry time for the given
// subject (i.e. username).
type AuthGenerator interface {
	Generate(sub string, exp time.Time) (*http.Cookie, error)
}

// JWTGenerator can be used to generate a JWT token that is valid until
// the expiry time for the given subject (i.e. username).
type JWTGenerator struct{ key string }

// NewJWTGenerator creates and returns a new JWT cookie generator.
func NewJWTGenerator(key string) JWTGenerator { return JWTGenerator{key: key} }

// Generate generates a JWT token that is valid until the expiry time for the
// given subject (i.e. username)
func (g JWTGenerator) Generate(sub string, exp time.Time) (*http.Cookie, error) {
	if token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": sub,
		"exp": exp.Unix(),
	}).SignedString([]byte(g.key)); err != nil {
		return nil, err
	} else {
		return &http.Cookie{
			Name:    "authToken",
			Value:   token,
			Expires: exp.UTC(),
		}, nil
	}
}
