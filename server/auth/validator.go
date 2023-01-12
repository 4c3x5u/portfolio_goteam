package auth

import "github.com/golang-jwt/jwt/v4"

// TokenValidator describes a type that can be used to validate an auth token.
type TokenValidator interface {
	Validate(string) (string, error)
}

// JWTValidator can be used to validate a JWT.
type JWTValidator struct{ key []byte }

// NewJWTValidator creates and returns a new JWTValidator.
func NewJWTValidator(key string) JWTValidator {
	return JWTValidator{key: []byte(key)}
}

// Validate validates a JWT and returns its subject.
func (v JWTValidator) Validate(token string) (string, error) {
	claims := &jwt.RegisteredClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (any, error) {
		return v.key, nil
	})
	return claims.Subject, err
}
