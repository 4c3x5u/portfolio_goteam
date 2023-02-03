//go:build utest

package auth

import (
	"testing"

	"server/assert"
)

// TestJWTValidator tests the Validate method of the JWTValidator to assert
// that the correct errors get returned based on the token passed in.
func TestJWTValidator(t *testing.T) {
	validKey := "ASLDJFLASKDJFLAKSDJFALSDKJAFLSDK"
	for _, c := range []struct {
		name    string
		key     string
		token   string
		wantSub string
	}{
		{
			name:    "InvalidSignature",
			key:     "INVALIDKEYINVALIDKEYINVALIDKEY",
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2IyMSJ9.k6QDVjyaHxPYixeoQBLixC5c79VK-WZ_kD9u4fjX_Ks",
			wantSub: "",
		},
		{
			name:    "MalformedToken",
			key:     validKey,
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCb2IyMSJ9.k6QDVjyaHxPYixeoQBLixC5c79VK-WZ_kD9u4fjX_Ks",
			wantSub: "",
		},
		{
			name:    "Expired",
			key:     validKey,
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzM1MzQ2NDIsInN1YiI6ImJvYjIxIn0.6Ii9QWGyjY5Q1TMoI6W5QdiTB4Fhy87aD3QZYbxpmn4",
			wantSub: "",
		},
		{
			name:    "Success",
			key:     validKey,
			token:   "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2IyMSJ9.k6QDVjyaHxPYixeoQBLixC5c79VK-WZ_kD9u4fjX_Ks",
			wantSub: "bob21",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			sut := NewJWTValidator(c.key)

			sub := sut.Validate(c.token)

			if err := assert.Equal(c.wantSub, sub); err != nil {
				t.Error(err)
			}
		})
	}
}
