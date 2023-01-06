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

// Generate generates a token for the user associated with a given username.
func (g *GeneratorToken) Generate(username string) (string, error) {
	now := time.Now()
	return jwt.NewWithClaims(g.methodSigning, jwt.MapClaims{
		"sub": username,
		"exp": now.Add(1 * time.Hour).Unix(),
	}).SignedString([]byte(g.key))
}
