//go:build utest

package auth

import (
	"testing"
	"time"

	"github.com/kxplxn/goteam/pkg/assert"

	"github.com/golang-jwt/jwt/v4"
)

// TestJWTGenerator tests the JWTCookieGenerator's Generate method to ensure
// that the generated JWT and the format of the returned *http.Cookie are valid.
func TestJWTGenerator(t *testing.T) {
	var (
		username = "bob123"
		expiry   = time.Now().Add(1 * time.Hour)
		sut      = NewJWTGenerator("d16889c5-5e2e-48ed-87c4-d29b8ee23fad")
	)

	token, err := sut.Generate(username, expiry)
	if err != nil {
		t.Fatal(err)
	}

	claims := &jwt.RegisteredClaims{}
	if _, err = jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (any, error) { return sut.key, nil },
	); err != nil {
		t.Error(err)
	}
	assert.Equal(t.Error, claims.Subject, username)
	assert.Equal(t.Error, claims.ExpiresAt.Unix(), expiry.Unix())
}
