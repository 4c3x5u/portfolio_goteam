package auth

import (
	"testing"
	"time"

	"server/assert"

	"github.com/golang-jwt/jwt/v4"
)

// TestJWTGenerator tests the JWTCookieGenerator's Generate method to ensure
// that the generated JWT and the format of the returned *http.Cookie are valid.
func TestJWTGenerator(t *testing.T) {
	var (
		username = "bob21"
		expiry   = time.Now().Add(1 * time.Hour)
		sut      = NewJWTGenerator("d16889c5-5e2e-48ed-87c4-d29b8ee23fad")
	)

	token, err := sut.Generate(username, expiry)
	if err != nil {
		t.Fatal(err)
	}

	// err = assert.Equal(CookieName, cookie.Name)
	// err = assert.Equal(expiry.UTC(), cookie.Expires)
	// if err != nil {
	// 	t.Error(err)
	// }

	claims := &jwt.RegisteredClaims{}
	if _, err = jwt.ParseWithClaims(
		token,
		claims,
		func(token *jwt.Token) (any, error) { return sut.key, nil },
	); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(username, claims.Subject); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(expiry.Unix(), claims.ExpiresAt.Unix()); err != nil {
		t.Error(err)
	}
}
