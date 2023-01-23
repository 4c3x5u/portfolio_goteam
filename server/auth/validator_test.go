package auth

import (
	"testing"

	"server/assert"

	"github.com/golang-jwt/jwt/v4"
)

// TestJWTValidator tests the Validate method of the JWTValidator to assert
// that the correct errors get returned based on the token passed in.
func TestJWTValidator(t *testing.T) {
	validKey := "ASLDJFLASKDJFLAKSDJFALSDKJAFLSDK"
	for _, c := range []struct {
		name          string
		sub           string
		key           string
		token         string
		validationErr error
	}{
		{
			name:          "InvalidSignature",
			sub:           "bob21",
			key:           "INVALIDKEYINVALIDKEYINVALIDKEY",
			token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2IyMSJ9.k6QDVjyaHxPYixeoQBLixC5c79VK-WZ_kD9u4fjX_Ks",
			validationErr: jwt.ErrSignatureInvalid,
		},
		{
			name:          "MalformedToken",
			sub:           "",
			key:           validKey,
			token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCb2IyMSJ9.k6QDVjyaHxPYixeoQBLixC5c79VK-WZ_kD9u4fjX_Ks",
			validationErr: jwt.ErrTokenMalformed,
		},
		{
			name:          "Expired",
			sub:           "bob21",
			key:           validKey,
			token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzM1MzQ2NDIsInN1YiI6ImJvYjIxIn0.6Ii9QWGyjY5Q1TMoI6W5QdiTB4Fhy87aD3QZYbxpmn4",
			validationErr: jwt.ErrTokenExpired,
		},
		{
			name:          "Success",
			sub:           "bob21",
			key:           validKey,
			token:         "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2IyMSJ9.k6QDVjyaHxPYixeoQBLixC5c79VK-WZ_kD9u4fjX_Ks",
			validationErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			sut := NewJWTValidator(c.key)
			sub, err := sut.Validate(c.token)

			if c.validationErr == nil {
				if err = assert.Nil(err); err != nil {
					t.Error(err)
				}
			} else {
				// Due to the way that validation errors are implemented in jwt-go,
				// this seems to be the most reasonable way to check what error we
				// get.
				if err = assert.True(err.(*jwt.ValidationError).Is(c.validationErr)); err != nil {
					t.Error(err)
				}
			}
			if err = assert.Equal(c.sub, sub); err != nil {
				t.Error(err)
			}
		})
	}
}
