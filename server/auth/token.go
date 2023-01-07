package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
)

// GeneratorToken can be used to generate a JWT token for a given user that is
// valid for an hour.
type GeneratorToken struct {
	key           string
	methodSigning jwt.SigningMethod
}

// NewGeneratorToken is the constructor for GeneratorToken.
func NewGeneratorToken(key string, methodSigning jwt.SigningMethod) *GeneratorToken {
	return &GeneratorToken{key: key, methodSigning: methodSigning}
}

// Generate generates a JWT for the user associated with a given username.
// The sub argument is the subject ID (e.g. username), and exp is for expiry.
func (g *GeneratorToken) Generate(sub string, exp time.Time) (string, error) {
	return jwt.NewWithClaims(g.methodSigning, jwt.MapClaims{
		"sub": sub,
		"exp": exp.Unix(),
	}).SignedString([]byte(g.key))
}
