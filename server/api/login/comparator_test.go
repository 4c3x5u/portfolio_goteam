//go:build utest

package login

import (
	"testing"

	"github.com/kxplxn/goteam/server/assert"

	"golang.org/x/crypto/bcrypt"
)

// TestPasswordComparator tests the PasswordComparator's Compare method to assert that all it
// does is to call the bcrypt.CompareHashAndPassword and return whatever error
// that returns.
func TestPasswordComparator(t *testing.T) {
	sut := NewPasswordComparator()

	for _, c := range []struct {
		name        string
		inHash      []byte
		inPlaintext string
		wantErr     error
	}{
		{
			name:        "WrongPassword",
			inPlaintext: "password",
			inHash: []byte(
				"$2a$04$ngqMWrzBWyg8KO3MGk1cnOISt3wyeBwbFlkvghSHKBkSYOeO2.7XG",
			),
			wantErr: bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:        "BcryptError",
			inPlaintext: "password",
			inHash:      []byte("$2a$04$ngqMWrzBWyg8K"),
			wantErr:     bcrypt.ErrHashTooShort,
		},
		{
			name:        "Success",
			inPlaintext: "password",
			inHash: []byte(
				"$2a$04$W4ABZofxx5uoJVgTlYS1wuFHz1LLQaBfoO0iwz/04WWmg9LQdCPsS",
			),
			wantErr: nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			err := sut.Compare(c.inHash, c.inPlaintext)

			assert.Equal(t.Error, err, c.wantErr)
		})
	}
}
