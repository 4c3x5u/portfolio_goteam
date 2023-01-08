// Package token contains code for working with authentication tokens.
package token

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// Generator describes a type that can be used to generate a string token for
// a given subject and expiry time.
type Generator interface {
	Generate(sub string, exp time.Time) (string, error)
}

// JWTGenerator can be used to generate a JWT token for a given user.
type JWTGenerator struct {
	key           string
	signingMethod jwt.SigningMethod
}

// NewJWTGenerator is the constructor for JWTGenerator.
func NewJWTGenerator(key string, signingMethod jwt.SigningMethod) JWTGenerator {
	return JWTGenerator{key: key, signingMethod: signingMethod}
}

// Generate generates a JWT for the user associated with a given username.
// The sub argument is the subject ID (e.g. username), and exp is for expiry.
func (g JWTGenerator) Generate(sub string, exp time.Time) (string, error) {
	return jwt.NewWithClaims(g.signingMethod, jwt.MapClaims{
		"sub": sub,
		"exp": exp.Unix(),
	}).SignedString([]byte(g.key))
}
