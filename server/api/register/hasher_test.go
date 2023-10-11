//go:build utest

package register

import (
	"testing"

	"github.com/kxplxn/goteam/server/assert"

	"golang.org/x/crypto/bcrypt"
)

// TestPasswordHasher tests the Hash method of the password hasher. It uses
// bcrypt.CompareHashAndPassword to assert that the result was generated from
// the given plaintext and doesn't match another plaintext string.
func TestPasswordHasher(t *testing.T) {
	sut := NewPasswordHasher()

	for _, c := range []struct {
		name           string
		inPlaintext    string
		matchPlaintext string
		wantErr        error
	}{
		{
			name:           "NoMatch",
			inPlaintext:    "password",
			matchPlaintext: "differentPassword",
			wantErr:        bcrypt.ErrMismatchedHashAndPassword,
		},
		{
			name:           "Match",
			inPlaintext:    "password",
			matchPlaintext: "password",
			wantErr:        nil,
		},
	} {
		t.Run(c.name, func(t *testing.T) {
			hash, err := sut.Hash(c.inPlaintext)

			if err = assert.Nil(err); err != nil {
				t.Error(err)
			}
			err = bcrypt.CompareHashAndPassword(hash, []byte(c.matchPlaintext))
			if err = assert.Equal(c.wantErr, err); err != nil {
				t.Error(err)
			}
		})
	}
}
