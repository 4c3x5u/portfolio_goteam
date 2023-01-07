package auth

import (
	"testing"
	"time"

	"server/assert"

	"github.com/golang-jwt/jwt/v4"
)

func TestGeneratorToken(t *testing.T) {
	var (
		username = "bob21"
		expiry   = time.Now().Add(1 * time.Hour)
		sut      = NewGeneratorToken("d16889c5-5e2e-48ed-87c4-d29b8ee23fad", jwt.SigningMethodHS256)
	)

	tokenStr, err := sut.Generate(username, expiry)
	assert.Nil(t, err)

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		assert.Equal(t, true, ok)
		return []byte(sut.key), nil
	})
	assert.Nil(t, err)
	claims, ok := token.Claims.(jwt.MapClaims)
	assert.Equal(t, true, ok)
	assert.Equal(t, nil, claims.Valid())
	assert.Equal(t, username, claims["sub"].(string))
	assert.Equal(t, float64(expiry.Unix()), claims["exp"].(float64))
}
