package cookie

import (
	"testing"
	"time"

	"server/assert"

	"github.com/golang-jwt/jwt/v4"
)

func TestJWTGenerator(t *testing.T) {
	var (
		username = "bob21"
		expiry   = time.Now().Add(1 * time.Hour)
		sut      = NewJWTGenerator("d16889c5-5e2e-48ed-87c4-d29b8ee23fad")
	)

	cookie, err := sut.Generate(username, expiry)
	if err != nil {
		t.Fatal(err)
	}

	err = assert.Equal("authToken", cookie.Name)
	err = assert.Equal(expiry.UTC(), cookie.Expires)
	if err != nil {
		t.Error(err)
	}

	token, err := jwt.Parse(cookie.Value, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if err = assert.True(ok); err != nil {
			t.Error(err)
		}
		return []byte(sut.key), nil
	})
	if err = assert.Nil(err); err != nil {
		t.Error(err)
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if err = assert.True(ok); err != nil {
		t.Error(err)
	}
	if err = assert.Nil(claims.Valid()); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(username, claims["sub"].(string)); err != nil {
		t.Error(err)
	}
	if err = assert.Equal(float64(expiry.Unix()), claims["exp"].(float64)); err != nil {
		t.Error(err)
	}
}
