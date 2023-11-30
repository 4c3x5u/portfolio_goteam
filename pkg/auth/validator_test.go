//go:build utest

package auth

import (
	"testing"

	"github.com/kxplxn/goteam/pkg/assert"
)

// TestJWTValidator tests the Validate method of the JWTValidator to assert
// that the correct errors get returned based on the token passed in.
func TestJWTValidator(t *testing.T) {
	sut := NewJWTValidator("ASLDJFLASKDJFLAKSDJFALSDKJAFLSDK")

	for _, c := range []struct {
		name    string
		token   string
		wantSub string
	}{
		{
			name: "InvalidSignature",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2IyMSJ9" +
				".w_B9yQkrWU3s5vdD7YJn4hAutfUMxtb4JfdQvpfeiP0",
			wantSub: "",
		},
		{
			name: "MalformedToken",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCb2IyMSJ9.k6QDVjyaHxPYix" +
				"eoQBLixC5c79VK-WZ_kD9u4fjX_Ks",
			wantSub: "",
		},
		{
			name: "Expired",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NzM1MzQ2" +
				"NDIsInN1YiI6ImJvYjIxIn0.6Ii9QWGyjY5Q1TMoI6W5QdiTB4Fhy87aD3QZ" +
				"Ybxpmn4",
			wantSub: "",
		},
		{
			name: "Success",
			token: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiJib2IyMSJ9" +
				".k6QDVjyaHxPYixeoQBLixC5c79VK-WZ_kD9u4fjX_Ks",
			wantSub: "bob21",
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			sub := sut.Validate(c.token)

			assert.Equal(t.Error, c.wantSub, sub)
		})
	}
}
